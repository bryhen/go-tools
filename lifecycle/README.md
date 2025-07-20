# Lifecycle

Lifecycle manages the sequential startup, monitoring, and shutdown of an application.

A webserver could configure Lifecycle to the following or any other variation of startup and shutdown.
1) Connect to a database.
2) Load database configurations for an async process that runs in the background and start it.
3) Start an HTTP server to begin serving live traffic.
4) Wait for a shutdown signal (ie ctrl+c).
5) Shutdown the server to stop serving live traffic.
6) Shutdown the async process, making sure that the last batch of stuff in the async process gets handled.

See /_examples for more.