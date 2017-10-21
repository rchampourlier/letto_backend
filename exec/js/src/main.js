// Main JS execution script
//
// This script will load the environment to be provided
// to the workflows (`data`, `secrets` and `context`)
// and execute the scripts contained in 
// `workflows/<group>`.

const fs = require('fs');


// Loading the environment
const data = require("./data");
const secrets = require("./secrets");
const context = require(process.argv[2]);
const group = context.Group;

// Directory containing the workflows to be executed
const rootDir = "./workflows/";
const dir = rootDir + group + "/";
console.log("[JS] Processing workflows for group `" + group + "`");

fs.readdir(dir, (err, files) => {
  if (err != null) {
    console.log("[JS] Failed to read the `" + dir + "` directory: %s", err);
    return
  }

  var filesCount = files.length;
  for (var i = 0; i < filesCount; i++) {
    // Load and execute each workflow
    var workflow = files[i];
    console.log("[JS] Loading `" + workflow + "`");
    var workflowFunc = require(dir + workflow);
    console.log("[JS] Running `" + workflow + "`");
    workflowFunc(data, secrets, context);
    console.log("[JS] Finished `" + workflow + "`");
  }
});
