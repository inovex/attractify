package events

import (
	"encoding/json"
	"strings"

	"attractify.io/platform/db"
)

type properties struct {
	Key        string `json:"key"`
	Type       string `json:"type"`
	Pattern    string `json:"pattern"`
	IsRequired bool   `json:"isRequired"`
}

type structure struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Children   []structure `json:"children"`
	Properties properties
}

type jsonSchema struct {
	Type                 string                 `json:"type,omitempty"`
	Properties           map[string]*jsonSchema `json:"properties,omitempty"`
	Items                []jsonSchema           `json:"items,omitempty"`
	Required             []string               `json:"required,omitempty"`
	Pattern              string                 `json:"pattern,omitempty"`
	AdditionalProperties bool                   `json:"additionalProperties"`
	AdditionalItems      bool                   `json:"additionalItems"`
}

type Definition struct {
	properties      []db.EventProperty
	parsedStructure []structure
}

func NewDefinition() *Definition {
	return &Definition{}
}

func (d *Definition) Prepare(structure json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	if err := d.parseStructure(structure); err != nil {
		return nil, nil, err
	}

	js, err := d.generateSchema()
	if err != nil {
		return nil, nil, err
	}

	j, err := json.Marshal(js)
	if err != nil {
		return nil, nil, err
	}

	d.properties = []db.EventProperty{}
	d.extractorFunc(d.parsedStructure, []string{})
	props, err := json.Marshal(d.properties)
	if err != nil {
		return nil, nil, err
	}

	return json.RawMessage(j), json.RawMessage(props), nil
}

func (d *Definition) parseStructure(structure json.RawMessage) error {
	if err := json.Unmarshal(structure, &d.parsedStructure); err != nil {
		return err
	}

	return nil
}

func (d *Definition) generateSchema() (*jsonSchema, error) {
	props := map[string]*jsonSchema{}
	required := d.generatorFunc(d.parsedStructure, props)

	js := jsonSchema{
		Type:                 "object",
		Properties:           props,
		Required:             required,
		AdditionalItems:      false,
		AdditionalProperties: false,
	}

	return &js, nil
}

func (d *Definition) generatorFunc(structure []structure, props map[string]*jsonSchema) []string {
	req := []string{}
	for _, v := range structure {
		if len(v.Children) > 0 {
			js := &jsonSchema{
				Type:       "object",
				Properties: map[string]*jsonSchema{},
			}
			props[v.Properties.Key] = js
			js.Required = d.generatorFunc(v.Children, props[v.Properties.Key].Properties)
		} else {
			t := v.Properties.Type
			switch v.Properties.Type {
			case "float":
				t = "number"
			case "dateTime":
				t = "string"
			}
			props[v.Properties.Key] = &jsonSchema{
				Type:    t,
				Pattern: v.Properties.Pattern,
			}
		}

		if v.Properties.IsRequired {
			req = append(req, v.Properties.Key)
		}
	}

	return req
}

func (d *Definition) extractorFunc(structure []structure, path []string) {
	for _, v := range structure {
		if len(v.Children) > 0 {
			d.extractorFunc(v.Children, append(path, v.Properties.Key))
		} else {
			d.properties = append(d.properties, db.EventProperty{
				Key:  strings.Join(append(path, v.Properties.Key), "."),
				Type: v.Properties.Type,
			})
		}
	}
}
