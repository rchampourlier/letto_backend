// Workflow
// 
// ### Trigger:
//
// When a card is:
//   - commented,
//   - a label is added or removed.
//
// We ignore when:
//   - the card is created, because the label is not set yet,
//   - the card is updated, otherwise we would trigger again when the due 
//     date is updated and go into a loop.
//
// ### Workflow:
//
// Set the _due date_ to a new value. The new value is determined by adding 
// a number of days corresponding to the label to the current date (today).
//
// Mapping:
//   - "action to take": 7 days
//   - "active contact": 30 days
//   - "passive contact": 90 days

module.exports = (data, secrets, context, console) => {

  trelloAPICall = (key, token, verb, path, queryString, onSuccess, onError) => {
    path += '?key=' + key;
    path += '&token=' + token;
    if (queryString.length > 0) {
      path += '&' + queryString;
    }

    const options = {
      hostname: 'api.trello.com',
      port: 443,
      path: path,
      method: verb
    };

    const req = https.request(options, (res) => {
      var dataStr = '';
      res.on('data', (chunk) => {
        dataStr += chunk;
      });

      res.on('end', function () {
        onSuccess(dataStr);
      });
    });

    req.on('error', (e) => {
      onError(e);
    });

    req.end();
  };

  // Main workflow code
  const https = require('https');
  const body = context.Request.Body;
  var actionType = body.action.type;

  // Trigger:
  //   - commented: "commentCard"
  //   - added label: "addLabelToCard"
  //   - removed label: "removeLabelFromCard"
  //
  // Action types are documented here:
  //   [Trello reference](https://developers.trello.com/v1.0/reference#action-types)
  //
  if (actionType == 'commentCard' ||
      actionType == 'addLabelToCard' ||
      actionType == 'removeLabelFromCard') {

    var cardID = body.action.data.card.id;

    // We need to fetch the card to get its labels.
    // GET /cards/<id>
    trelloAPICall(secrets.trello.api_key, secrets.trello.api_token, "GET", "/1/cards/" + cardID, '',
      
    (data) => { // success callback
      var cardData = JSON.parse(data);
      var labels = cardData.labels;
      var labelNames = labels.map( (item) => { return item.name; } );

      var nDays;
      if (labelNames.includes("action to take")) {
        nDays = 7;
      }
      else if (labelNames.includes("active contact")) {
        nDays = 30;
      }
      else if (labelNames.includes("passive contact")) {
        nDays = 90;
      }

      // TODO: use the Trello event's time instead of now
      var date = new Date();
      var newDate = new Date(date.setTime(date.getTime() + nDays * 86400000));
      // LIMIT: arbitrary adding time, without taking calendar rules into account
      //   (e.g. timezones). Should be ok though in the context of this workflow,
      //   since we only want to set a due date, with no time.

      // Update the card with the new due date
      trelloAPICall(secrets.trello.api_key, secrets.trello.api_token, "PUT", "/1/cards/" + cardID, 'due=' + newDate.toISOString(), () => {}, (error) => {
        console.log("Could not update the card (" + error + ")");
      });
    }, 
      
    (error) => { // error callback
      console.error("Could not fetch the card (" + error + ")");
    });
  };
};
