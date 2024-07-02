package schema

import (
	"encoding/json"
	"log"
	"strings"
	"syscall/js"
)

var URL = "http://localhost:8000"

func GeneralDecoder(args []js.Value) Args {
	var (
		decArgs Args
		arg1    Arg1

		arg2Arr []Arg2
		arg2    Arg2

		_   = json.NewDecoder(strings.NewReader(parser(args))).Decode(&arg1)
		dec = json.NewDecoder(strings.NewReader(args[1].String()))
	)

	dec.Token()
	for dec.More() {
		if err := dec.Decode(&arg2); err != nil {
			log.Println(err)
		}

		arg2.Arg2Val1 = strings.ReplaceAll(arg2.Arg2Val1, " ", "_")
		arg2.Arg2Val2 = strings.ReplaceAll(arg2.Arg2Val2, " ", "_")

		arg2Arr = append(arg2Arr, arg2)
	}
	dec.Token()

	decArgs = Args{arg1, arg2Arr}
	return decArgs
}

func parser(args []js.Value) string {
	args1Split := strings.Split(args[0].String(), " ")
	for i := 1; i < len(args1Split); i++ {
		args1Split[i] = strings.ToLower(args1Split[i])
	}

	return strings.Join(args1Split, "_")
}

func StrSanitize(this js.Value, args []js.Value) interface{} {
	strArr := strings.Split(args[0].String(), "_")
	for i := 0; i < len(strArr); i++ {
		strArr[i] = strings.ToUpper(string(strArr[i][0])) + strings.TrimPrefix(strArr[i], string(strArr[i][0]))
	}

	return strings.Join(strArr, " ")
}

func StrParse(this js.Value, args []js.Value) interface{} {
	return parser(args)
}

// general schema
type Arg2 struct {
	Arg2Val1 string
	Arg2Val2 string
}

type Arg1 struct {
	Arg1Val string
}

type Args struct {
	Args1 Arg1
	Args2 []Arg2
}
