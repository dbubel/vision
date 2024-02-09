package middleware

import (
	"net/http"
	"time"

	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func  Logging(log *logrus.Logger, next intake.Handler) intake.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		t := time.Now()
		defer func() {
			var code int
			if err := intake.FromContext(r, "response-code", &code); err != nil {
				log.WithError(err).Error("error getting response code from context")
			}

			log.WithFields(logrus.Fields{
				"method":         r.Method,
				"route":          r.URL.Path,
				"responseTimeMs": time.Now().Sub(t).String(),
				"code":           code,
			}).Debug("handled request")
		}()
		next(w, r, params)
	}
}