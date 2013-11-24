var logger = require('log-hanging-fruit').defaultLogger;
var mongoConnect = require('mongoconnect');

var init = function(dbConnString) {

  mongoConnect.setup(dbConnString, null, true);
  mongoConnect.execute(function(err, db) {

    if(err) {
      logger.error('Could not perform DB initialization tasks. Error: ' + err);
      return;
    }

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


module.exports = {
  init: init
}