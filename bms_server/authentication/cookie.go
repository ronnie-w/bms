package authentication

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/ronnie-w/bms-server/schema"
	"github.com/ronnie-w/bms-server/utils"
)

var (
	cookieHandler     *securecookie.SecureCookie
	cookieHandlerName = base64.StdEncoding.EncodeToString([]byte(`hd3J>9I!6QxGv(&+%|BtZ""~Ip8],[bms]E0g_e2:9KK"Tuqj?m5uJ<:<>(B})TM`))
)

func init() {
	gen64 := securecookie.GenerateRandomKey(64)
	gen32 := securecookie.GenerateRandomKey(32)
	cookieHandler = securecookie.New(gen64, gen32)
}

func CreateCookie(rw http.ResponseWriter, name, value string) {
	var (
		base64enc, err = cookieHandler.Encode(cookieHandlerName, value)
		//expiry         = time.Now().Add(1 * time.Hour)
		cookie = &http.Cookie{
			Name:  name,
			Value: base64enc,
			//Expires:  expiry,
			SameSite: http.SameSiteLaxMode,
		}
	)

	if err != nil {
		log.Println("Error encoding cookie: ", err)
	}

	http.SetCookie(rw, cookie)
}

func ReadCookie(r *http.Request, name string) string {
	var (
		cookie, cookieErr = r.Cookie(name)
		cookieValue       string
	)

	if cookieErr == nil {
		if err := cookieHandler.Decode(cookieHandlerName, cookie.Value, &cookieValue); err != nil {
			log.Println("Error decoding cookie: ", err)
		}
	}

	return cookieValue
}

func Authenticate(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	var dbUUID, dbName string
	return func(rw http.ResponseWriter, r *http.Request) {
		authDetails := strings.Split(ReadCookie(r, "uuid"), "***")
		uuid := authDetails[1]
		if err := schema.DevSelect([]string{"uuid", "name"}, "admin", "uuid", uuid, []interface{}{&dbUUID, &dbName}); err != nil {
			rw.WriteHeader(http.StatusTeapot)
			rw.Write([]byte("Not enough permissions"))
			return
		}

		fn(rw, r)
	}
}

func FetchActiveUser(rw http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(rw, ReadCookie(r, "uuid"))
}
