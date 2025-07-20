package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"time"
)

type ExitReason struct {
	OsSignal     os.Signal
	StartupErr   error
	RuntimeErr   error
	ShutdownErrs []error
}

var rte = make(chan error, 1)

// Helps run an application by handling graceful startup and shutdown.
//
// Returns the ExitReason struct which contains information about why the program exited.
//
// This function will:
//
// 1. Run startup function sequentially up to the startupMaxDur duration.
//   - If any of these functions returns an error, this function will return immediately.
//
// 2. Block until an OS shutdown signal (ie ctrl+c) or runtime error (ie your code calls runtimeErr(yourError)) is received.
//   - Only the first runtime error received (if any) will be returned. All others are discarded.
//
// 3. Run the shutdown functions sequentially up to the shutdownMaxDur duration.
func Start(startupMaxDur time.Duration, shutdownMaxDur time.Duration, startupFns []func(ctx context.Context) error, shutdownFns []func(ctx context.Context) error) *ExitReason {
	er := &ExitReason{}
	fnErrs := make(chan error)

	// Startup the application and exit early if any errors occur.
	stCtx, stCancel := context.WithTimeout(context.Background(), startupMaxDur)

	go func() {

		for _, fn := range startupFns {
			if err := fn(stCtx); err != nil {
				fnErrs <- err
				stCancel()
				return
			}
		}

		fnErrs <- nil
	}()

	select {
	case er.StartupErr = <-fnErrs:
	case <-stCtx.Done():
		er.StartupErr = stCtx.Err()
	}

	stCancel()

	if er.StartupErr != nil {
		return er
	}

	// Monitor the application/OS and document why we're shutting down.
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, os.Interrupt)

	select {
	case er.RuntimeErr = <-rte:
	case er.OsSignal = <-osSig:
	}

	// Shutdown the application and collect all the errors that occurred during shutdown.
	sdCtx, sdCancel := context.WithTimeout(context.Background(), shutdownMaxDur)
	defer sdCancel()

	go func() {
		for _, fn := range shutdownFns {
			if err := fn(sdCtx); err != nil {
				fnErrs <- err
			}
		}
	}()

Shutdown:
	for range shutdownFns {
		select {
		case e := <-fnErrs:
			if e != nil {
				er.ShutdownErrs = append(er.ShutdownErrs, e)
			}
		case <-sdCtx.Done():
			er.ShutdownErrs = append(er.ShutdownErrs, sdCtx.Err())
			break Shutdown
		}
	}

	return er
}
