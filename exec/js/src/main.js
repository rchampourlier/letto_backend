// Main JS execution script
//
// This script will load the environment to be provided
// to the workflows (`data`, `secrets` and `context`)
// and execute the scripts contained in `workflows`.

const fs = require('fs');

const workflowsDir = "./workflows/";

// Loading the environment
const data = require("./data");
const secrets = require("./secrets");
const context = require(process.argv[2]);

fs.readdir(workflowsDir, (err, files) => {
  if (err != null) {
    console.log("[JS] Failed to read the `workflows` directory: %s", err);
    return
  }

  var filesCount = files.length;
  for (var i = 0; i < filesCount; i++) {
    // Load and execute each workflow
    var workflow = files[i];
    console.log("[JS] Loading `" + workflow + "`");
    var workflowFunc = require(workflowsDir + workflow);
    console.log("[JS] Running `" + workflow + "`");
    workflowFunc(data, secrets, context);
    console.log("[JS] Finished `" + workflow + "`");
  }
});
