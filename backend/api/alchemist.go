package api

type AlchemistRequestDto struct {
	Name      string `json:"name"`
	Age       int32  `json:"age"`
	Specialty string `json:"specialty"`
	Rank      string `json:"rank"`
}

type AlchemistResponseDto struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Specialty string `json:"specialty"`
	Rank      string `json:"rank"`
	CreatedAt string `json:"created_at"`
}
