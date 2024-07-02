package dialer

import (
	"bytes"
	"encoding/json"
	"log"
	"syscall/js"

	"github.com/google/go-querystring/query"
	"github.com/ronnie-w/bms-app/schema"
	fetch "marwan.io/wasm-fetch"
)

func DefaultDialer(endPoint string, dynamic bool) func(this js.Value, args []js.Value) interface{} {
	var route string
	jsFunc := func(this js.Value, args []js.Value) interface{} {
		if dynamic {
			route = endPoint + "/" + args[2].String()
		} else {
			route = endPoint
		}

		form, err := query.Values(schema.GeneralDecoder(args))
		if err != nil {
			log.Println(err)
		}
		formData, err := json.Marshal(form)
		if err != nil {
			log.Println(err)
		}

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			var (
				res     *fetch.Response
				resolve = args[0]
				reject  = args[1]
			)

			go func() {
				res, err = fetch.Fetch(schema.URL+route, &fetch.Opts{
					Method:      fetch.MethodPost,
					Mode:        fetch.ModeCORS,
					Credentials: fetch.CredentialsInclude,
					Headers: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					Body: bytes.NewReader(formData),
				})
				if err != nil {
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}

				responseConstructor := js.Global().Get("Response")
				response := responseConstructor.New(string(res.Body))

				resolve.Invoke(response)
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	}

	return jsFunc
}
