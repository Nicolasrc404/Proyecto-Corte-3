package api

type MaterialRequestDto struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Quantity float64 `json:"quantity"`
}

type MaterialEditRequestDto struct {
	Name     *string  `json:"name,omitempty"`
	Category *string  `json:"category,omitempty"`
	Quantity *float64 `json:"quantity,omitempty"`
}

type MaterialResponseDto struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Quantity  float64 `json:"quantity"`
	CreatedAt string  `json:"created_at"`
}
