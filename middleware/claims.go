package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserInContext struct {
	UserDisplayName string
	UserId          string
}

//handler func gets custom claims from claims map and adds them to the context
//making them available to next handler in the pipeline
func GetClaimsContext(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//todo this code needs a cleanup and error checking
		//todo get roles for user
		//parse the Auth header
		tokenHeader := r.Header.Get("Authorization")
		splitted := strings.Split(tokenHeader, " ")
		tokenPart := splitted[1]
		token, err := jwt.ParseWithClaims(tokenPart, jwt.MapClaims{}, tokenMiddleware.Options.ValidationKeyGetter)
		if err != nil {
			fmt.Println(err)
			return
		}
		if token == nil {
			return
		}

		//cast Claims object to a map
		claimsMap := token.Claims.(jwt.MapClaims)

		//create userincontext and set fields using the map, which contains our custom claims
		userInContext := UserInContext{}
		userDisplayName, _ := claimsMap["http://MyCustomClaimType/UserDisplayName"]
		userId, _ := claimsMap["http://MyCustomClaimType/UserID"]

		if userDisplayName != nil {
			userInContext.UserDisplayName = userDisplayName.(string)
		}
		if userId != nil {
			userInContext.UserId = userId.(string)
		}

		//get the context and create a new one with our custom user type populate from our custom claims
		claimsContext := context.WithValue(r.Context(), "CustomUser", userInContext)

		//give the request our newly created context, with the userInContext type and key of 'CustomUser'
		//and pass it on to the next handler in pipeline
		next.ServeHTTP(w, r.WithContext(claimsContext))
	}

}
