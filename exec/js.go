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
func RunJS(dir string) error {
	cfg := config(dir)
	js.Run(cfg)
	return nil
}

func config(dir string) js.DockerConfig {
	var cfg = js.DockerConfig{
		Image:      "node:4",
		Command:    []string{"node", "./test.js"},
		Volumes:    map[string]struct{}{"/usr/src/app": {}},
		WorkingDir: "/usr/src/app",
		Binds:      []string{dir + ":/usr/src/app"},
	}
	return cfg
}
