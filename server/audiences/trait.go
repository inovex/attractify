package audiences

type traits []property

func (t traits) generate() []string {
	var conditions []string
	for _, tr := range t {
		if tr.Operator == "exists" || tr.Operator == "not_exists" {
			tr.Target = "exists"
		}
		if tr.Source == "custom" {
			conditions = append(conditions, tr.generateForTraits("i.custom_traits"))
		} else {
			conditions = append(conditions, tr.generateForTraits("i.computed_traits"))
		}
	}

	return conditions
}
