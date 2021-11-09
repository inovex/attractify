package audiences

import (
	"fmt"
	"strings"
)

func escape(v string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(v)
}

func quote(v string) string {
	return fmt.Sprintf("'%s'", v)
}

func sqlID(id string) string {
	return "e_" + strings.ReplaceAll(id, "-", "")
}
