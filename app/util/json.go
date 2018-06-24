package util

import (
	"encoding/json"
	"fmt"
	"io"
)

func JSON(w io.Writer, data interface{}) error {
	j, err := json.MarshalIndent(data, "    ", "    ")
	if err != nil {
		return err
	}
	fmt.Fprint(w, string(j))
	return nil
}
