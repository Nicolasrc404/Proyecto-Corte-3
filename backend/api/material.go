package api

type MaterialRequestDto struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Quantity float64 `json:"quantity"`
}

type MaterialResponseDto struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Quantity  float64 `json:"quantity"`
	CreatedAt string  `json:"created_at"`
}
