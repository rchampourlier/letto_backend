package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/events"
	"github.com/rchampourlier/letto_go/util"
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
	timestamp := util.Timestamp(time.Now())
	fileName := fmt.Sprintf("%s.json", timestamp)
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

func logTraceError(err error) error {
	fmt.Printf("Failed to write event file (%s)\n", err)
	return err
}
