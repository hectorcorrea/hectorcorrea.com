var db = require('./userDb');
var cryptoUtil = require('./cryptoUtil');


exports.initialize = function(data, cb) {

  var initData = {
    user: data.user, 
    password: cryptoUtil.createHash(data.password, data.salt) 
  };

  db.initialize(initData, function(err) {
    cb(err);
  });

};


exports.changePassword = function(data, cb) {

  var newData = {
    user: data.user, 
    oldPassword: cryptoUtil.createHash(data.oldPassword, data.salt), 
    newPassword: cryptoUtil.createHash(data.newPassword, data.salt) 
  };

  db.changePassword(newData, function(err) {
    cb(err);
  });
  
};


exports.login = function(data, cb) {

  var loginData = {
    user: data.user, 
    password: cryptoUtil.createHash(data.password, data.salt) 
  };

  db.login(loginData, function(err, authKey) {
    cb(err, authKey);
  });
  
};


exports.validateSession = function(data, cb) {

  var sessionData = {
    user: data.user,
    token: data.token
  };

  db.validateSession(sessionData, function(err) {
    cb(err);
  });

};


exports.killSession = function(data, cb) {

  var sessionData = {
    user: data.user,
    token: data.token
  };

  db.killSession(sessionData, function(err) {
    cb(err);
  });

};



