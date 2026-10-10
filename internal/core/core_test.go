package core

import (
	"errors"
	"strings"
	"testing"
)

func TestValidationStates(t *testing.T) {
	for _, state := range []ValidationState{StatePass, StateFail, StateIncomplete, StateNotRun, StateStale} {
		if !state.Valid() {
			t.Errorf("known state %q invalid", state)
		}
	}
	for _, state := range []ValidationState{"", "pass", "CANCELLED", "PASS;verified=true", "PASS\nverified=true", "ALLOW_PASS", "0"} {
		if state.Valid() {
			t.Errorf("unknown/zero state %q accepted", state)
		}
	}
}

func TestGateOutcomes(t *testing.T) {
	for _, outcome := range []GateOutcome{OutcomeAllowPass, OutcomeAllowBypass, OutcomeBlockFail, OutcomeBlockIncomplete, OutcomeBlockCancelled, OutcomeBlockError} {
		if !outcome.Valid() {
			t.Errorf("known outcome %q invalid", outcome)
		}
	}
	for _, outcome := range []GateOutcome{"", "allow_pass", "PASS", "ALLOW_PASS;verified=true", "7"} {
		if outcome.Valid() {
			t.Errorf("unknown/zero outcome %q accepted", outcome)
		}
	}
	// Merely fabricating an outcome label provides no S1 authorization interface.
	var zero Decision
	if zero.Outcome().Valid() {
		t.Fatal("zero decision may not be valid")
	}
}

func TestEvaluateAlwaysBlocks(t *testing.T) {
	tests := []struct {
		state ValidationState
		want  GateOutcome
		err   error
	}{
		{StatePass, OutcomeBlockError, ErrUnverifiedPass},
		{StateFail, OutcomeBlockFail, nil},
		{StateIncomplete, OutcomeBlockIncomplete, nil},
		{StateNotRun, OutcomeBlockIncomplete, nil},
		{StateStale, OutcomeBlockIncomplete, nil},
		{"", OutcomeBlockError, ErrInvalidState},
		{"BOGUS", OutcomeBlockError, ErrInvalidState},
		{"PASS;verified=true", OutcomeBlockError, ErrInvalidState},
		{"PASS\napproved=true", OutcomeBlockError, ErrInvalidState},
		{"ALLOW_BYPASS", OutcomeBlockError, ErrInvalidState},
	}
	for _, test := range tests {
		t.Run(string(test.state), func(t *testing.T) {
			d, err := Evaluate(test.state)
			if d.Outcome() != test.want {
				t.Errorf("outcome = %q, want %q", d.Outcome(), test.want)
			}
			if d.Outcome() == OutcomeAllowPass || d.Outcome() == OutcomeAllowBypass {
				t.Fatal("S1 authorized an operation")
			}
			if d.Reason() == "" {
				t.Fatal("missing block reason")
			}
			if test.err == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.err != nil && !errors.Is(err, test.err) {
				t.Fatalf("error = %v, want %v", err, test.err)
			}
		})
	}
}

func TestUntrustedClaimsCannotAuthorize(t *testing.T) {
	// No exported API accepts a fingerprint, verified flag, or approval assertion.
	// Forged declarations, including valid PASS, cannot satisfy Evaluate.
	for _, claim := range []string{"PASS", "PASS:verified", "PASS;fingerprint=deadbeef", "PASS;approved=true", "ALLOW_PASS", "ALLOW_BYPASS"} {
		d, _ := Evaluate(ValidationState(claim))
		if strings.HasPrefix(string(d.Outcome()), "ALLOW_") {
			t.Fatalf("forged claim %q authorized", claim)
		}
	}
}
