package handlers

import (
	"net/http"

	"github.com/dbubel/vision/models/requests"

	"github.com/dbubel/intake"
	"github.com/dbubel/vision/pkg/validate"
	"github.com/julienschmidt/httprouter"
)

func (c *App) CreateWorkspace(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newWorkspaceRequest requests.CreateWorkspaceRequest
	err := validate.UnmarshalJSON(r.Body, &newWorkspaceRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "invalid create workspace request",
		})
		return
	}
	res, err := c.DB.CreateWorkspace(r.Context(), newWorkspaceRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "error creating workspace",
		})
		return
	}

	intake.RespondJSON(w, r, http.StatusCreated, res)
}
