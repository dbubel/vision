package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dbubel/intake"
	"github.com/dbubel/vision/config"
	"github.com/dbubel/vision/pkg/api/handlers"
	"github.com/dbubel/vision/pkg/middleware"
	"github.com/dbubel/vision/store"
	"github.com/sirupsen/logrus"
)

type ServeCommand struct {
	Config config.Config
	Log    *logrus.Logger
}

func (c *ServeCommand) Help() string {
	return ""
}

func (c *ServeCommand) Synopsis() string {
	return "Runs the vision server"
}

func (c *ServeCommand) Run(args []string) int {
	c.Log.WithFields(logrus.Fields{"env": c.Config.Environment}).Info("API Starting up...")

	// Init the app api
	app := intake.New(c.Log)
	app.AddGlobalMiddleware(middleware.Logging(c.Log))
	// Handle CORS for OPTIONS requests
	app.Router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	// add a custom not found handler and nr
	app.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	store, err := store.New(c.Log, c.Config.ReaderConnStr(), c.Config.WriterConnStr())
	fmt.Println(c.Config.ReaderConnStr(), c.Config.WriterConnStr())
	if err != nil {
		c.Log.WithError(err).Error("error creating store")
		return 1
	}

	//==========================================================================================
	visionAPI := handlers.NewApp(c.Config, c.Log, store)
	endpoints := handlers.Endpoints(visionAPI)
	app.AddEndpoints(endpoints)

	app.Run(&http.Server{
		Addr:           fmt.Sprintf(":%d", 3000),
		Handler:        app.Router,
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		MaxHeaderBytes: 1 << 20,
	})

	return 0
}
