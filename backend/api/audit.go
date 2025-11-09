package api

type AuditRequestDto struct {
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	EntityID  uint   `json:"entity_id"`
	UserEmail string `json:"user_email"`
}

type AuditEditRequestDto struct {
	Action    *string `json:"action,omitempty"`
	Entity    *string `json:"entity,omitempty"`
	EntityID  *uint   `json:"entity_id,omitempty"`
	UserEmail *string `json:"user_email,omitempty"`
}

type AuditResponseDto struct {
	ID        int    `json:"id"`
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	EntityID  uint   `json:"entity_id"`
	UserEmail string `json:"user_email"`
	CreatedAt string `json:"created_at"`
}
