// Package reporter provides check result reporting.
package reporter

// Reporter describes a hook for sending summarized results to a remote API.
type Reporter interface {
	SetStatus(string, string, string, string) error
}
