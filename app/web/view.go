package web

import (
	. "github.com/netbrain/cloudfiler/app/conf"
	"github.com/netbrain/cloudfiler/app/web/fvgo"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	templatesDir = "/ui/web/view/"
)

var FormValidator = fvgo.NewFormValidator()
var Views = make(map[string]*template.Template)

func init() {
	//loadTemplates()
	go func() {
		for {
			FormValidator = fvgo.NewFormValidator()
			loadTemplates()
			time.Sleep(5000 * time.Millisecond)
		}
	}()
}

/*
Loads all template files as templates and caches them, as well as parsing
any form definitions and creating validation rules from them.
*/
func loadTemplates() {
	log.Println("Parsing templates in: " + ViewDir())
	filepath.Walk(ViewDir(), func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			if !strings.HasPrefix(info.Name(), "_") && strings.HasSuffix(info.Name(), ".html") {
				log.Println(path)

				relPath, _ := filepath.Rel(ViewDir(), path)
				view := relPath[0:strings.LastIndex(relPath, ".html")]

				baseSrc, err := ioutil.ReadFile(baseFilePath())
				if err != nil {
					panic(err)
				}

				tmplSrc, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}
				tmpl := parseTemplate(string(baseSrc), string(tmplSrc))
				Views[view] = tmpl

				FormValidator.ParseTemplate(path)
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

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		panic(err)
	}
}

func baseFilePath() string {
	return filepath.Join(ViewDir(), "_base.html")
}

func parseTemplate(src ...string) *template.Template {
	t := template.New("*").Funcs(templateFunctions)
	for _, s := range src {
		t = template.Must(t.Parse(s))
	}
	return t
}
