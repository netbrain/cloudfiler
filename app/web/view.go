package web

import (
	. "github.com/netbrain/cloudfiler/app/conf"
	"github.com/netbrain/cloudfiler/app/web/fvgo"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	templatesDir = "/ui/web/view/"
)

var FormValidator = fvgo.NewFormValidator()
var Views = make(map[string]*template.Template)

func init() {
	loadTemplates()
}

/*
Loads all template files as templates and caches them, as well as parsing
any form definitions and creating validation rules from them.
*/
func loadTemplates() {
	log.Println("Parsing templates...")
	filepath.Walk(ViewDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), ".html") {
				log.Println(path)

				tmpl, err := template.ParseFiles(path)

				if err != nil {
					panic(err)
				}
				//TODO add a proper funcmap
				tmpl.Funcs(template.FuncMap{})

				FormValidator.ParseTemplate(path)

				relPath, _ := filepath.Rel(ViewDir(), path)
				view := relPath[0:strings.LastIndex(relPath, ".html")]
				Views[view] = tmpl

				log.Printf("Created view: %s", view)

			}
		}
		return nil
	})
}

func ViewFilePath(view string) string {
	return filepath.Join(ViewDir(), view+".html")
}

func ViewDir() string {
	return filepath.Join(Config.ApplicationHome, templatesDir)
}

func ViewExists(view string) bool {
	if _, exist := Views[view]; exist {
		return true
	}

	return false
}

func RenderView(view string, w io.Writer, data interface{}) {
	var tmpl *template.Template
	var exist bool

	if tmpl, exist = Views[view]; !exist {
		panic("view doesn't exist, create it first! " + view)
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
