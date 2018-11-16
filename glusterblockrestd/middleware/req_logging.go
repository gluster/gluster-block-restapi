package middleware

import (
	"net/http"

	"github.com/gluster/gluster-block-restapi/pkg/context"

	"github.com/gorilla/handlers"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

// WithReqLogger is a http middleware which will log incoming
// http request.
func WithReqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			reqID     = uuid.NewRandom()
			logWriter = log.WithField("reqid", reqID.String())
			ctx       = r.Context()
		)
		{
			ctx = context.WithReqID(ctx, reqID)
			ctx = context.WithReqLogger(ctx, logWriter)
		}
		handlers.LoggingHandler(logWriter.Writer(), next).ServeHTTP(w, r.WithContext(ctx))
	})
}
