package debugging

import (
	"attractify.io/platform/actions"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
)

type ActionSimulation struct {
	User   *requests.ActionSimulationUser
	Action *actions.Action
}

func (sim ActionSimulation) GetSteps() []responses.Step {
	steps := []responses.Step{}

	// whats the best way to get access to these Functions?
	// actions.Action has private atributes. Might make these public

	/*
		if sim.Action.Action.State == db.StateInactive || sim.Action.State == "" {
			return false
		}

		// State == staging with testusers and matching channel
		if sim.Action.State == db.StateStaging {
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
		}*/

	return steps
}
