package config

type Config struct {
	Address                     string  `json:"address"`
	Database                    string  `json:"database"`
	KillDuration                int     `json:"kill_duration"`
	KillDurationWithDescription int     `json:"kill_duration_with_desc"`
	RedisAddress                string  `json:"redis_address"`
	VerificationIntervalMinutes int     `json:"verification_interval_minutes"`
	PendingTransmutationHours   int     `json:"pending_transmutation_hours"`
	MaterialLowStockThreshold   float64 `json:"material_low_stock_threshold"`
}
