package exsimple

/*
	This code can be copied and ran. It will start an HTTP server at localhost:8080.
	Requests at any path except /error will return 200 and echo the path.
	If the path is /error the server will report an error and the application will shut down.
	If the application is ended with ctrl+c, the server will be shut down.
*/

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bryhen/go-tools/lifecycle"
)

const (
	// How long to wait before startup should fail since it's taking too long. See advanced example for more details.
	startupMaxDur = time.Second * 30

	// How long to wait before shutdown should fail since it's taking too long. See advanced example for more details.
	shutdownMaxDur = time.Second * 300
)

var (
	// The functions to run during startup, in order. See advanced example for more details.
	startupFns = []func(context.Context) error{startServer}

	// The functions to run during shutdown, in order. See advanced example for more details.
	shutdownFns = []func(context.Context) error{stopServer}

	// EXAMPLE CODE: An http server
	server = &http.Server{
		Addr: ":8080",
		// EXAMPLE CODE: An echo or error handler.
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/error" {
				// If the path is /error, shut down the application
				fmt.Println("Error path hit. Calling lifecycle.RunetimeErr().")
				lifecycle.RuntimeErr(fmt.Errorf("/error path was hit"))
				w.Write([]byte("Server shutting down because /error was hit."))
			} else {
				// If it isn't /error, echo the path
				w.Write([]byte(r.URL.Path))
			}
		}),
	}
)

// EXAMPLE CODE: A main function.
func main() {

	fmt.Println("Starting the application...")
	exitReason := lifecycle.Run(startupMaxDur, shutdownMaxDur, startupFns, shutdownFns)

	if exitReason.StartupErr != nil {
		// A startup function returned an error therefore the application should not ever run.
		fmt.Printf("Failed to start: %+v\n", exitReason.StartupErr)
	}

	if exitReason.RuntimeErr != nil {
		// The application ran and encountered an error that was reported by lifecycle.RuntimeErr().
		fmt.Printf("Runtime error triggered shutdown: %+v\n", exitReason.RuntimeErr)
	} else if exitReason.OsSignal != nil {
		// The application ran and received an OS shutdown signal.
		fmt.Printf("OS signal triggered shutdown: %+v\n", exitReason.OsSignal)
	}

	if len(exitReason.ShutdownErrs) > 0 {
		// 1 or more shutdown functions errored.
		fmt.Printf("Shutdown process errors: %+v\n", exitReason.ShutdownErrs)
	}

	fmt.Println("Exiting the application.")
}

func startServer(ctx context.Context) error {
	// Start the webserver in a separate goroutine and start serving HTTP traffic.
	go func() {
		fmt.Printf("Server started on -> %s\n", server.Addr)

		if err := server.ListenAndServe(); !errors.Is(http.ErrServerClosed, err) {
			// If it's not ErrServerClosed, something bad happened.
			lifecycle.RuntimeErr(err)
		}
	}()

	return nil
}

func stopServer(ctx context.Context) error {
	// Shutdown the webserver to stop serving HTTP traffic.
	return server.Shutdown(ctx)
}
