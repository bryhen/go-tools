package lifecycle

// Passes the provided error to Run(), which signals that a critical runtime error has occurred and the application should exit.
//
// Only the first error passed to RuntimeErr() will be propogated. Safe to call concurrently. If called after shutdown has begun, error will be discarded.
func RuntimeErr(runtimeErr error) {
	select {
	case rte <- runtimeErr:
	default:
	}
}
