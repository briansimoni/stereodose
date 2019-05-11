package util

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// TransactionIDKey can be used to grab the unique ID from request context
const TransactionIDKey = "TransactionID"

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// addTransactionID takes an http request and adds a unique ID to the context
// It attempts to grab the ID from an AWS load balancer, but adds one if not
// running in an environment with such a load balancer.
func addTransactionID(r *http.Request) (*http.Request, string) {
	var transactionID string
	transactionID = r.Header.Get("X-Amzn-Trace-Id")
	if transactionID == "" {
		transactionID = fmt.Sprintf("%s-%d", uuid.New().String(), time.Now().Unix())
	}
	ctx := context.WithValue(r.Context(), TransactionIDKey, transactionID)
	req := r.WithContext(ctx)
	return req, transactionID
}

// RequestLogger complies to the gorilla webkit Middleware interface
// It is used to log detailed request data and adds a unique transaction ID
func RequestLogger(h http.Handler) http.Handler {
	next := func(w http.ResponseWriter, r *http.Request) {
		req, id := addTransactionID(r)
		sw := &statusWriter{ResponseWriter: w}
		h.ServeHTTP(sw, req)
		log.WithFields(logrus.Fields{
			"RemoteAddress":   req.RemoteAddr,
			"ForwardedScheme": req.Proto,
			"Host":            req.Host,
			"Method":          req.Method,
			"Path":            req.URL.Path,
			"UserAgent":       req.UserAgent(),
			"Referer":         req.Referer(),
			"TransactionID":   id,
			"StatusCode":      sw.status,
			"Length":          sw.length,
		}).Info("RequestLog")
	}
	return http.HandlerFunc(next)
}
