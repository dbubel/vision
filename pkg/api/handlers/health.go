package handlers

import (
	"net/http"

	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
)

type Err map[string]interface{}

func (c *App) Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	intake.RespondJSON(w, r, http.StatusOK, "OK")
	// intake.RespondJSON(w, r, http.StatusOK, struct {
	// 	BuildTag     string `json:"buildTag"`
	// 	BuildDate    string `json:"buildDate"`
	// 	GitHash      string `json:"commitHash"`
	// 	InstanceName string `json:"instanceName"`
	// 	UpTime       string `json:"upTime"`
	// }{
	// 	BuildDate:    c.BuildDate,
	// 	BuildTag:     c.BuildTag,
	// 	GitHash:      c.GitHash,
	// 	InstanceName: c.InstanceName,
	// 	UpTime:       time.Now().Sub(c.UpTime).String(),
	// })
}
func (c *App) Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	intake.RespondJSON(w, r, http.StatusOK, "OK")
	// intake.RespondJSON(w, r, http.StatusOK, struct {
	// 	BuildTag     string `json:"buildTag"`
	// 	BuildDate    string `json:"buildDate"`
	// 	GitHash      string `json:"commitHash"`
	// 	InstanceName string `json:"instanceName"`
	// 	UpTime       string `json:"upTime"`
	// }{
	// 	BuildDate:    c.BuildDate,
	// 	BuildTag:     c.BuildTag,
	// 	GitHash:      c.GitHash,
	// 	InstanceName: c.InstanceName,
	// 	UpTime:       time.Now().Sub(c.UpTime).String(),
	// })
}
