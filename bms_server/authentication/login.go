package authentication

import (
	"fmt"
	"net/http"

	"github.com/ronnie-w/bms-server/database"
	"github.com/ronnie-w/bms-server/schema"
	"github.com/ronnie-w/bms-server/utils"
	"golang.org/x/crypto/bcrypt"
)

var db = database.Conn()

func Login(rw http.ResponseWriter, r *http.Request) {
	var (
		uuid           string
		hashedPassword string
	)

	loginData := schema.GeneralDecoder(r)
	usernameOrId := loginData.Args2[0].Arg2Val1
	password := loginData.Args2[1].Arg2Val1

	authenticate := func(table string) {
		defaultSelect := func(clauseId string) error {
			return schema.DevSelect([]string{"uuid", "password"}, table, clauseId, usernameOrId, []interface{}{&uuid, &hashedPassword})
		}

		defaultAuth := func() error {
			if err := defaultSelect("name"); err != nil {
				if err := defaultSelect("id"); err != nil {
					utils.ErrResponse(rw, nil, "Invalid name or work id")
					return err
				}
			}

			return nil
		}

		switch table {
		case "admin":
			if err := defaultAuth(); err != nil {
				fmt.Println(err)
			} else if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
				utils.ErrResponse(rw, nil, "Incorrect password")
				fmt.Println(err)
			} else {
				CreateCookie(rw, "uuid", usernameOrId+"***"+uuid)
				utils.SuccessResponse(rw, "Login successful: "+usernameOrId+"")
			}
		case "staff":
			if err := defaultAuth(); err != nil {
				fmt.Println(err)
			} else if hashedPassword != password {
				utils.ErrResponse(rw, nil, "Invalid password")
				fmt.Println(err)
			} else {
				CreateCookie(rw, "uuid", usernameOrId+"***"+uuid)
				utils.SuccessResponse(rw, "Login successful: "+usernameOrId+"")
			}
		}
	}

	if usernameOrId == "admin" {
		authenticate("admin")
	} else {
		authenticate("staff")
	}
}
