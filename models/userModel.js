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


var validateSession = function(data, cb) {

  var sessionData = {
    user: data.user,
    token: data.token
  };

  db.setup(dbUrl);
  db.validateSession(sessionData, function(err) {
    cb(err);
  });

};


var publicApi = {
  initialize: initialize,
  changePassword: changePassword,
  login: login,
  validateSession: validateSession
};


module.exports.user = function(dbConnString) {
  dbUrl = dbConnString;
  return publicApi;
};

