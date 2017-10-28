module.exports = (data, secrets, context, console) => {
  var ovh = require('ovh')(secrets.ovh);
  ovh.request('GET', '/sms', function (err, serviceNames) {
    if (err) {
      console.error(err, serviceName);
    }
    else {
      if (serviceNames.constructor === Array) {
        serviceName = serviceNames[0];
      }
      else {
        serviceName = serviceNames;
      }
      console.log("Retrieved /sms service name: `" + serviceName + "`");
      console.log("Posting new SMS");
      ovh.request('POST', '/sms/' + serviceName + '/jobs', {
        message: 'Hello World!',
        senderForResponse: true,
        receivers: [secrets.myPhoneNumber]
      }, function (errsend, result) {
        console.error(errsend, result);
      });
    }
  });
};
