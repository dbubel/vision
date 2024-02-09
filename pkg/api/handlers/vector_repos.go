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
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &newVectorRequest)
	// err := json.Marshaler(r.Body, &newVectorRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "invalid create vector request",
		})
		return
	}
	err = c.DB.InsertVectors(r.Context(), "vec_285waarjusddapcf9gmreqo7h", newVectorRequest)
	if err != nil {
		intake.RespondJSON(w, r, http.StatusBadRequest, Err{
			"Error":       err.Error(),
			"Description": "error inserting vector",
		})
		return
	}

	intake.Respond(w, r, http.StatusCreated, []byte("ok"))
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

	intake.RespondJSON(w, r, http.StatusCreated, res)
}
