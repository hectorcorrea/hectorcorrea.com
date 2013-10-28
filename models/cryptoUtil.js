var crypto = require('crypto');

var createHash = function(text, salt) {

  var sha = crypto.createHash('sha1');
  sha.update(text);
  if(salt) {
    sha.update(salt);
  }

  var digest =  sha.digest('hex');
  console.log(text);
  console.log(salt);
  console.log(digest);
  console.log('------');
  return digest;

};

module.exports = {
  createHash: createHash
}

