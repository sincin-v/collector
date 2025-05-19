package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"

	b64 "encoding/base64"

	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/server/config"
)

func CheckHashRequest(cfg config.Config) func(h http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {
		checkerFunction := func(w http.ResponseWriter, r *http.Request) {
			requestHash := r.Header.Get("HashSHA256")
			if requestHash != "" {

				receivedBody, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Log.Error("Error read request body: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				hmac := hmac.New(sha256.New, []byte(cfg.SecretKey))
				hmac.Write(receivedBody)
				hashSum := hmac.Sum(nil)
				resultHash := b64.StdEncoding.EncodeToString([]byte(hashSum))
				if requestHash != resultHash {
					logger.Log.Error("Invalid receive hash!")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(receivedBody))
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(checkerFunction)
	}
}
