package exec

import (
	"os"

	"github.com/rchampourlier/letto_go/exec/js"
)

func config() js.DockerConfig {
	dir, err := os.Getwd()
	if err != nil {
		panic("Could not get current working dir. Abandoning.")
	}

	var cfg = js.DockerConfig{
		Image:      "node:4",
		Command:    []string{"node", "./exec/js/test.js"},
		Volumes:    map[string]struct{}{"/usr/src/app": {}},
		WorkingDir: "/usr/src/app",
		Binds:      []string{dir + ":/usr/src/app"},
	}
	return cfg
}

// RunJS execute the specified script in a JS execution
// environment. The script is passed the `body` string
// (made available j
func RunJS(group string, dir string) error {
	cfg := config()
	js.Run(cfg)
	return nil
}
