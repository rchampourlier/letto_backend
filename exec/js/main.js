const data = require("./data");
const credentials = require("./credentials");
const request = require("./request");
const body = request.body;
const headers = request.headers;
const https = require('https');

console.log("Trello workflow: START");
console.log("");
//console.log("** BODY **");
//console.log(body);

var actionType = body.action.type;
if (actionType == 'createCard') {
  var cardID = body.action.data.card.id;
  var cardName = body.action.data.card.name;

  var newName = cardName + ' (edited by Letto)';
  newName = encodeURIComponent(newName);

  var path = '/1/cards/' + cardID;
  path = path + '?name=' + newName;
  path = path + '&key=' + credentials.trello.api_key;
  path = path + '&token=' + credentials.trello.api_token;

  const options = {
    hostname: 'api.trello.com',
    port: 443,
    path: path,
    method: 'PUT'
  };

  const req = https.request(options, (res) => {
    console.log('statusCode:', res.statusCode);
    console.log('headers:', res.headers);

    res.on('data', (d) => {
      process.stdout.write(d);
    });
  });

  req.on('error', (e) => {
    console.error(e);
  });

  req.end();
}


console.log("Trello workflow: END");
