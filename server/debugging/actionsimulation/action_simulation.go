package debugging

import (
	"attractify.io/platform/actions"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
)

type ActionSimulation struct {
	User   *requests.ActionSimulationUser
	Action *actions.Action
}

func (sim ActionSimulation) GetSteps() []responses.Step {
	steps := []responses.Step{}

	step := responses.Step{
		Name:      "ActionState",
		UserValue: string(sim.Action.Action.State),
		DataValue: string(sim.Action.Action.State),
		Error:     "",
	}

	if sim.Action.Action.State == db.StateInactive || sim.Action.Action.State == "" {
		step.Error = "Actionstate is inactive"
	}

	steps = append(steps, step)

	step = responses.Step{
		Name:      "ActionState",
		UserValue: string(sim.Action.Action.State),
		DataValue: string(sim.Action.Action.State),
		Error:     "",
	}

	// State == staging with testusers and matching channel
	if sim.Action.Action.State == db.StateStaging {
		if !sim.Action.HasTestUser(sim.User.UserID.String(), sim.User.Channel) {
			step.Error = "Actionstate is staging, User is not a testuser"
		}

		if sim.Action.SkipTargeting(sim.User.UserID.String(), sim.User.Channel) {
		}
	}

	// Time range
	if !sim.Action.InTimeRange(sim.User.Time, "UTC") { // TODO: get timezone
		//return false
	}

	// Trait conditions
	if !sim.Action.TraitConditions() {
		//return false
	}

	// Context conditions
	if !sim.Action.ContextConditions(sim.User.Channel, sim.User.Context) {
		//return false
	}

	// Is in audience
	if len(sim.Action.Targeting.Audiences) > 0 && !sim.Action.IsInAudiences() {
		//return false
	}

	// Capping
	if len(sim.Action.Capping) > 0 && !sim.Action.HasNoCapping() {
		//return false
	}

	return steps
}
