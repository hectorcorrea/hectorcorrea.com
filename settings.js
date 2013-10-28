var fs = require('fs');

exports.loadSync = function(fileName) {

  if(!fs.existsSync(fileName)) {
    text = "{}";
  }
  else {
    text = fs.readFileSync(fileName, 'utf8');
  }

  return JSON.parse(text);
}
