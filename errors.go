package claw402

import "fmt"

// Error represents a non-2xx API response from claw402.
type Error struct {
	Status int
	Body   string
}

func (e *Error) Error() string {
	return fmt.Sprintf("claw402 API error %d: %s", e.Status, e.Body)
}
