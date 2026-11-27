package ruleset

// PatternError is an error in one pattern of a rule. Its message is the message of the underlying
// error, so wrapping an error in it does not change what is reported.
type PatternError struct {
	RuleID       string
	PatternIndex int
	// IsFilter is true when the error is in the pattern's filters rather than the pattern itself
	IsFilter bool
	Err      error
}

func (err *PatternError) Error() string {
	return err.Err.Error()
}

func (err *PatternError) Unwrap() error {
	return err.Err
}
