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
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	if status > 400 && status < 600 {
		w.Write(js)
	}
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
