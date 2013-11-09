var crypto = require('crypto');

var createHash = function(text, salt) {

  var sha = crypto.createHash('sha1');
  sha.update(text);
  if(salt) {
    sha.update(salt);
  }

  var digest =  sha.digest('hex');
  return digest;

};

module.exports = {
  createHash: createHash
}

