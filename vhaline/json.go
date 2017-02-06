package vhaline

import (
	"bytes"
	"encoding/json"
)

func prettyPrintJson(input []byte) []byte {
	var prettyBB bytes.Buffer
	jsErr := json.Indent(&prettyBB, input, "      ", "    ")
	if jsErr != nil {
		return input
	} else {
		return prettyBB.Bytes()
	}
}
