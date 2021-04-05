package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"picapp/context"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// NewView sets the template for shared components to be used
// on each .gohtml file specified within the function's parameters.
// The templatized files will be appended to the parameterized 'files' list.
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r,nil)
}

// Render is used to render a view with a predefined layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	//do nothing
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		http.Error(w,"Something went wrong. Please contact me if this persists:"+
			"adamwoolhether@gmail.com",http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// layoutFiles returns a slice of strings representing
// all layout files used in PicApp
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes a slice of strings representing templates' file paths.
// It prepends TemplateDir dir to each of the slice's string
// Eg input {"home"} results in output {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTempalteExt takes a slice of strings representing templates' file paths.
// It appends template ext to each string in the slice.
// Eg. input {"home"} results in output {"home.gohtml"} if TemplateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
