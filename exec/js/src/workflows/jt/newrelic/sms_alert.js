module.exports = (data, secrets, context, console) => {
  var ovh = require('ovh')(secrets.ovh);
  ovh.request('GET', '/sms', function (err, serviceName) {
    if (err) {
      console.log(err, serviceName);
    }
    else {
      ovh.request('POST', '/sms/' + serviceName + '/jobs', {
        message: 'Hello World!',
        senderForResponse: true,
        receivers: [secrets.myPhoneNumber]
      }, function (errsend, result) {
        console.log(errsend, result);
      });
    }
  });
};
