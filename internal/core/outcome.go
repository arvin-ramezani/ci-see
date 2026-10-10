package core

// GateOutcome is a semantic result, distinct from the CLI's numeric exit code.
// ALLOW values are contract names only: S1 cannot produce an authorized allow.
type GateOutcome string

const (
	OutcomeAllowPass       GateOutcome = "ALLOW_PASS"
	OutcomeAllowBypass     GateOutcome = "ALLOW_BYPASS"
	OutcomeBlockFail       GateOutcome = "BLOCK_FAIL"
	OutcomeBlockIncomplete GateOutcome = "BLOCK_INCOMPLETE"
	OutcomeBlockCancelled  GateOutcome = "BLOCK_CANCELLED"
	OutcomeBlockError      GateOutcome = "BLOCK_ERROR"
)

func (o GateOutcome) Valid() bool {
	switch o {
	case OutcomeAllowPass, OutcomeAllowBypass, OutcomeBlockFail,
		OutcomeBlockIncomplete, OutcomeBlockCancelled, OutcomeBlockError:
		return true
	default:
		return false
	}
}
