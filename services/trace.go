package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/events"
)

const tracesDirPath = "./traces"

// Trace stores the service's config.
type Trace struct {
	Fs afero.Fs
}

// NewTrace creates a new TraceService with
// the specified `afero.Fs` filesystem abstraction.
func NewTrace(fs afero.Fs) *Trace {
	return &Trace{Fs: fs}
}

// OnReceivedWebhook write a trace of the specified trigger, allowing
// for future reuse (e.g. to try a new version of a workflow or
// inspect the payload).
func (s *Trace) OnReceivedWebhook(event events.ReceivedWebhook) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return logTraceError(err)
	}
	dirPath := path.Join(rootDir, tracesDirPath, event.Group)
	err = s.Fs.MkdirAll(dirPath, 0777)
	if err != nil {
		return logTraceError(err)
	}
	fileName := fmt.Sprintf("%s-trigger.json", event.UniqueID)
	filePath := path.Join(dirPath, fileName)

	// Write the content of the event to a file
	eventAsJSON, err := json.Marshal(event)
	if err != nil {
		return logTraceError(err)
	}
	err = afero.WriteFile(s.Fs, filePath, eventAsJSON, 0777)
	if err != nil {
		return logTraceError(err)
	}

	return nil
}

// OnCompletedWorkflows writes a trace of the workflows execution, using
// provided `events.CompletedWorkflows` data.
func (s *Trace) OnCompletedWorkflows(event events.CompletedWorkflows) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return logTraceError(err)
	}
	dirPath := path.Join(rootDir, tracesDirPath, event.Group)
	err = s.Fs.MkdirAll(dirPath, 0777)
	if err != nil {
		return logTraceError(err)
	}
	fileName := fmt.Sprintf("%s-logs.json", event.TriggerUniqueID)
	filePath := path.Join(dirPath, fileName)

	// Write the logs to a file
	logs := logsTrace{
		Stdout: event.Stdout,
		Stderr: event.Stderr,
		error:  event.Error,
	}
	logsAsJSON, err := json.Marshal(logs)
	if err != nil {
		return logTraceError(err)
	}
	err = afero.WriteFile(s.Fs, filePath, logsAsJSON, 0777)
	if err != nil {
		return logTraceError(err)
	}

	return nil
}

func logTraceError(err error) error {
	fmt.Printf("Failed to write file (%s)\n", err)
	return err
}

type logsTrace struct {
	Stdout string
	Stderr string
	error  string
}
