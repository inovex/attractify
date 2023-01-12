package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

type Action struct {
	Ctx             context.Context
	App             *app.App
	Action          *db.Action
	Profile         *db.Profile
	profileIdentity *db.ProfileIdentity
	organizationID  uuid.UUID
	tags            []string
	properties      []db.ActionProperty
	Targeting       *db.ActionTargeting
	Capping         []db.ActionCapping
	hooks           []db.ActionHook
	TestUsers       []db.ActionTestUser
}

func New(ctx context.Context, app *app.App, orgID uuid.UUID, action *db.Action, profile *db.Profile, profileIdentity *db.ProfileIdentity) *Action {
	a := Action{
		Ctx:             ctx,
		App:             app,
		organizationID:  orgID,
		Action:          action,
		Profile:         profile,
		profileIdentity: profileIdentity,
	}

	json.Unmarshal(a.Action.Tags, &a.tags)
	json.Unmarshal(a.Action.Properties, &a.properties)
	json.Unmarshal(a.Action.Targeting, &a.Targeting)
	json.Unmarshal(a.Action.Capping, &a.Capping)
	json.Unmarshal(a.Action.TestUsers, &a.TestUsers)

	return &a
}

func (a *Action) ShouldDisplay(actionType string, tags []string, channel string, userID string, context json.RawMessage, time time.Time, timezone string) bool {
	// State
	if a.Action.State == db.StateInactive || a.Action.State == "" {
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
	if a.Action.State == db.StateStaging {
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
	if !a.CustomTraitConditions() {
		return false
	}

	if !a.ComputedTraitConditions() {
		return false
	}

	// Context conditions
	if !a.ContextConditions(channel, context) {
		return false
	}

	// Is in audience
	if len(a.Targeting.Audiences) > 0 && !a.IsInAudiences() {
		return false
	}

	// Capping
	if len(a.Capping) > 0 && !a.HasNoCapping() {
		return false
	}

	return true
}

func (a *Action) IsAllowedToAccept(channel string, userID string, time time.Time, timezone string) error {
	// State
	if a.Action.State == db.StateInactive || a.Action.State == "" {
		return errors.New("action inactive")
	}

	// Channel
	if !a.HasChannel(channel) {
		return errors.New("wrong channel")
	}

	// State == staging with testusers and matching channel
	if a.Action.State == db.StateStaging && !a.HasTestUser(userID, channel) {
		return errors.New("user is not a testuser")
	}

	// Time range
	if !a.InTimeRange(time, timezone) {
		return errors.New("wrong timerange")
	}

	// Is in audience
	if len(a.Targeting.Audiences) > 0 && !a.IsInAudiences() {
		return errors.New("user not in audience")
	}

	// Trait conditions
	if !a.CustomTraitConditions() {
		return errors.New("wrong customtraits")
	}

	if !a.ComputedTraitConditions() {
		return errors.New("wrong computedtraits")
	}

	// Capping
	fmt.Println(a.Action.ID)
	fmt.Println(a.Capping)
	fmt.Println(userID)
	if len(a.Capping) > 0 && !a.HasNoAcceptCapping() {
		fmt.Println("Error")
		return errors.New("capping")
	}

	return nil
}

func (a Action) HasTestUser(userID, channel string) bool {
	var testUser *db.ActionTestUser
	for _, t := range a.TestUsers {
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
	var testUser db.ActionTestUser
	var foundTestUser bool
	for _, t := range a.TestUsers {
		if t.UserID == userID {
			for _, c := range t.Channels {
				if c == channel {
					testUser = t
					foundTestUser = true
					break
				}
			}
		}
	}

	if !foundTestUser {
		return false
	}

	return testUser.SkipTargeting
}

func (a Action) HasNoCapping() bool {
	for _, c := range a.Capping {
		args := analytics.GetReactionCountParams{
			OrganizationID: a.Action.OrganizationID,
			ActionID:       a.Action.ID,
			Channels:       c.Channels,
			Event:          c.Event,
			IsUser:         c.Group == db.GroupUser,
			ProfileID:      a.Profile.ID,
			Within:         c.Within,
		}
		count, err := a.App.Analytics.GetReactionCount(args)
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
	for _, c := range a.Capping {
		fmt.Println("Capping test")
		fmt.Println(c.Event)
		if c.Event != analytics.ReactionEventAccepted {
			continue
		}
		args := analytics.GetReactionCountParams{
			OrganizationID: a.Action.OrganizationID,
			ActionID:       a.Action.ID,
			Channels:       c.Channels,
			Event:          c.Event,
			IsUser:         c.Group == db.GroupUser,
			ProfileID:      a.Profile.ID,
			Within:         c.Within,
		}
		count, err := a.App.Analytics.GetReactionCount(args)
		if err != nil {
			fmt.Println("Capping test - err")
			return false
		}
		if count >= c.Count {
			fmt.Println("Capping test - count")
			return false
		}
	}
	return true
}

func (a Action) HasChannel(channel string) bool {
	for _, c := range a.Targeting.Channels {
		if c == channel {
			return true
		}
	}
	return false
}

func (a Action) HasAndMatchesType(actionType string) bool {
	if len(actionType) > 0 {
		return a.Action.TypeName == actionType
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
	if a.Targeting.Start.Date != nil && a.Targeting.Start.Time == nil {
		state = 1
	}

	if a.Targeting.Start.Date == nil && a.Targeting.Start.Time != nil {
		state = 2
	}

	if a.Targeting.Start.Date != nil && a.Targeting.Start.Time != nil {
		state = 3
	}

	if state == 1 {
		sd, _ := time.ParseInLocation("2006-01-02", *a.Targeting.Start.Date, loc)
		sd = sd.UTC()
		if now.Before(sd) {
			return false
		}
	}

	if state == 2 {
		st, _ := time.ParseInLocation("15:04", *a.Targeting.Start.Time, zone)
		st = st.UTC()
		if now.Hour() < st.Hour() || now.Hour() == st.Hour() && now.Minute() < st.Minute() {
			return false
		}
	}

	if state == 3 {
		ts, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", *a.Targeting.Start.Date, *a.Targeting.Start.Time), loc)
		ts = ts.UTC()
		if now.Before(ts) {
			return false
		}
	}

	state = 0
	if a.Targeting.End.Date != nil && a.Targeting.End.Time == nil {
		state = 1
	}

	if a.Targeting.End.Date == nil && a.Targeting.End.Time != nil {
		state = 2
	}

	if a.Targeting.End.Date != nil && a.Targeting.End.Time != nil {
		state = 3
	}

	if state == 1 {
		ed, _ := time.ParseInLocation("2006-01-02", *a.Targeting.End.Date, loc)
		ed = ed.UTC()
		if now.After(ed) {
			return false
		}
	}

	if state == 2 {
		et, _ := time.ParseInLocation("15:04", *a.Targeting.End.Time, zone)
		et = et.UTC()
		if now.Hour() > et.Hour() || now.Hour() == et.Hour() && now.Minute() >= et.Minute() {
			return false
		}
	}

	if state == 3 {
		ts, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", *a.Targeting.End.Date, *a.Targeting.End.Time), loc)
		ts = ts.UTC()
		if now.After(ts) {
			return false
		}
	}

	return true
}

func (a Action) ComputedTraitConditions() bool {
	for _, c := range a.Targeting.TraitConditions {
		tc := targetingCondition{
			Key:      c.Key,
			Operator: c.Operator,
			Type:     c.Type,
			Value:    c.Value,
		}
		if c.Source != db.TraitConditionTypeCustom {
			if !tc.eval(a.Profile.ComputedTraits) {
				return false
			}
		}
	}
	return true
}

func (a Action) CustomTraitConditions() bool {
	for _, c := range a.Targeting.TraitConditions {
		tc := targetingCondition{
			Key:      c.Key,
			Operator: c.Operator,
			Type:     c.Type,
			Value:    c.Value,
		}
		if c.Source == db.TraitConditionTypeCustom {
			if !tc.eval(a.Profile.CustomTraits) {
				return false
			}
		}
	}
	return true
}

func (a Action) ContextConditions(channel string, context json.RawMessage) bool {
	for _, c := range a.Targeting.ContextConditions {
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
	ap, err := a.App.DB.GetAudienceProfile(
		a.Ctx,
		a.organizationID,
		a.Profile.ID,
		a.Targeting.Audiences,
	)
	if err != nil || ap == nil {
		return false
	}

	return true
}
