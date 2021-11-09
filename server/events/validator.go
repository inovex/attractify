package events

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func Validate(schema, document json.RawMessage) error {
	schemaLoader := gojsonschema.NewStringLoader(string(schema))
	documentLoader := gojsonschema.NewStringLoader(string(document))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}

	var e []string
	for _, i := range result.Errors() {
		e = append(e, i.String())
	}

	return fmt.Errorf("failed to validate document: %s", strings.Join(e, ", "))
}
