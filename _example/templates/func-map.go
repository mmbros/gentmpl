package templates

import (
	"html/template"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"weekday": formatItWeekday,
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
}

// formatItWeekday is an exemple of custom template func
//
//lint:ignore U1000 Ignore unused
func formatItWeekday(t time.Time) string {
	s := [...]string{"dom", "lun", "mar", "mer", "gio", "ven", "sab"}
	return s[t.Weekday()]
}
