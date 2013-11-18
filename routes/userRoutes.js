var logger = require('log-hanging-fruit').defaultLogger;
var model = require('../models/userModel');

var oneHour = 1000 * 60 * 60;
var oneMonth = oneHour * 24 * 30;


var initialize = function(req, res) {

  var error;
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
        res.status(500).send('Users database already initialized.');
      }
      else {
        logger.info('Initialized OK');
        res.status(200).send('OK');
      }
    });

  }
  else {

    logger.warn('user.initialize - cannot be executed. ' + 
      'Make sure environment variables BLOG_USER, BLOG_PASSWORD, ' + 
      'and BLOG_SALT are defined.');
    res.status(401).send('Not authorized to initialize');

  }

};


var changePassword = function(req, res) {

  if(!req.isAuth) {
    logger.warn('user.changePassword - user not authenticated');
    return res.status(401).send('Cannot change password');
  }

  var user = req.body.user;
  var oldPassword = req.body.oldPassword;
  var newPassword = req.body.newPassword;

  if(!user || !oldPassword || !newPassword) {
    logger.warn('user.changePassword - not all parameters were received');
    return res.status(401).send('Cannot change password');
  }

  logger.info('user.changePassword');
  var salt = process.env.BLOG_SALT || null;
  var data = {user: user, oldPassword: oldPassword, newPassword: newPassword, salt: salt};
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
    res.clearCookie('authToken');
    return res.status(401).send('Cannot login without a username');
  }

  if(!password) {
    logger.warn('user.login - no password received');
    res.clearCookie('authToken');
    return res.status(401).send('Cannot login without a password');
  }

  logger.info('user.login');
  var salt = process.env.BLOG_SALT || null;
  var data = {user: user, password: password, salt: salt};
  var m = model.user(req.app.settings.config.dbUrl);

  m.login(data, function(err, authToken) {
    if(err) {
      logger.error(err);
      res.clearCookie('authToken');
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


var logout = function(req, res, next) {

  var user = req.cookies.user;
  var token = req.cookies.authToken;
  if(!req.isAuth || !user || !token) {
    logger.warn('user.logout - user not logged in');
    req.isAuth = false;
    res.clearCookie('authToken');
    return res.status(401).send('Not logged in');
  }

  logger.info('Logging user out [' + user + ']');
  var data = {user: user, token: token};
  var m = model.user(req.app.settings.config.dbUrl);
  m.killSession(data, function(err) {

    req.isAuth = false;
    res.clearCookie('authToken');

    if(err) {
      logger.error(err);
      return res.status(401).send('Error logging out');
    }

    logger.info('Logged out OK');
    return res.status(200).send('OK');
  });

}


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
      res.clearCookie('authToken');
    }
    else {
      logger.info('Session is OK');
      req.isAuth = true;
    }

    next();
  });

};


module.exports = {
  initialize: initialize,
  changePassword: changePassword,
  login: login,
  logout: logout,
  validateSession: validateSession
}
