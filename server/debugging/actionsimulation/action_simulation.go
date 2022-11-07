package debugging

import (
	"fmt"

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
		Name:      "Action State",
		UserValue: string(sim.Action.Action.State),
		DataValue: string(sim.Action.Action.State),
		Blocking:  false,
		Info:      "",
	}

	if sim.Action.Action.State == db.StateInactive || sim.Action.Action.State == "" {
		step.Info = "Actionstate is inactive"
		step.Blocking = true
	}

	steps = append(steps, step)

	step = responses.Step{
		Name:      "Test User",
		UserValue: "",
		DataValue: "",
		Blocking:  false,
		Info:      "",
	}

	// State == staging with testusers and matching channel
	if sim.Action.Action.State == db.StateStaging {
		if !sim.Action.HasTestUser(sim.User.UserID.String(), sim.User.Channel) {
			step.Info = "User is not a testuser of this staging action."
			step.Blocking = true
		} else if sim.Action.SkipTargeting(sim.User.UserID.String(), sim.User.Channel) {
			step.Info = "User skips targeting."
		}
	} else {
		step.Info = "Testuser are disabled. Actionstate is not staging."
	}

	steps = append(steps, step)

	// Time range
	step = responses.Step{
		Name:      "Time Range",
		UserValue: "",
		DataValue: "",
		Blocking:  false,
		Info:      "",
	}
	orgDB := db.New(sim.Action.App.Analytics.DB)
	organization, err := orgDB.GetOrganization(sim.Action.Ctx, sim.Action.Action.OrganizationID)
	if err != nil {
		step.Info = "Internal server error! Could not get organizations timezone."
	} else if !sim.Action.InTimeRange(sim.User.Time, organization.Timezone) { // TODO: get timezone
		step.Info = "Time range does not match"
		step.Blocking = true
		step.UserValue = sim.User.Time.UTC().String()
		step.DataValue = *sim.Action.Targeting.Start.Date + " " + *sim.Action.Targeting.Start.Time + " -> " + *sim.Action.Targeting.End.Date + " " + *sim.Action.Targeting.End.Time
	}

	steps = append(steps, step)

	// Trait conditions
	sim.Action.Profile.ComputedTraits = sim.User.ComputedTraits
	sim.Action.Profile.CustomTraits = sim.User.CustomTraits

	step = responses.Step{
		Name:      "Computedtrait Conditions",
		UserValue: "",
		DataValue: "",
		Blocking:  false,
		Info:      "",
	}

	customTraitCondition, tc := sim.Action.CustomTraitConditions()

	computedTraitCondition, tc := sim.Action.ComputedTraitConditions()

	if !customTraitCondition {
		step.Info = "Computedtrait conditions are not fulfilled"
		step.UserValue = string(sim.User.ComputedTraits)
		step.DataValue = tc.Key + " " + tc.Operator + " " + fmt.Sprintf("%v", tc.Value)
		step.Blocking = true
	}
	steps = append(steps, step)

	step = responses.Step{
		Name:      "Customtrait Conditions",
		UserValue: "",
		DataValue: "",
		Blocking:  false,
		Info:      "",
	}
	if !computedTraitCondition {
		step.Info = "Computedtrait conditions are not fulfilled"
		step.UserValue = string(sim.User.CustomTraits)
		step.DataValue = tc.Key + " " + tc.Operator + " " + fmt.Sprintf("%v", tc.Value)
		step.Blocking = true
	}
	steps = append(steps, step)

	// Context conditions
	step = responses.Step{
		Name:      "Context Conditions",
		UserValue: "",
		DataValue: "",
		Blocking:  false,
		Info:      "",
	}
	if !sim.Action.ContextConditions(sim.User.Channel, sim.User.Context) {
		step.UserValue = string(sim.User.Context)
		step.DataValue = "" //sim.Action.Action.Context TODO
		step.Info = "Context conditions are not fulfilled"
		step.Blocking = true
	}

	steps = append(steps, step)

	// Is in audience
	step = responses.Step{
		Name:      "Audience",
		UserValue: "",
		DataValue: "",
		Blocking:  false,
		Info:      "",
	}
	if len(sim.Action.Targeting.Audiences) > 0 && !sim.Action.IsInAudiences() {
		//return false
	}
	steps = append(steps, step)

	// Capping
	step = responses.Step{
		Name:      "Capping",
		UserValue: "",
		DataValue: "",
		Blocking:  false,
		Info:      "",
	}
	if len(sim.Action.Capping) > 0 && !sim.Action.HasNoCapping() {
		//return false
	}
	steps = append(steps, step)

	return steps
}
