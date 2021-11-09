package events

func (e *Event) validateProperties() error {
	return Validate(e.event.JSONSchema, *e.params.Properties)
}
