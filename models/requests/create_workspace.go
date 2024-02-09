package requests

type CreateWorkspaceRequest struct {
	Name        string `json:"name" db:"name" validate:"required"`
	Description string `json:"description" db:"description"`
}
