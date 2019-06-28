package controllers

import (
	b64 "encoding/base64"
	"errors"
	"net/http"

	"github.com/iainkiloh/examplegoapi/configurations"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var personController PersonController
var healthController HealthController
var signingKey string
var decodedKey []byte

//next func handler for check Jwt is valid
//allows us to inject this into the route pipeline before the route is handled
func checkJwt(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := jwtMiddleware.CheckJWT(w, r)
		if err != nil {
			//if theres a token validation error then return and dont execute the next handler
			return
		} else {
			//token is fine, move to next handler
			next.ServeHTTP(w, r)
		}
	}
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		//set the token issuer and audience from configuration
		iss, aud := configurations.GetTokenIssuerAndAudience()
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("Invalid audience.")
		}
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}
		return []byte(decodedKey), nil
	},

	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodHS256,
})

//registers routes and sets jwt required variables
func Startup() {
	personController.registerRoutes()
	healthController.registerRoutes()
	signingKey = configurations.GetTokenSigningKey()
	var err error
	decodedKey, err = b64.StdEncoding.DecodeString(signingKey)
	if err != nil {
		panic(err)
	}
}

type MyClaimsType struct {
	*jwt.MapClaims
	UserId string `json:"http://hmm/hmm/UserID"`
}

/*
type Credentials struct {
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
*/

/*
var jwtHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	//fmt.Fprintf(w, "This is an authenticated request")
	//fmt.Fprintf(w, "Claim content:\n")
	fmt.Println("user info:", user)
	//for k, v := range user.(*jwt.Token).Claims {
	//	fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	//}
})
*/
