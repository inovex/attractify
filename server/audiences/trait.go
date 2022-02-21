package audiences

import (
	"strings"
)

type traits []property

func (t traits) generate() string {
	var tl []string
	for _, tr := range t {
		if tr.Source == "custom" {
			tl = append(tl, tr.generateForTraits("i.custom_traits"))
		} else {
			tl = append(tl, tr.generateForTraits("i.computed_traits"))
		}
	}

	if len(tl) == 0 {
		return ""
	}

	return strings.Join(tl, "\nAND ")
}
