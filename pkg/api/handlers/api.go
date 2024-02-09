package handlers

import (
	"github.com/dbubel/vision/config"
	"github.com/dbubel/vision/store"
	"github.com/sirupsen/logrus"
)

type App struct {
	Cfg config.Config
	Log *logrus.Logger
	DB  *store.Store
}

func NewApp(c config.Config, l *logrus.Logger, s *store.Store) *App {
	return &App{
		Cfg: c,
		Log: l,
		DB:  s,
	}
}
