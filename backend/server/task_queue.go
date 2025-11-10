package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"backend-avanzada/logger"
	"backend-avanzada/models"
	"backend-avanzada/repository"
)

const (
	taskTypeProcessTransmutation = "process_transmutation"
	taskTypeRegisterAudit        = "register_audit"
	taskTypeDailyVerification    = "daily_verification"
)

type queueTask struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type processTransmutationPayload struct {
	TransmutationID uint   `json:"transmutation_id"`
	RequestedBy     string `json:"requested_by"`
}

type registerAuditPayload struct {
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	EntityID  uint   `json:"entity_id"`
	UserEmail string `json:"user_email"`
	Details   string `json:"details"`
}

type dailyVerificationPayload struct {
	ExecutedAt time.Time `json:"executed_at"`
}

// TaskQueue orchestrates all background work for the application. It provides
// helpers for HTTP handlers to enqueue jobs and executes them in a dedicated
// worker that relies on Redis for coordination.
type TaskQueue struct {
	redis              *RedisClient
	logger             *logger.Logger
	ctx                context.Context
	cancel             context.CancelFunc
	transRepo          *repository.TransmutationRepository
	auditRepo          *repository.AuditRepository
	missionRepo        *repository.MissionRepository
	materialRepo       *repository.MaterialRepository
	verificationTicker *time.Ticker
	verificationEvery  time.Duration
	pendingThreshold   time.Duration
	lowStockThreshold  float64
	started            bool
}

func NewTaskQueue(redisAddr string, log *logger.Logger) *TaskQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskQueue{
		redis:             NewRedisClient(redisAddr),
		logger:            log,
		ctx:               ctx,
		cancel:            cancel,
		started:           false,
		lowStockThreshold: 5,
		verificationEvery: 24 * time.Hour,
		pendingThreshold:  24 * time.Hour,
	}
}

func (q *TaskQueue) WithRepositories(
	transRepo *repository.TransmutationRepository,
	auditRepo *repository.AuditRepository,
	missionRepo *repository.MissionRepository,
	materialRepo *repository.MaterialRepository,
) {
	q.transRepo = transRepo
	q.auditRepo = auditRepo
	q.missionRepo = missionRepo
	q.materialRepo = materialRepo
}

func (q *TaskQueue) ConfigureThresholds(verificationEvery, pendingThreshold time.Duration, lowStockThreshold float64) {
	if verificationEvery > 0 {
		q.verificationEvery = verificationEvery
	}
	if pendingThreshold > 0 {
		q.pendingThreshold = pendingThreshold
	}
	if lowStockThreshold > 0 {
		q.lowStockThreshold = lowStockThreshold
	}
}

// Start spins up the worker that consumes jobs from Redis.
func (q *TaskQueue) Start() error {
	if q.started {
		return nil
	}
	if err := q.redis.Ping(q.ctx); err != nil {
		return fmt.Errorf("async queue is not available: %w", err)
	}
	q.started = true
	go q.worker()
	return nil
}

// Stop gracefully cancels the worker and ticker.
func (q *TaskQueue) Stop() {
	q.cancel()
	if q.verificationTicker != nil {
		q.verificationTicker.Stop()
	}
}

// ScheduleDailyVerification enqueues verification jobs at the configured interval.
func (q *TaskQueue) ScheduleDailyVerification() {
	if !q.started {
		return
	}
	q.logger.Printf("[async] programando verificaciones cada %s", q.verificationEvery)
	q.verificationTicker = time.NewTicker(q.verificationEvery)
	go func() {
		// Ejecutar una verificación inmediata al arrancar.
		if err := q.enqueueDailyVerification(); err != nil {
			q.logger.Printf("[async] no se pudo encolar verificación inicial: %v", err)
		}
		for {
			select {
			case <-q.ctx.Done():
				return
			case <-q.verificationTicker.C:
				if err := q.enqueueDailyVerification(); err != nil {
					q.logger.Printf("[async] error encolando verificación diaria: %v", err)
				}
			}
		}
	}()
}

// EnqueueTransmutationProcessing schedules the heavy processing of a transmutation.
func (q *TaskQueue) EnqueueTransmutationProcessing(transmutationID uint, requestedBy string) error {
	payload := processTransmutationPayload{TransmutationID: transmutationID, RequestedBy: requestedBy}
	return q.enqueue(taskTypeProcessTransmutation, payload)
}

// EnqueueAudit registers an audit asynchronously so handlers do not block on DB writes.
func (q *TaskQueue) EnqueueAudit(action, entity string, entityID uint, userEmail, details string) error {
	payload := registerAuditPayload{Action: action, Entity: entity, EntityID: entityID, UserEmail: userEmail, Details: details}
	return q.enqueue(taskTypeRegisterAudit, payload)
}

func (q *TaskQueue) enqueueDailyVerification() error {
	payload := dailyVerificationPayload{ExecutedAt: time.Now().UTC()}
	return q.enqueue(taskTypeDailyVerification, payload)
}

