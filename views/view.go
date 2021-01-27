package views

import "html/template"

// NewView sets the template for shared components to be used
// on each .gohtml file specified within the function's parameters.
// The templatized files will be appended to the parameterized 'files' list.
func NewView(layout string, files ...string) *View {
	// always include the following appended layouts
	files = append(files,
		"views/layouts/footer.gohtml",
		"views/layouts/bootstrap.gohtml",
		)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout: layout,
	}
}

type View struct {
	Template *template.Template
	Layout	string
}
