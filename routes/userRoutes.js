var logger = require('log-hanging-fruit').defaultLogger;
var model = require('../models/userModel');

var oneHour = 1000 * 60 * 60;
var oneMonth = oneHour * 24 * 30;

var initialize = function(req, res) {

  var user = process.env.BLOG_USER;
  var password = process.env.BLOG_PASSWORD;
  if(user && password) {

    logger.info('user.initialize');

    var salt = process.env.BLOG_SALT || null;
    var data = {user: user, password: password, salt: salt};
    var m = model.user(req.app.settings.config.dbUrl);

    m.initialize(data, function(err) {
      if(err) {
        logger.error(err);
        res.status(500).send('Error initializing users');
      }
      else {
        logger.info('Initialized OK');
        res.status(200).send('OK');
      }
    });

  }
  else {

    logger.warn('user.initialize - cannot be executed');
    res.status(401).send('Not authorized to initialize');

  }

};


var changePassword = function(req, res) {

  var user = req.body.user;
  var password = req.body.password;

  if(!user) {
    logger.warn('user.changePassword - no user received');
    return res.status(401).send('Cannot change password');
  }

  if(!password) {
    logger.warn('user.changePassword - no password received');
    return res.status(401).send('Cannot change password');
  }

  logger.info('user.changePassword');
  var salt = process.env.BLOG_SALT || null;
  var data = {user: user, password: password, salt: salt};
  var m = model.user(req.app.settings.config.dbUrl);

  m.changePassword(data, function(err) {
    if(err) {
      logger.error(err);
      res.status(500).send('Error changing password');
    }
    else {
      logger.info('user.changePassword - Password changed OK');
      res.status(200).send('OK');
    }
  });

};


var login = function(req, res) {

  var user = req.body.user;
  var password = req.body.password;

  if(!user) {
    logger.warn('user.login - no user received');
    res.cookie('authToken', null, {maxAge: oneHour});
    return res.status(401).send('Cannot login without a username');
  }

  if(!password) {
    logger.warn('user.login - no password received');
    res.cookie('authToken', null, {maxAge: oneHour});
    return res.status(401).send('Cannot login without a password');
  }

  logger.info('user.login');
  var salt = process.env.BLOG_SALT || null;
  var data = {user: user, password: password, salt: salt};
  var m = model.user(req.app.settings.config.dbUrl);

  m.login(data, function(err, authToken) {
    if(err) {
      logger.error(err);
      res.cookie('authToken', null, {maxAge: oneHour});
      res.status(500).send('Cannot login');
    }
    else {
      logger.info('Logged in OK');
      res.cookie('authToken', authToken, {maxAge: oneHour});
      res.cookie('user', user, {maxAge: oneMonth});
      res.status(200).send(authToken);
    }
  });

};


var validateSession = function(req, res, next) {

  var user = req.cookies.user;
  if(!user) {
    logger.warn('No user to validate session');
    return next();
  }

  var token = req.cookies.authToken;
  if(!token) {
    logger.warn('No authToken to validate session');
    return next();
  }

  logger.info('Validating session for [' + user + ']');
  var data = {user: user, token: token};
  var m = model.user(req.app.settings.config.dbUrl);
  m.validateSession(data, function(err) {

    if(err) {
      logger.error(err);
      res.cookie('authToken', null, {maxAge: oneHour});
    }
    else {
      logger.info('Session is OK');
    }

    next();
  });

};


module.exports = {
  initialize: initialize,
  changePassword: changePassword,
  login: login,
  validateSession: validateSession
}
