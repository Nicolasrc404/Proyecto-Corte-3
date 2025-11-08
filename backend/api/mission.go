package api

type MissionRequestDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	AssignedTo  uint   `json:"assigned_to"` // AlchemistID
}

type MissionEditRequestDto struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Difficulty  *string `json:"difficulty,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type MissionResponseDto struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Status      string `json:"status"`
	AssignedTo  uint   `json:"assigned_to"`
	CreatedAt   string `json:"created_at"`
}
