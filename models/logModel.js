var fs = require('fs');
var logger = require('log-hanging-fruit').defaultLogger;
var path = require('path');

var getLog = function(logPath, date, callback) {

  date = date.replace(/-/g,'_');
  var fileName = path.join(logPath, date + '.log');

  console.log('reading log:' + fileName);
  fs.readFile(fileName, function(err, text) {

    if(err) {
      callback('Error reading log file: ' + err);
      return;
    }

    text = text.toString().replace(/\r\n/g,'<br/>');
    callback(null, text);

  });

}

module.exports = {
  getLog: getLog
}