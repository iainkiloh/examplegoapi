package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
)

var tokenMiddleware *jwtmiddleware.JWTMiddleware

func SetJwtMiddleware(middleware *jwtmiddleware.JWTMiddleware) {
	tokenMiddleware = middleware
}

//next func handler for check Jwt is valid
//allows us to inject this into the route pipeline before the route is handled
func CheckJwt(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tokenMiddleware.CheckJWT(w, r)
		if err != nil {
			//if theres a token validation error then return and dont execute the next handler
			return
		} else {
			//token is fine, move to next handler
			next.ServeHTTP(w, r)
		}
	}
}
