package exadv

/*
	This code can be copied and ran. It will start an HTTP server at localhost:8080.
	Requests at any path except /error will return 200, echo that path, and mock batch processing with duration mockDbBatchSize * durPerReq.
	If the path is /error the server will report an error and the application will shut down.
	When ending the application with ctrl+c, any remaining batch items will be processed and logged as normal (unless timed out at shutdownMaxDur), despite the SIGINT command produced by ctrl+c.
*/

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bryhen/go-tools/lifecycle"
)

const (
	// How long to wait before startup should fail since it's taking too long.
	// In real code, make sure this variable provides more than enough time for a slow, sequential startup.
	// Adjust this variable to change how long the maximum duration for startup should take. In this example, all startup functions return essentially immediately.
	startupMaxDur = time.Second * 10

	// How long to wait before shutdown should fail since it's taking too long.
	// In real code, make sure this variable is longer than the maximum amount of time all processes should take, sequentially, to exit.
	// Adjust this variable to change how long the maximum duration for shutdown should take.
	shutdownMaxDur = time.Second * 300
)

var (
	// The functions to run during startup, in order.
	// In this example, if startDB came after startReqBatchProcDone, startup would error and the appication would immediately exit.
	startupFns = []func(context.Context) error{startDB, startReqBatchProcDone, startServer}

	// The functions to run during shutdown, in order.
	// In this example, stopServer must come before stopReqBatchProcDone or else requests would hang then error and logs would be missed.
	shutdownFns = []func(context.Context) error{stopServer, stopReqBatchProcDone}

	// EXAMPLE CODE: Send that an http request has been received.
	reqCh = make(chan string)

	// EXAMPLE CODE: Closed when the reqBatchProcDone process exits.
	reqBatchProcDone = make(chan struct{})

	// EXAMPLE CODE: The number of http requests received.
	reqCount = 0

	// EXAMPLE CODE: An http server that communicates with reqBatchProc, which is running in a separate goroutine, whenever it receives a request.
	// You can make requests to see this in action at localhost:8080
	server = &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCh <- r.URL.Path
			w.Write([]byte(r.URL.Path))
		}),
	}

	// EXAMPLE CODE: Mocks the database call. If not initialized, anything that calls this will error.
	// startDB sets this function's value so that startReqBatchProcDone does not error in this example.
	mockDbQuery = func() (int, error) {
		return 0, fmt.Errorf("database query failed because startDB has not been ran")
	}

	// EXAMPLE CODE:
	// Adjust this variable to change how many items each batch processes. It's a mock of a Database queried value for batch size.
	// In real code, this would probably be a env, config, or build file variable, but it serves this example.
	// It represents anything that depends on anything else to have started first.
	mockDbBatchSize = 5

	// EXAMPLE CODE:
	// Adjust this variable to change how long to wait for a batch to complete based on how many items it contains.
	// In real code, this would be the processing time of whatever process is actually being ran.
	// It represents an asnyc process taking time and it can be helpful to see latency in action.
	mockDurPerReq = time.Second * 2
)

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

func startDB(ctx context.Context) error {
	// Mock connect to the db.
	fmt.Println("Connected to DB.")

	// Mock the db query so that this function returns a constant integer and no error.
	mockDbQuery = func() (int, error) {
		return mockDbBatchSize, nil
	}

	return nil
}

func startReqBatchProcDone(ctx context.Context) error {
	batchSize, err := mockDbQuery() // Returns an error if startDB() hasn't been ran.
	if err != nil {
		close(reqBatchProcDone) // We're not starting ReqBatchProc because startup failed so it's already shut down.
		return err
	}

	go func() {
		fmt.Println("ReqBatchProc has started running in the background...")
		batch := make([]string, 0, batchSize)

		// Receive from the channel that a request was handled.
		for path := range reqCh {
			fmt.Printf("\tReqBatchProc received a request that will be handled in a batch later at path: %s\n", path)
			batch = append(batch, path)
			if len(batch) == batchSize {
				handleBatch(batch)
				batch = batch[:0]
			}
		}

		// The for loop above exited meaning that reqCh has been closed. Handle the last batch then signal that this process is shutting down.
		fmt.Println("ReqBatchProc has exited its for loop and is handling the last batch.")
		handleBatch(batch)
		close(reqBatchProcDone)

		fmt.Printf("Total requests received is: %d\n", reqCount)
	}()

	return nil
}

func handleBatch(batch []string) {
	dur := mockDurPerReq * time.Duration(len(batch))
	fmt.Printf("Working on a batch that will take %.1f seconds...\n", dur.Seconds())

	// Simulate some work for handling the batch.
	time.Sleep(dur)

	reqCount += len(batch)
	fmt.Printf("Batch Size: %d. Req count at: %d. Paths this batch: %s", len(batch), reqCount, "\n\t"+strings.Join(batch, "\n\t")+"\n\t\n")
}

func startServer(ctx context.Context) error {
	// Start the webserver in a separate goroutine...
	go func() {
		fmt.Printf("Server started on %s\n\n", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(http.ErrServerClosed, err) {
			// If the server isn't closed, something bad happened that we want to know about.
			lifecycle.RuntimeErr(err)
		}
	}()

	return nil
}

func stopServer(ctx context.Context) error {
	// Shutdown the webserver...
	fmt.Println("Shutting down server...")
	return server.Shutdown(ctx)
}

func stopReqBatchProcDone(ctx context.Context) error {
	fmt.Println("Shutting down ReqBatchProc...")

	// Signals the ReqBatchProc to exit its for loop.
	close(reqCh)

	// Signals that the ReqBatchProc has exited.
	<-reqBatchProcDone

	return nil
}
