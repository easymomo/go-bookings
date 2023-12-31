package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/easymomo/go-bookings/internal/config"
	"github.com/easymomo/go-bookings/internal/models"
	"github.com/justinas/nosurf"
)

// functions is a variable that holds all the functions we want to make available to our go templates
var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, gohtml string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	var err error

	// Get the requested template from cache
	t, ok := tc[gohtml]
	if !ok {
		// log.Println("Error, could not get the template from template cache")
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err = t.Execute(buf, td)
	if err != nil {
		fmt.Print("Error while executing the template")
		log.Println(err)
		return err
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Print("Error while writing to the IO")
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// Get all the files named *.page.gohtml from ./template
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		fmt.Print("Error while getting the *.page.gohtml from ./template")
		return myCache, err
	}

	// range through all the files ending with *.page.gohtml
	for _, page := range pages {
		// Extract the name from the full path
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			fmt.Print("Error while parsing the pages")
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			fmt.Print("Error while searching for the layout files with filepath.Glob")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				fmt.Print("Error while parsing the layout files with filepath.Glob")
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
