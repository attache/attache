package attache

import "fmt"

type BootstrapError struct {
	Cause error
	Phase string
}

func (b BootstrapError) Error() string {
	return fmt.Sprintf("bootstrap: %s: %s", b.Phase, b.Cause)
}
