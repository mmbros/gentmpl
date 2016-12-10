package main

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

func formatItWeekday(t time.Time) string {
	s := [...]string{"dom", "lun", "mar", "mer", "gio", "ven", "sab"}
	return s[t.Weekday()]
}
