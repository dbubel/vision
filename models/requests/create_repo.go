package requests

type CreateRepoRequest struct {
	Name      string `json:"name" validate:"required"`
	Workspace string `json:"workspace" validate:"required"`
	Dimension int    `json:"dimension" validate:"required"`
}
