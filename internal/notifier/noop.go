package notifier

// SetStatus is a noop func.
func (n *Noop) SetStatus(state, policy, check, message string) error {
	return nil
}
