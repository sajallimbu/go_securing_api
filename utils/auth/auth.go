package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sajallimbu/go_securing_api/models"
)

//JwtVerify ... our middleware that checks if the JWT token is valid and if true gives access to the api handlers
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("Authorization") // Grab the token from the user request
		header = strings.TrimSpace(header)

		if header == "" {
			// Token is missing. Return with error code 403 not authorized
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&models.Response{Success: false, ResponseCode: http.StatusForbidden, Message: "Missing authentication token"})
			return
		}

		// Token usually come in the format
		// Bearer <toke>
		// so we have to parse the header into the bearer part and the token part
		splitHeader := strings.Split(header, " ")
		if len(splitHeader) != 2 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&models.Response{Success: false, ResponseCode: http.StatusForbidden, Message: "Invalid/Malformed authentication token"})
			return
		}

		tokenPart := splitHeader[1]
		tk := &models.Token{}

		secret := os.Getenv("jwtSecret")

		_, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&models.Response{Success: false, ResponseCode: http.StatusForbidden, Message: "Invalid authentication token"})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
