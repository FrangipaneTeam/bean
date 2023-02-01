package template

import (
	"html/template"
	"log"
	"os"
)

// Process processes a template file and writes the output to a file.
func Process(rawTemplate string, outputFile string, data interface{}) {
	var err error

	// Read the template file
	t := template.Must(template.New("listTested").Parse(rawTemplate)) // .Funcs(funcMap)

	// Parse the template
	t, err = t.Parse(rawTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// create a new file
	file, _ := os.Create(outputFile)
	defer file.Close()

	// Execute the template
	err = t.ExecuteTemplate(file, "listTested", data)
	if err != nil {
		// TODO: handle error in app
		log.Print(err)
	}
}
