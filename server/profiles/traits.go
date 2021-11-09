package profiles

import (
	"encoding/json"

	"attractify.io/platform/db"
)

func (p Profile) validateTraits() error {
	schema, err := p.getTraitSchema()
	if err != nil {
		return err
	}

	return Validate(schema, *p.params.Traits)
}

func (p Profile) getTraitSchema() (json.RawMessage, error) {
	schema, err := p.app.DB.GetCustomTraitsSchema(p.ctx, p.params.OrganizationID)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func (p *Profile) prepareTraits() error {
	if err := json.Unmarshal(*p.params.Traits, &p.traits); err != nil {
		return err
	}

	return nil
}

func (p *Profile) mergeTraits(leader, follower json.RawMessage) (json.RawMessage, error) {
	lt, err := p.unmarshalTraits(leader)
	if err != nil {
		return nil, err
	}

	ft, err := p.unmarshalTraits(follower)
	if err != nil {
		return nil, err
	}

	for k, v := range lt {
		if _, ok := ft[k]; !ok {
			ft[k] = v
		}
	}
	return p.marshalTraits(lt)
}

func (p *Profile) unmarshalTraits(traits json.RawMessage) (db.Traits, error) {
	var t db.Traits
	return t, json.Unmarshal(traits, &t)
}

func (p *Profile) marshalTraits(traits db.Traits) (json.RawMessage, error) {
	return json.Marshal(traits)
}
