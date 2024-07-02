package authentication

import (
	"fmt"
	"net/http"

	"github.com/ronnie-w/bms-server/schema"
	"github.com/ronnie-w/bms-server/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(rw http.ResponseWriter, r *http.Request) {
	var values []interface{}

	signupCredentials := schema.GeneralDecoder(r)
	table := signupCredentials.Args1.Arg1Val

	for i, s := range signupCredentials.Args2 {
		if i == 2 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Arg2Val1), 10)
			if err != nil {
				utils.ErrResponse(rw, nil, "Could not sign you up at the moment, try again later")
				fmt.Println(err)
				break
			} else {
				values = append(values, string(hashedPassword))
			}
			continue
		}

		values = append(values, s.Arg2Val1)
	}

	insertRes := schema.InsertQuery(table, values)
	if insertRes.Err != nil {
		utils.ErrResponse(rw, nil, insertRes.Err)
	} else {
		utils.SuccessResponse(rw, "Sign up successful")
	}
}

func ResetPassword(rw http.ResponseWriter, r *http.Request) {
	var email string
	resetCredentials := schema.GeneralDecoder(r)
	schema.DevSelect([]string{"email"}, "admin", "name", "admin", []interface{}{&email})

	parsedEmail, newPassword := resetCredentials.Args2[0].Arg2Val1, resetCredentials.Args2[1].Arg2Val1
	if len(newPassword) < 4 {
		utils.ErrResponse(rw, nil, "Password is too short")
	} else if parsedEmail != email {
		utils.ErrResponse(rw, nil, "Email entered is incorrect")
	} else {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
		if err := schema.DevUpdate("admin", "email", []string{"password"}, []interface{}{string(hashedPassword), email}); err == nil {
			utils.SuccessResponse(rw, "Password changed successfully")
		}
	}
}
