package main

import (
	"encoding/json"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

func main() {
	var resp map[string]interface{}
	in, err := os.Open("codegen/user.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(in)
	b, err := io.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Name   string
		Fields map[string]interface{}
	}{
		"User",
		resp,
	}

	tpl, _ := template.New("template.tpl").Funcs(template.FuncMap{
		"Title": func(word string) string {
			words := strings.Split(word, "_")
			for idx, w := range words {
				words[idx] = convertWord(w)
			}
			return strings.Join(words, "")
		},
		"TypeOf": func(v interface{}) string {
			if v == nil {
				return "string"
			}
			return strings.ToLower(reflect.TypeOf(v).String())
		},
	}).ParseFiles("codegen/template.tpl")

	out, _ := os.Create("codegen/out.gen.go")
	defer out.Close()

	tpl.Execute(out, data)
}

func convertWord(w string) string {
	switch w {
	case "id":
		return "ID"
	}
	return cases.Title(language.English, cases.Compact).String(w)
}
