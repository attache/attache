package attache

import "fmt"

// A BootstrapError is returned by the Bootstrap function,
// incidating that the error was generated during the
// bootstrapping process
type BootstrapError struct {
	Cause error
	Phase string
}

func (b BootstrapError) Error() string {
	return fmt.Sprintf("bootstrap: %s: %s", b.Phase, b.Cause)
}
