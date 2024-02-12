package handlers

import (
	"net/http"

	"github.com/dbubel/intake"
	"github.com/dbubel/vision/models/requests"
	"github.com/dbubel/vision/pkg/validate"
	"github.com/julienschmidt/httprouter"
)

func (c *App) Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newVectorSearchRequest requests.Search
	err := validate.UnmarshalJSON(r.Body, &newVectorSearchRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "invalid create vector request",
		})
		return
	}

	vecTable, err := c.DB.ValidateInsert(r.Context(), newVectorSearchRequest.Workspace, newVectorSearchRequest.VectorRepo)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "workspace repo pair not found",
		})
		return
	}
	res, err := c.DB.Search(r.Context(), vecTable, newVectorSearchRequest.Data.Data)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "error searching for vector",
		})
		return
	}

	intake.RespondJSON(w, r, http.StatusOK, res)
}
