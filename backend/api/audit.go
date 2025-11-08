package api

type AuditResponseDto struct {
	ID        int    `json:"id"`
	Entity    string `json:"entity"`
	EntityID  uint   `json:"entity_id"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}
