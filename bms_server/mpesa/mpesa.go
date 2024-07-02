package mpesa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ronnie-w/bms-server/database"
	"github.com/ronnie-w/bms-server/utils"
)

type Authorization struct {
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expires_in"`
}

type RegisterPayload struct {
	ShortCode       int
	ResponseType    string
	ConfirmationURL string
	ValidationURL   string
}

func ConfirmC2BMpesa(rw http.ResponseWriter, r *http.Request) {
	var (
		client    = &http.Client{}
		basicAuth = base64.StdEncoding.EncodeToString([]byte("BASIC AUTH"))

		accessToken Authorization
		mpesaTill   int
	)

	if err := database.Conn().QueryRow("SELECT mpesa_till_number FROM settings").Scan(&mpesaTill); err != nil {
		utils.ErrResponse(rw, nil, "Could not find till number")
		return
	}

	authTokenReq, err := http.NewRequest("GET", "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials", nil)
	if err != nil {
		fmt.Println(err)
	}
	authTokenReq.Header.Add("Content-Type", "application/json")
	authTokenReq.Header.Add("Authorization", "Bearer "+basicAuth+"")

	authTokenRes, err := client.Do(authTokenReq)
	if err != nil {
		fmt.Println(err)
	}
	defer authTokenRes.Body.Close()
	json.NewDecoder(authTokenRes.Body).Decode(&accessToken)

	payloadBytes, err := json.Marshal(RegisterPayload{
		ShortCode:       mpesaTill,
		ResponseType:    "Completed",
		ConfirmationURL: "https://webhook.site/627d718d-8393-407c-8ed8-e699f2fb94f9",//test url,
		ValidationURL:   "https://webhook.site/627d718d-8393-407c-8ed8-e699f2fb94f9",//test url,
	})
	if err != nil {
		fmt.Println(err)
	}

	confirmationReq, err := http.NewRequest("POST", "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/registerurl", strings.NewReader(string(payloadBytes)))
	if err != nil {
		fmt.Println(err)
	}
	confirmationReq.Header.Add("Content-Type", "application/json")
	authTokenReq.Header.Add("Authorization", "Bearer "+accessToken.AccessToken+"")

	confirmationRes, err := client.Do(confirmationReq)
	if err != nil {
		fmt.Println(err)
	}
	defer confirmationRes.Body.Close()

	utils.SuccessResponse(rw, "Transaction complete")
}
