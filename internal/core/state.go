package core

// ValidationState describes the result for a requested validation state.
// A state label alone is never sufficient evidence to allow a Git operation.
type ValidationState string

const (
	StatePass       ValidationState = "PASS"
	StateFail       ValidationState = "FAIL"
	StateIncomplete ValidationState = "INCOMPLETE"
	StateNotRun     ValidationState = "NOT_RUN"
	StateStale      ValidationState = "STALE"
)

func (s ValidationState) Valid() bool {
	switch s {
	case StatePass, StateFail, StateIncomplete, StateNotRun, StateStale:
		return true
	default:
		return false
	}
}
