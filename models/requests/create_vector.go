package requests

type (
	// Vector  []float32
	Vectors struct {
		Workspace  string      `json:"workspace" validate:"required"`
		VectorRepo string      `json:"vectorRepo" validate:"required"`
		Data       [][]float32 `json:"data" validate:"required"`
	}
)
