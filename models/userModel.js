var db = require('./userDb');
var cryptoUtil = require('./cryptoUtil');
var dbUrl = null;

var initialize = function(data, cb) {

  var initData = {
    user: data.user, 
    password: cryptoUtil.createHash(data.password, data.salt) 
  };

  db.setup(dbUrl);
  db.initialize(initData, function(err) {
    cb(err);
  });

};


var changePassword = function(data, cb) {

  var newData = {
    user: data.user, 
    password: cryptoUtil.createHash(data.password, data.salt) 
  };

  db.setup(dbUrl);
  db.changePassword(newData, function(err) {
    cb(err);
  });
  
};


var login = function(data, cb) {

  var loginData = {
    user: data.user, 
    password: cryptoUtil.createHash(data.password, data.salt) 
  };

  db.setup(dbUrl);
  db.login(loginData, function(err, authKey) {
    cb(err, authKey);
  });
  
};

var publicApi = {
  initialize: initialize,
  changePassword: changePassword,
  login: login
};


module.exports.user = function(dbConnString) {
  dbUrl = dbConnString;
  return publicApi;
};

