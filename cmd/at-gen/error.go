package main

type ArgumentError struct {
	Cause error
	Issue string
}

func (u ArgumentError) Error() string  { return u.Issue }
func (u ArgumentError) String() string { return u.Error() }
