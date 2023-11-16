package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/easymomo/go-bookings/pkg/config"
	"github.com/easymomo/go-bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	var err error

	// Get the requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Error, could not get the template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err = t.Execute(buf, td)
	if err != nil {
		fmt.Print("Error while executing the template")
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Print("Error while writing to the IO")
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// Get all of the files named *.page.tmpl from ./template
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		fmt.Print("Error while getting the *.page.tmpl from ./template")
		return myCache, err
	}

	// range through all the files ending with *.page.tmpl
	for _, page := range pages {
		// Extract the name from the full path
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			fmt.Print("Error while parsing the pages")
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			fmt.Print("Error while searching for the layout files with filepath.Glob")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				fmt.Print("Error while parsing the layout files with filepath.Glob")
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
