package api

type AlchemistRequestDto struct {
	Name      string `json:"name"`
	Age       int32  `json:"age"`
	Specialty string `json:"specialty"`
	Rank      string `json:"rank"`
}

type AlchemistEditRequestDto struct {
	Name      *string `json:"name,omitempty"`
	Age       *int32  `json:"age,omitempty"`
	Specialty *string `json:"specialty,omitempty"`
	Rank      *string `json:"rank,omitempty"`
}

type AlchemistResponseDto struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Specialty string `json:"specialty"`
	Rank      string `json:"rank"`
	CreatedAt string `json:"created_at"`
}
