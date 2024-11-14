package middlewares

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"status_code"`
	Message string `json:"message"`
}

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		// Custom response writer to intercept errors
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		// Handle any errors that were written to the response writer
		if rr.statusCode >= 400 {
			handleHTTPError(w, rr.statusCode)
		}

	})
}

func handleHTTPError(w http.ResponseWriter, statusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := errorResponse{
		Code:    statusCode,
		Message: http.StatusText(statusCode),
	}

	json.NewEncoder(w).Encode(resp)
}

// Custom response writer to capture status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}
