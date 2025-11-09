package api

type TransmutationRequestDto struct {
	AlchemistID uint   `json:"alchemist_id"`
	MaterialID  uint   `json:"material_id"`
	Formula     string `json:"formula"`
}

type TransmutationResponseDto struct {
	ID          int    `json:"id"`
	AlchemistID uint   `json:"alchemist_id"`
	MaterialID  uint   `json:"material_id"`
	Status      string `json:"status"`
	Result      string `json:"result"`
	CreatedAt   string `json:"created_at"`
}

type TransmutationEditRequestDto struct {
	Formula *string `json:"formula,omitempty"`
	Status  *string `json:"status,omitempty"`
	Result  *string `json:"result,omitempty"`
}
