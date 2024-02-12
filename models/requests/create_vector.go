package requests

type Vectors struct {
	Workspace  string `json:"workspace" validate:"required"`
	VectorRepo string `json:"vectorRepo" validate:"required"`
	Data       []Data `json:"data" validate:"required"`
}

type Data struct {
	Data     []float32
	Metadata any
}
type Search struct {
	Workspace  string `json:"workspace" validate:"required"`
	VectorRepo string `json:"vectorRepo" validate:"required"`
	Data       Data   `json:"data" validate:"required"`
}
