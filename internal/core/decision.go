package core

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidState   = errors.New("unknown validation state")
	ErrUnverifiedPass = errors.New("trusted exact-state validation unavailable in S1")
)

// Decision has no externally writable outcome. There is no S1 authorization
// constructor, fingerprint-verification interface, or approval input.
type Decision struct {
	outcome GateOutcome
	reason  string
}

func (d Decision) Outcome() GateOutcome { return d.outcome }
func (d Decision) Reason() string       { return d.reason }

// Evaluate is deliberately block-only until later slices implement independently
// verified exact-state proof. Never interpret caller-supplied PASS as an allow.
func Evaluate(state ValidationState) (Decision, error) {
	switch state {
	case StatePass:
		return Decision{OutcomeBlockError, ErrUnverifiedPass.Error()}, ErrUnverifiedPass
	case StateFail:
		return Decision{OutcomeBlockFail, "validation failed"}, nil
	case StateIncomplete:
		return Decision{OutcomeBlockIncomplete, "validation incomplete"}, nil
	case StateNotRun:
		return Decision{OutcomeBlockIncomplete, "validation not run"}, nil
	case StateStale:
		return Decision{OutcomeBlockIncomplete, "validation state stale"}, nil
	default:
		return Decision{OutcomeBlockError, ErrInvalidState.Error()}, fmt.Errorf("%w: %q", ErrInvalidState, state)
	}
}
