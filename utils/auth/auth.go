package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sajallimbu/go_securing_api/models"
)

// Exception ... our exception message struct
type Exception struct {
	ResponseCode int
	Message      string
}

//JwtVerify ... our middleware that checks if the JWT token is valid and if true gives access to the api handlers
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("Authorization") // Grab the token from the user request
		header = strings.TrimSpace(header)

		if header == "" {
			// Token is missing. Return with error code 403 not authorized
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{ResponseCode: http.StatusNotFound, Message: "Missing authentication token"})
			return
		}

		// Token usually come in the format
		// Bearer <toke>
		// so we have to parse the header into the bearer part and the token part
		splitHeader := strings.Split(header, " ")
		if len(splitHeader) != 2 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "Invalid/Malformed auth token"})
			return
		}

		tokenPart := splitHeader[1]
		tk := &models.Token{}

		_, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
