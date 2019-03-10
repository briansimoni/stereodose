package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// JSON is a utility function that handlers can use to turn arbitrary structs
// into JSON and write the response
func JSON(w io.Writer, data interface{}) error {
	j, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}
	writer, ok := w.(http.ResponseWriter)
	if !ok {
		return errors.New("JSON write failed: writer was not an http.ResponseWriter")
	}
	writer.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(j))
	return nil
}
