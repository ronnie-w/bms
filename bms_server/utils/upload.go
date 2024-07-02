package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Uploader(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(500 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		ErrResponse(rw, err, "Could not upload file")
	}
	defer file.Close()

	destination := "uploads/" + strings.Split(handler.Filename, "-")[0]
	if err := os.MkdirAll(destination, 0750); err != nil {
		fmt.Println(err)
	}

	fBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	if err := os.WriteFile(destination+"/"+handler.Filename, fBytes, 0660); err != nil {
		ErrResponse(rw, err, "Could not upload file")
	}

	json.NewEncoder(rw).Encode("Upload success")
}
