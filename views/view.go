package views

import (
	"html/template"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
)

// NewView sets the template for shared components to be used
// on each .gohtml file specified within the function's parameters.
// The templatized files will be appended to the parameterized 'files' list.
func NewView(layout string, files ...string) *View {
	// always include the following appended layouts
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

// layoutFiles returns a slice of strings representing
// all layout files used in PicApp
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}
