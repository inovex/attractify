package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

type Action struct {
	ctx             context.Context
	app             *app.App
	action          *db.Action
	profile         *db.Profile
	profileIdentity *db.ProfileIdentity
	organizationID  uuid.UUID
	tags            []string
	properties      []db.ActionProperty
	targeting       *db.ActionTargeting
	capping         []db.ActionCapping
	hooks           []db.ActionHook
	testUsers       []db.ActionTestUser
}

func New(ctx context.Context, app *app.App, orgID uuid.UUID, action *db.Action, profile *db.Profile, profileIdentity *db.ProfileIdentity) *Action {
	a := Action{
		ctx:             ctx,
		app:             app,
		organizationID:  orgID,
		action:          action,
		profile:         profile,
		profileIdentity: profileIdentity,
	}

	json.Unmarshal(a.action.Tags, &a.tags)
	json.Unmarshal(a.action.Properties, &a.properties)
	json.Unmarshal(a.action.Targeting, &a.targeting)
	json.Unmarshal(a.action.Capping, &a.capping)
	json.Unmarshal(a.action.TestUsers, &a.testUsers)

	return &a
}

func (a *Action) ShouldDisplay(actionType uuid.UUID, tags []string, channel string, userID string, context json.RawMessage, time time.Time, timezone string) bool {
	// State
	if a.action.State == db.StateInactive || a.action.State == "" {
		return false
	}

	// Channel
	if !a.HasChannel(channel) {
		return false
	}

	// Type
	if !a.HasAndMatchesType(actionType) {
		return false
	}

	// Tags
	if len(tags) > 0 && !a.HasTags(tags) {
		return false
	}

	// State == staging with testusers and matching channel
	if a.action.State == db.StateStaging {
		if !a.HasTestUser(userID, channel) {
			return false
		}
		if a.SkipTargeting(userID, channel) {
			return true
		}
	}

	// Time range
	if !a.InTimeRange(time, timezone) {
		return false
	}

	// Trait conditions
	if !a.TraitConditions() {
		return false
	}

	// Context conditions
	if !a.ContextConditions(channel, context) {
		return false
	}

	// Is in audience
	if len(a.targeting.Audiences) > 0 && !a.IsInAudiences() {
		return false
	}

	// Capping
	if len(a.capping) > 0 && !a.HasNoCapping() {
		return false
	}

	return true
}

func (a *Action) IsAllowedToAccept(channel string, userID string, time time.Time, timezone string) bool {
	// State
	if a.action.State == db.StateInactive || a.action.State == "" {
		return false
	}

	// Channel
	if !a.HasChannel(channel) {
		return false
	}

	// State == staging with testusers and matching channel
	if a.action.State == db.StateStaging && !a.HasTestUser(userID, channel) {
		return false
	}

	// Time range
	if !a.InTimeRange(time, timezone) {
		return false
	}

	// Is in audience
	if len(a.targeting.Audiences) > 0 && !a.IsInAudiences() {
		return false
	}

	// Trait conditions
	if !a.TraitConditions() {
		return false
	}

	// Capping
	if len(a.capping) > 0 && !a.HasNoAcceptCapping() {
		return false
	}

	return true
}

func (a Action) HasTestUser(userID, channel string) bool {
	var testUser *db.ActionTestUser
	for _, t := range a.testUsers {
		if t.UserID == userID {
			for _, c := range t.Channels {
				if c == channel {
					testUser = &t
					break
				}
			}
		}
	}

	if testUser == nil {
		return false
	}

	return true
}

func (a Action) SkipTargeting(userID, channel string) bool {
	var testUser *db.ActionTestUser
	for _, t := range a.testUsers {
		if t.UserID == userID {
			for _, c := range t.Channels {
				if c == channel {
					testUser = &t
					break
				}
			}
		}
	}

	if testUser == nil {
		return false
	}

	return testUser.SkipTargeting
}

func (a Action) HasNoCapping() bool {
	for _, c := range a.capping {
		args := analytics.GetReactionCountParams{
			OrganizationID: a.action.OrganizationID,
			ActionID:       a.action.ID,
			Channels:       c.Channels,
			Event:          c.Event,
			IsUser:         c.Group == db.GroupUser,
			ProfileID:      a.profile.ID,
			Within:         c.Within,
		}
		count, err := a.app.Analytics.GetReactionCount(args)
		if err != nil {
			return false
		}
		if count >= c.Count {
			return false
		}
	}
	return true
}

