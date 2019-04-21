package ppsync

// Lock is a generic interface for multiple
// Lock implementations (TAS, TTAS, A, EB).
type Lock interface {
	Lock()
	Unlock()
}
