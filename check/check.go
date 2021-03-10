package check


// Check if the connection etc is available before middleware
type Check interface {
	Check() error
}

