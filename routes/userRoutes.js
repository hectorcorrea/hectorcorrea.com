var logger = require('log-hanging-fruit').defaultLogger;
var model = require('../models/userModel');
var oneHour = 1000 * 60 * 60;
var oneMonth = oneHour * 24 * 30;

var error = function(req, res, title, err) {
  logger.error(title + ' ' + err);
  res.status(500).render('error', {title: title, error: err});
};


exports.changePassword = function(req, res) {
  if(!req.isAuth) {
    logger.warn('user.changePassword - user not authenticated');
    return error(req, res, 'Cannot change password');
  }
  res.render('changePassword', {user: req.cookies.user})
}


exports.changePasswordPost = function(req, res) {

  if(!req.isAuth) {
    logger.warn('user.changePassword - user not authenticated');
    return error(req, res, 'Cannot change password');
  }

  var user = req.body.user;
  var oldPassword = req.body.oldPassword;
  var newPassword = req.body.newPassword;
  var repeatPassword = req.body.repeatPassword;

  if(!user || !oldPassword || !newPassword || !repeatPassword) {
    logger.warn('user.changePassword - not all parameters were received');
    return error(req, res, 'Cannot change password');
  }

  if(newPassword !== repeatPassword) {
    logger.warn("user.changePassword - newPassword and repeatPassword don't match");
    return error(req, res, 'Cannot change password');
  }

  logger.info('user.changePassword');
  var salt = process.env.BLOG_SALT || null;
  var data = {user: user, oldPassword: oldPassword, newPassword: newPassword, salt: salt};

  model.changePassword(data, function(err) {
    if(err) {
      logger.error(err);
      return error(req, res, 'Error changing password');
    }
    else {
      logger.info('user.changePassword - Password changed OK');
      res.redirect(301, "/")
    }
  });

};


exports.login = function(req, res) {
  res.render('login')
}


exports.loginPost = function(req, res) {

  var user = req.body.user;
  var password = req.body.password;

  if(!user) {
    logger.warn('user.login - no user received');
    res.clearCookie('authToken');
    error(req, res, 'Cannot login without a username');
  }

  if(!password) {
    logger.warn('user.login - no password received');
    res.clearCookie('authToken');
    error(req, res, 'Cannot login without a password');
  }

  logger.info('user.login');
  var salt = process.env.BLOG_SALT || null;
  var data = {user: user, password: password, salt: salt};

  model.login(data, function(err, authToken) {
    if(err) {
      logger.error(err);
      res.clearCookie('authToken');
      error(req, res, "Cannot login");
    }
    else {
      logger.info('Logged in OK');
      res.cookie('authToken', authToken, {maxAge: oneHour});
      res.cookie('user', user, {maxAge: oneMonth});
      res.redirect(301, "/?action=login")
    }
  });

};


exports.logout = function(req, res, next) {

  var user = req.cookies.user;
  var token = req.cookies.authToken;
  if(!req.isAuth || !user || !token) {
    logger.warn('user.logout - user not logged in');
    req.isAuth = false;
    res.clearCookie('authToken');
    error(req, res, 'Not logged in');
  }

  logger.info('Logging user out [' + user + ']');
  var data = {user: user, token: token};

  model.killSession(data, function(err) {

    req.isAuth = false;
    res.clearCookie('authToken');

    if(err) {
      logger.error(err);
      error(req, res, 'Error logging out');
    }

    logger.info('Logged out OK');
    res.redirect(301, "/?action=logout")
  });

}


exports.validateSession = function(req, res, next) {

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

  model.validateSession(data, function(err) {

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

