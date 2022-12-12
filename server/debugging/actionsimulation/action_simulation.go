package debugging

import (
	"fmt"
	"time"

	"attractify.io/platform/actions"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
)

type ActionSimulation struct {
	User   *requests.ActionSimulationUser
	Action *actions.Action
}

func (sim ActionSimulation) GetSteps() ([]responses.Step, bool) {
	steps := []responses.Step{}

	step := responses.Step{
		Name:  "Action State",
		State: "correct",
		Info:  "",
	}

	if sim.Action.Action.State == db.StateInactive || sim.Action.Action.State == "" {
		step.Info = "Actionstate is inactive"
		step.State = "error"
	}

	steps = append(steps, step)

	// Channel
	step = responses.Step{
		Name:  "Channel",
		State: "correct",
		Info:  "",
	}

	if !sim.Action.HasChannel(sim.User.Channel) {
		step.Info = "Action is not active for the users Channel"
		step.State = "error"
	}

	steps = append(steps, step)

	// Test User
	step = responses.Step{
		Name:  "Test User",
		State: "correct",
		Info:  "",
	}

	// State == staging with testusers and matching channel
	if sim.Action.Action.State == db.StateStaging {
		if !sim.Action.HasTestUser(sim.User.UserID.String(), sim.User.Channel) {
			step.Info = "User is not a testuser of this staging action."
			step.State = "error"
		} else if sim.Action.SkipTargeting(sim.User.UserID.String(), sim.User.Channel) {
			step.Info = "User skips targeting."
		}
	} else {
		step.Info = "Testusers are disabled. Actionstate is not staging."
	}

	steps = append(steps, step)

	// Time range
	step = responses.Step{
		Name:  "Time Range",
		State: "correct",
		Info:  "",
	}

	organization, err := sim.Action.App.DB.GetOrganization(sim.Action.Ctx, sim.Action.Action.OrganizationID)
	if err != nil {
		fmt.Println(err)
		step.Info = "Internal server error! Could not get organizations timezone."
		step.State = "server_error"
	} else if !sim.Action.InTimeRange(time.Now(), organization.Timezone) {
		step.Info = "Time range does not match"
		step.State = "error"
	}

	steps = append(steps, step)

	// Trait conditions
	sim.Action.Profile.ComputedTraits = sim.User.ComputedTraits
	sim.Action.Profile.CustomTraits = sim.User.CustomTraits

	step = responses.Step{
		Name:  "Computedtrait Conditions",
		State: "correct",
		Info:  "",
	}

	customTraitCondition := sim.Action.CustomTraitConditions()

	computedTraitCondition := sim.Action.ComputedTraitConditions()

	if !customTraitCondition {
		step.Info = "Computedtrait conditions are not fulfilled"
		step.State = "error"
	}
	steps = append(steps, step)

	step = responses.Step{
		Name:  "Customtrait Conditions",
		State: "correct",
		Info:  "",
	}
	if !computedTraitCondition {
		step.Info = "Computedtrait conditions are not fulfilled"
		step.State = "error"
	}
	steps = append(steps, step)

	// Context conditions
	step = responses.Step{
		Name:  "Context Conditions",
		State: "correct",
		Info:  "",
	}
	if !sim.Action.ContextConditions(sim.User.Channel, sim.User.Context) {
		step.Info = "Context conditions are not fulfilled"
		step.State = "error"
	}

	steps = append(steps, step)

	// Is in audience
	step = responses.Step{
		Name:  "Audience",
		State: "correct",
		Info:  "",
	}
	if len(sim.Action.Targeting.Audiences) > 0 && !sim.Action.IsInAudiences() {
		step.Info = "User is not in a given audience"
		step.State = "error"
	}
	steps = append(steps, step)

	// Capping
	step = responses.Step{
		Name:  "Capping",
		State: "correct",
		Info:  "",
	}
	if len(sim.Action.Capping) > 0 && !sim.Action.HasNoCapping() {
		step.Info = "Capping blocks the action"
		step.State = "error"
	}
	steps = append(steps, step)

	display := true

	for _, step := range steps {
		if step.State == "error" {
			display = false
			break
		}
	}

	return steps, display
}
