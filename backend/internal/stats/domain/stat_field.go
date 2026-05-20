package domain

import "fmt"

// StatField enumerates the user-stats counters that can be incremented.
type StatField int

const (
	StatFieldUnknown StatField = iota
	StatFieldProjects
	StatFieldActivity
	StatFieldNotifications
)

// NewStatField parses the wire-level field name into a typed value.
// Unknown names are rejected.
func NewStatField(s string) (StatField, error) {
	switch s {
	case "projects":
		return StatFieldProjects, nil
	case "activity":
		return StatFieldActivity, nil
	case "notifications":
		return StatFieldNotifications, nil
	default:
		return StatFieldUnknown, fmt.Errorf("unknown stat field %q (want one of: projects, activity, notifications)", s)
	}
}

// String returns the wire-level name.
func (f StatField) String() string {
	switch f {
	case StatFieldProjects:
		return "projects"
	case StatFieldActivity:
		return "activity"
	case StatFieldNotifications:
		return "notifications"
	default:
		return "unknown"
	}
}
