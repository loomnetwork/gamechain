package view

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
)

var v View
var LayoutDir string = "view/layout"

type View struct {
	Template *template.Template
	Layout   string
}

type ViewData struct {
	Data    interface{}
	BaseUrl string
}

func getPage(page string) string {
	var buffer bytes.Buffer
	layout := packr.New("layout", "./layout")
	layoutFiles := layout.List()
	pages := packr.New("pages", "./pages")

	for i := 0; i < len(layoutFiles); i++ {
		b, err := layout.FindString(layoutFiles[i])
		if err != nil {
			panic(err)
		}
		buffer.WriteString(b)
	}
	b, err := pages.FindString(page)
	if err != nil {
		panic(err)
	}
	buffer.WriteString(b)
	return buffer.String()
}

func NewView(layout string, file string) *View {

	page := getPage(file)
	t, err := template.New(layout).Parse(page)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "/*.html")
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
	return files
}

func (v *View) Render(w http.ResponseWriter, data ...interface{}) error {
	vd := ViewData{
		Data: data,
	}
	return v.Template.ExecuteTemplate(w, v.Layout, vd)
}
