package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dbubel/intake"
	"github.com/dbubel/vision/models/requests"
	"github.com/dbubel/vision/pkg/validate"
	"github.com/julienschmidt/httprouter"
)

func (c *App) InsertVector(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newVectorRequest requests.Vectors
	// err := validate.UnmarshalJSON(r.Body, &newVectorRequest)
	b, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(b, &newVectorRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "invalid create vector request",
		})
		return
	}

	vecTable, err := c.DB.ValidateInsert(r.Context(), newVectorRequest.Workspace, newVectorRequest.VectorRepo)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "workspace repo pair not found",
		})
		return
	}

	err = c.DB.InsertVectors(r.Context(), vecTable, newVectorRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "error inserting vector",
		})
		return
	}

	intake.Respond(w, r, http.StatusAccepted, nil)
}

func (c *App) CreateRepo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newRepoRequest requests.CreateRepoRequest
	err := validate.UnmarshalJSON(r.Body, &newRepoRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "invalid create repo request",
		})
		return
	}
	res, err := c.DB.CreateVectorTableMapping(r.Context(), newRepoRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "error creating vector repo",
		})
		return
	}
	err = c.DB.CreateVectorTable(r.Context(), res.VectorTable, res.Dimension)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "error creating vector table",
		})
		return
	}
	// s, err := c.DB.GetIndexListSize(res.VectorTable)
	// if err != nil {
	// 	intake.RespondJSON(w, r, http.StatusBadRequest, Err{
	// 		"Error":       err.Error(),
	// 		"Description": "error getting index size",
	// 	})
	// 	return
	// }
	//
	// err = c.DB.CreateIndex(r.Context(), res.VectorTable, s)
	// if err != nil {
	// 	intake.RespondJSON(w, r, http.StatusBadRequest, Err{
	// 		"Error":       err.Error(),
	// 		"Description": "error creating index",
	// 	})
	// 	return
	// }
	intake.RespondJSON(w, r, http.StatusCreated, res)
}
