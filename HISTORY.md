# History

## 2017-11-05

- `hostDataDir` is now always passed as the 1st argument of the application.

## 2017-11-03

- Added `SequenceID` to event to trace a sequence of events initiated by
  the same event.
- Fixed status code of `triggers/webhook` endpoint response.
- Updated Goa and regenerated code.

## 2017-11-02

- Refactoring events to have a single point of logging, the event bus.
  This will also make it easier to replay past triggers by simply
  firing the event again.

## 2017-10-31

- Fixed implementation to enable building and running the app inside
  the `docker` container.
- The current implementation requires setting the correct `...DataDir`
  paths in `main.go`. This needs to be addressed.

## 2017-10-29

- Making changes to enable it running inside a Docker container and spawning
  sibling containers.
- For now only works when run with `go run main.go`. When run using
  `docker-compose`, the `execjs` container is not executed.

## 2017-10-22

- Tracing standard and error logs from JS execution, writing them to `traces/...`
  files. Injecting customized `Console` in JS scripts to prefix the logs with
  a timestamp and the context, enabling differentiating logs from different
  workflows.
- Added support custom container images. The image is built when starting the
  main program, using the `Dockerfile` in the `exec/js/docker` directory.

## 2017-10-21

- Completed refactoring of JS execution.
- Implemented support for multiple workflows.

## 2017-10-15

- Wrote a more complex test workflow with Trello.
- Finished implementing tracing of received webhooks.
  - Enriched received webhooks' trace (added Method, URL, Host).
  - TriggersController is now injected with a `afero.Fs` to enable
    it injecting it in the services (`trace` and `run_workflows`).
  - Replaced `WorkflowData` by the event structure `events.ReceivedWebhook`.
- Started refactoring JS execution

## 2017-10-09

- Started implementing tracing of received webhooks.

## 2017-10-08

- Injecting credentials using non-committed `credentials.js`
  file.
- Finalized implementation of end-to-end workflow: Trello
  webhook triggering JS workflow execution.

## 2017-10-01

- Connect API's webhook endpoint with JS execution.
- Injecting JS data and scripts into the container using
  a tmp dir.

## 2017-09-31

- Prepared for JS execution within a Docker container.
- Executing Docker container from Go code using Docker
  Go SDK.

## 2017-09-26

- Implemented one of the _workflow_ endpoints (_create_).
- Started implementation of the _webhook_ trigger endpoint.
- Started implementation of JS execution:
  - Tested [otto](https://github.com/robertkrimen/otto). Does not fit, because it's limited to JS execution and does not include some libraries provided by browsers or NodeJS (such as `XMLHttpRequest`).
  - Will have to test an approach using Docker.

## 2017-09-16

- Started implementation by tweaking the generated scaffold.

## 2017-08-06

- Initialized the project and started with [Goa's Getting Started](https://goa.design/learn/guide/)).
