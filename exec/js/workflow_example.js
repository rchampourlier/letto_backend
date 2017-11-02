module.exports = (data, secrets, event, console) => {

  const ovh = require('ovh');
  const request = require('request');

  console.log('Hello world!');

  console.log('Accessible data:');
  console.log(data);

  console.log('Event\'s data:');
  console.log(event);
};
