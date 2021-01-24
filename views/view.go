package views

import "html/template"

// NewView sets the template for shared components to be used
// on each .gohtml file specified within the function's parameters.
// The templatized files will be appended to the parameterized 'files' list.
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
	}
}

type View struct {
	Template *template.Template
}
