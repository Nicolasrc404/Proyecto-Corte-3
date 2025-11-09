package api

type AuditResponseDto struct {
	ID        int    `json:"id"`
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	EntityID  uint   `json:"entity_id"`
	UserEmail string `json:"user_email"`
	CreatedAt string `json:"created_at"`
}
