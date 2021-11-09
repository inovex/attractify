package analytics

import (
	"github.com/gofrs/uuid"
)

type AudienceProfile struct {
	ID uuid.UUID `db:"id"`
}

func (a Analytics) RunAudience(query string) ([]AudienceProfile, error) {
	var res []AudienceProfile
	if err := a.DB.Select(&res, query); err != nil {
		return nil, err
	}
	return res, nil
}
