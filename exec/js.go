package exec

import (
	"github.com/rchampourlier/letto_go/exec/js"
)

// RunJS execute the specified script in a JS execution
// environment. The script is passed the `body` string.
//
// `dir` is expected to be the path to a local directory
// which contains the JS code to be executed. The container
// will run the `main.js` file present in this directory.
//
// `main` is the path, local to `dir`, of the main JS file
// to be run by the container.
func RunJS(dir string, main string) error {
	cfg := config(dir, main)
	js.Run(cfg)
	return nil
}

func config(dir string, main string) js.DockerConfig {
	var cfg = js.DockerConfig{
		Image:      "node:4",
		Command:    []string{"node", main},
		Volumes:    map[string]struct{}{"/usr/src/app": {}},
		WorkingDir: "/usr/src/app",
		Binds:      []string{dir + ":/usr/src/app"},
	}
	return cfg
}
