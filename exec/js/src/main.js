// Main JS execution script
//
// This script will load the environment to be provided
// to the workflows (`data`, `secrets` and `context`)
// and execute the scripts contained in 
// `workflows/<group>`.

const fs = require('fs');
const { Console } = require('console');

const timestamp = () => {
  return (new Date()).toISOString();
};

// Overriding `console` so that logs from workflows
// are prefixed with the workflow's name.
const workflowConsole = (group, workflow) => {
  c = new Console(process.stdout, process.stderr);
  rewriteLogFunc = (func) => {
    return (...args) => {
      args[0] = timestamp() + " [Exec/JS/" + group + "/" + workflow + "] " + args[0];
      func(...args);
    }
  }
  c.log = rewriteLogFunc(c.log);
  c.info = c.log;
  c.error = rewriteLogFunc(c.error);
  c.trace = rewriteLogFunc(c.trace);
  c.warn = rewriteLogFunc(c.warn);
  return c;
};

// Logger for main.js
mainLog = (message) => { console.log(timestamp() + " [Exec/JS] " + message); };

// Loading the environment
const data = require("./data");
const secrets = require("./secrets");
const context = require(process.argv[2]);
const group = context.Group;

// Directory containing the workflows to be executed
const rootDir = "./workflows/";
const dir = rootDir + group + "/";
mainLog("Starting workflows in group `" + group + "`");

fs.readdir(dir, (err, files) => {
  if (err != null) {
    mainLog("Failed to read the `" + dir + "` directory: %s", err);
    return
  }

  var filesCount = files.length;
  for (var i = 0; i < filesCount; i++) {
    var workflow = files[i];

    // Load and execute each workflow
    mainLog("Loading workflow `" + workflow + "`");
    var workflowFunc = require(dir + workflow);
    mainLog("Running workflow `" + workflow + "`");
    workflowFunc(data, secrets, context, workflowConsole(group, workflow));
    mainLog("Completed workflow `" + workflow + "`");
  }

});