func (a Action) HasNoAcceptCapping() bool {
	for _, c := range a.capping {
		if c.Event != analytics.ReactionEventAccepted {
			continue
		}
		args := analytics.GetReactionCountParams{
			OrganizationID: a.action.OrganizationID,
			ActionID:       a.action.ID,
			Channels:       c.Channels,
			Event:          c.Event,
			IsUser:         c.Group == db.GroupUser,
			ProfileID:      a.profile.ID,
			Within:         c.Within,
		}
		count, err := a.app.Analytics.GetReactionCount(args)
		if err != nil {
			return false
		}
		if count >= c.Count {
			return false
		}
	}
	return true
}

func (a Action) HasChannel(channel string) bool {
	for _, c := range a.targeting.Channels {
		if c == channel {
			return true
		}
	}
	return false
}

func (a Action) HasAndMatchesType(actionType uuid.UUID) bool {
	if len(actionType) > 0 {
		return a.action.Type == actionType
	}
	return true
}

func (a Action) HasTags(tags []string) bool {
	for _, tag := range tags {
		match := false
		for _, t := range a.tags {
			if tag == t {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}

func (a Action) InTimeRange(now time.Time, timezone string) bool {
	loc, _ := time.LoadLocation(timezone)
	now = now.UTC()

	// Use given time, convert it to the user's timezone to obtain the fixed zone
	name, offset := now.In(loc).Zone()
	zone := time.FixedZone(name, offset)

	state := 0
	if a.targeting.Start.Date != nil && a.targeting.Start.Time == nil {
		state = 1
	}

	if a.targeting.Start.Date == nil && a.targeting.Start.Time != nil {
		state = 2
	}

	if a.targeting.Start.Date != nil && a.targeting.Start.Time != nil {
		state = 3
	}

	if state == 1 {
		sd, _ := time.ParseInLocation("2006-01-02", *a.targeting.Start.Date, loc)
		sd = sd.UTC()
		if now.Before(sd) {
			return false
		}
	}

	if state == 2 {
		st, _ := time.ParseInLocation("15:04", *a.targeting.Start.Time, zone)
		st = st.UTC()
		if now.Hour() < st.Hour() || now.Hour() == st.Hour() && now.Minute() < st.Minute() {
			return false
		}
	}

	if state == 3 {
		ts, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", *a.targeting.Start.Date, *a.targeting.Start.Time), loc)
		ts = ts.UTC()
		if now.Before(ts) {
			return false
		}
	}

	state = 0
	if a.targeting.End.Date != nil && a.targeting.End.Time == nil {
		state = 1
	}

	if a.targeting.End.Date == nil && a.targeting.End.Time != nil {
		state = 2
	}

	if a.targeting.End.Date != nil && a.targeting.End.Time != nil {
		state = 3
	}

	if state == 1 {
		ed, _ := time.ParseInLocation("2006-01-02", *a.targeting.End.Date, loc)
		ed = ed.UTC()
		if now.After(ed) {
			return false
		}
	}

	if state == 2 {
		et, _ := time.ParseInLocation("15:04", *a.targeting.End.Time, zone)
		et = et.UTC()
		if now.Hour() > et.Hour() || now.Hour() == et.Hour() && now.Minute() >= et.Minute() {
			return false
		}
	}

	if state == 3 {
		ts, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", *a.targeting.End.Date, *a.targeting.End.Time), loc)
		ts = ts.UTC()
		if now.After(ts) {
			return false
		}
	}

	return true
}

func (a Action) TraitConditions() bool {
	for _, c := range a.targeting.TraitConditions {
		tc := targetingCondition{
			Key:      c.Key,
			Operator: c.Operator,
			Type:     c.Type,
			Value:    c.Value,
		}
		if c.Source == db.TraitConditionTypeCustom {
			if !tc.eval(a.profile.CustomTraits) {
				return false
			}
		} else {
			if !tc.eval(a.profile.ComputedTraits) {
				return false
			}
		}
	}
	return true
}

func (a Action) ContextConditions(channel string, context json.RawMessage) bool {
	for _, c := range a.targeting.ContextConditions {
		if channel != c.Channel {
			continue
		}
		tc := targetingCondition{
			Key:      c.Key,
			Operator: c.Operator,
			Type:     c.Type,
			Value:    c.Value,
		}
		if !tc.eval(context) {
			return false
		}
	}
	return true
}

func (a Action) IsInAudiences() bool {
	ap, err := a.app.DB.GetAudienceProfile(
		a.ctx,
		a.organizationID,
		a.profile.ID,
		a.targeting.Audiences,
	)
	if err != nil || ap == nil {
		return false
	}

	return true
}
