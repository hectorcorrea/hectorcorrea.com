// Manages connecting to MongoDB
// Typical usage is:
//  
//  var mongoConnect = require("./mongoConnect");
//  mongoConnect.setup("url to the DB", "name of collection");
//  mongoConnect.connect(function(err, db) {
//    do something with the db
//  })
//

var logger = require('log-hanging-fruit').defaultLogger;
var MongoClient = require('mongodb').MongoClient;
var dbUrl = null; 
var db = null;

var dbOptions = {
  db: {},
  server: {
    auto_reconnect: true,
    socketOptions: {keepAlive: 1}
  },
  replSet: {},
  mongos: {}
};


var _connectToDb = function(callback) {

  logger.debug("Connecting...");
  MongoClient.connect(dbUrl, dbOptions, function(err, dbConn) {
    
    if(err) {
    
      logger.error("Connect 1 of 2 failed: " + err);
      MongoClient.connect(dbUrl, dbOptions, function(err2, dbConn2) {
        
        if(err2) {
          logger.error("Connect 2 of 2 failed: " + err2);
          callback(err2);
          return;
        } 

        logger.debug("Connected! (in second attempt)");
        db = dbConn2;
        callback(null, db);

      });

      return;
    }
    
    logger.debug("Connected!");
    db = dbConn;
    callback(null, db);
  
  });

}



var _validateConnection = function(callback) {

  // Ping the server to make sure things are still OK.
  //
  // If the connection is dropped because it was idle
  // the ping will restore it and the client won't 
  // even notice we lost connectivity. 
  // Ideally, auto_reconnect = true in the dbOptions
  // should take care of this, but I've had mixed
  // results with that.
  //
  // This ping is wasteful when the connection is OK
  // but I am willing to take the hit in order to guarantee
  // that users don't notice any connect/disconnect 
  // issues.

  logger.debug("Already connected, about to ping...");
  db.admin().ping(function(err) {
    
    if (err) {
      logger.debug("...existing connection is broken.");
      db = null;
      _connectToDb(callback);
    }
    else {
      logger.debug("...existing connections is OK");
      callback(null, db);
    }

  }); 

}


var execute = function(callback) {

  var isAlreadyConnected = (db != null);
  if(isAlreadyConnected) {
    _validateConnection(callback);
    return;
  }

  _connectToDb(callback);
};


var setup = function(connString) {
  if(dbUrl == null) {
    dbUrl = connString;
  }
};


module.exports = {
  setup: setup,
  execute: execute
};

