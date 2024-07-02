package parser

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ronnie-w/bms-server/resources"
)

func Parser(tempPath string) func(http.ResponseWriter, *http.Request) {
	parser := func(rw http.ResponseWriter) *template.Template {
		parseTemplate, err := template.ParseFS(resources.Resources, "dist/"+tempPath+".html")
		if err != nil {
			fmt.Println(err)
		}

		err = parseTemplate.Execute(rw, nil)
		if err != nil {
			log.Println(err)
		}

		return parseTemplate
	}

	template := func(rw http.ResponseWriter, r *http.Request) {
		parser(rw)
		rw.Header().Set("Cache-Control", "max-age=604800")
	}

	return template
}
