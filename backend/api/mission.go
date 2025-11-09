package api

type MissionRequestDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	AssignedTo  uint   `json:"assigned_to"`
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

type MissionEditRequestDto struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Difficulty  *string `json:"difficulty,omitempty"`
	Status      *string `json:"status,omitempty"`
	AssignedTo  *uint   `json:"assigned_to,omitempty"`
}
