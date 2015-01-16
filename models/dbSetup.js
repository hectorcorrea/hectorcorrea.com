var logger = require('log-hanging-fruit').defaultLogger;
var mongoConnect = require('mongoconnect');

var addAdminUser = function(db, user) {

  logger.info('Validating users...');
  db.collection('users').count(function(err,count) {

    if(err) {
      logger.error('Error validating users: ' + err);
      return;
    }

    if(count > 0) {
      logger.info('...at least one user was found.');
      return;
    }

    if(user.user == null || user.password == null) {
      var errorMsg = "...no default BLOG_USER and BLOG_PWD found on the environment.\n" +
      "No default user will be added, run \n\n" +
      "BLOG_USER=xx BLOG_PWD=yy node server.js\n\n" +
      "to create a default user.";
      logger.error(errorMsg);
      return;
    }

    logger.info('...no users were found, adding initial user.');

    var model = require('../models/userModel');
    model.initialize(user, function(err) {
      if(err) {
        logger.error(err);
      }
      else {
        logger.info('...default user [' + user.user + '] added');
      }
    });

  });

}


exports.init = function(dbConnString, defaultUser) {

  mongoConnect.setup(dbConnString, null, true);
  
  mongoConnect.execute(function(err, db) {

    if(err) {
      logger.error('Could not perform DB initialization tasks. Error: ' + err);
      return;
    }

    addAdminUser(db, defaultUser);

    logger.info('Validating index by key...');
    db.collection('blog').ensureIndex({key:1}, function(err,ix) {

      if(err) {
        logger.error('Error creating index by key: ' + err);
        return;
      }

      logger.info('...done validating index by key.');
    });

    logger.info('Validating index by postedOn...');
    db.collection('blog').ensureIndex({postedOn:-1}, function(err,ix) {

      if(err) {
        logger.error('Error creating index by postedOn: ' + err);
        return;
      }

      logger.info('...done validating index by postedOn.');
    });

  }); 

};

