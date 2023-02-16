package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')
	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	//w.Header().Set("Content-Type", "application/json")
	fmt.Println(status)
	dat := string(js[:])
	fmt.Fprint(os.Stdout, dat)
	return nil
}
func (app *application) background(fn func()) {
	// Launch a background goroutine.
	go func() {
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.Print(fmt.Errorf("%s", err), nil)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}
