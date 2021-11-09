package profiles

import (
	"encoding/json"
	"errors"

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

	return errors.New("failed to validate document")
}
