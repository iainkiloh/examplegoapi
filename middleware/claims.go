package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserInContext struct {
	UserEmail string
	UserId    string
	Roles     []Role
}

type Role struct {
	RoleId   int    `json:"RoleId"`
	RoleName string `json:"RoleName"`
}

//handler func gets custom claims from claims map and adds them to the context
//making them available to next handler in the pipeline
func GetClaimsContext(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//parse the Auth header to retrieve the token
		tokenHeader := r.Header.Get("Authorization")
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) < 2 {
			return
		}
		tokenPart := splitted[1]
		token, err := jwt.ParseWithClaims(tokenPart, jwt.MapClaims{}, tokenMiddleware.Options.ValidationKeyGetter)
		if err != nil {
			//fmt.Println(err)
			return
		}
		if token == nil {
			return
		}

		//cast Claims object to a map
		claimsMap := token.Claims.(jwt.MapClaims)

		//create userincontext and set fields using the map, which contains our custom claims
		userInContext := UserInContext{}
		userEmail, _ := claimsMap["KilohApp/UserEmail"]
		userId, _ := claimsMap["KilohApp/UserId"]

		if userEmail != nil {
			userInContext.UserEmail = userEmail.(string)
		}
		if userId != nil {
			userInContext.UserId = userId.(string)
		}

		//check for custom role claims
		roles, ok := claimsMap["KilohApp/SystemRole"].([]interface{})

		if ok {
			for _, v := range roles {
				var role Role
				jsonRole, ok := v.(string)
				if ok {
					err := json.Unmarshal([]byte(jsonRole), &role)
					if err == nil {
						userInContext.Roles = append(userInContext.Roles, role)
					}
				}
			}
		}

		//get the context and create a new one with our custom user type populate from our custom claims
		claimsContext := context.WithValue(r.Context(), "CustomUser", userInContext)

		//give the request our newly created context, with the userInContext type and key of 'CustomUser'
		//and pass it on to the next handler in pipeline
		next.ServeHTTP(w, r.WithContext(claimsContext))
	}

}
