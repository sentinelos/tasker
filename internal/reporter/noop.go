package reporter

// Noop is a reporter that does nothing.
type Noop struct{}

// SetStatus is a noop func.
func (n *Noop) SetStatus(state, policy, check, message string) error {
	return nil
}
