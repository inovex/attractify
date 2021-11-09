package events

import (
	"database/sql"
	"encoding/json"
)

func (e *Event) validateContext() error {
	schema, err := e.getContextSchema()
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	return Validate(schema, *e.params.Context)
}

func (e *Event) getContextSchema() (json.RawMessage, error) {
	cs, err := e.app.DB.GetContextSchema(e.ctx, e.params.OrganizationID, e.params.Channel)
	if err != nil {
		return nil, err
	}
	return cs, nil
}