func (q *TaskQueue) enqueue(taskType string, payload interface{}) error {
	if !q.started {
		return errors.New("async queue has not been started")
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := queueTask{Type: taskType, Payload: data}
	raw, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return q.redis.LPUSH(q.ctx, redisQueueKey, raw)
}

func (q *TaskQueue) worker() {
	for {
		select {
		case <-q.ctx.Done():
			return
		default:
		}
		data, err := q.redis.BRPOP(q.ctx, redisQueueKey)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			q.logger.Printf("[async] error leyendo cola: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		var task queueTask
		if err := json.Unmarshal(data, &task); err != nil {
			q.logger.Printf("[async] payload inválido: %v", err)
			continue
		}
		if err := q.dispatch(task); err != nil {
			q.logger.Printf("[async] error ejecutando tarea %s: %v", task.Type, err)
		}
	}
}

func (q *TaskQueue) dispatch(task queueTask) error {
	switch task.Type {
	case taskTypeProcessTransmutation:
		var payload processTransmutationPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return err
		}
		return q.handleTransmutation(payload)
	case taskTypeRegisterAudit:
		var payload registerAuditPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return err
		}
		return q.handleAudit(payload)
	case taskTypeDailyVerification:
		return q.handleDailyVerification()
	default:
		return fmt.Errorf("tipo de tarea desconocido: %s", task.Type)
	}
}

func (q *TaskQueue) handleTransmutation(payload processTransmutationPayload) error {
	if q.transRepo == nil {
		return errors.New("transmutation repository is not configured")
	}
	transmutation, err := q.transRepo.FindById(int(payload.TransmutationID))
	if err != nil {
		return err
	}
	if transmutation == nil {
		return fmt.Errorf("transmutación %d no encontrada", payload.TransmutationID)
	}
	if strings.EqualFold(transmutation.Status, "completada") {
		return nil
	}

	// Simula un trabajo costoso.
	time.Sleep(3 * time.Second)
	transmutation.Status = "completada"
	transmutation.Result = fmt.Sprintf("Transmutación %d procesada exitosamente", transmutation.ID)
	if _, err := q.transRepo.Save(transmutation); err != nil {
		return err
	}

	if q.auditRepo != nil {
		audit := registerAuditPayload{
			Action:    "process_transmutation",
			Entity:    "transmutation",
			EntityID:  transmutation.ID,
			UserEmail: payload.RequestedBy,
			Details:   transmutation.Result,
		}
		return q.handleAudit(audit)
	}
	return nil
}

func (q *TaskQueue) handleAudit(payload registerAuditPayload) error {
	if q.auditRepo == nil {
		return errors.New("audit repository is not configured")
	}
	audit := &models.Audit{
		Action:    payload.Action,
		Entity:    payload.Entity,
		EntityID:  payload.EntityID,
		UserEmail: payload.UserEmail,
		Details:   payload.Details,
	}
	_, err := q.auditRepo.Save(audit)
	return err
}

func (q *TaskQueue) handleDailyVerification() error {
	if q.auditRepo == nil {
		return errors.New("audit repository is not configured")
	}
	var details []string

	if q.transRepo != nil {
		threshold := time.Now().Add(-q.pendingThreshold)
		pending, err := q.transRepo.FindPendingBefore(threshold)
		if err != nil {
			return err
		}
		if len(pending) > 0 {
			details = append(details, fmt.Sprintf("%d transmutaciones pendientes", len(pending)))
		}
	}

	if q.missionRepo != nil {
		threshold := time.Now().Add(-q.pendingThreshold)
		open, err := q.missionRepo.FindOpenBefore(threshold)
		if err != nil {
			return err
		}
		if len(open) > 0 {
			details = append(details, fmt.Sprintf("%d misiones sin cerrar", len(open)))
		}
	}

	if q.materialRepo != nil {
		scarce, err := q.materialRepo.FindScarce(q.lowStockThreshold)
		if err != nil {
			return err
		}
		if len(scarce) > 0 {
			details = append(details, fmt.Sprintf("%d materiales con stock crítico", len(scarce)))
		}
	}

	if len(details) == 0 {
		details = append(details, "Sin hallazgos críticos")
	}

	audit := &models.Audit{
		Action:    "daily_verification",
		Entity:    "system",
		Details:   strings.Join(details, "; "),
		UserEmail: "system",
	}
	_, err := q.auditRepo.Save(audit)
	return err
}

// asyncErrorReporter creates a helper that handlers can use to report async issues.
func (s *Server) asyncErrorReporter() func(path string, err error) {
	return func(path string, err error) {
		if err == nil {
			return
		}
		s.logger.Error(http.StatusInternalServerError, fmt.Sprintf("%s [async]", path), err)
	}
}

// currentUserExtractor returns the email stored in the JWT claims.
func currentUserExtractor(r *http.Request) string {
	if claims := GetAuthClaims(r); claims != nil {
		return claims.Email
	}
	return ""
}
