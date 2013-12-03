var crypto = require('crypto');

exports.createHash = function(text, salt) {

  var sha = crypto.createHash('sha1');
  sha.update(text);
  if(salt) {
    sha.update(salt);
  }

  var digest =  sha.digest('hex');
  return digest;

};


