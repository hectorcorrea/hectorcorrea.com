var mongoConnect = require("mongoconnect");


var setup = function(dbConnString) {
  mongoConnect.setup(dbConnString);
};


var initialize = function(data, callback) {

  mongoConnect.execute(function(err, db) {

    var collection = db.collection("users");
    collection.find({}).toArray(function(err, items){
      
      if(err) return callback(err);

      if(items.length != 0) {
        callback("Error: users collection already initialized");
        return;
      }

      collection.save(data, function(err, savedCount){
        callback(err);
      });

    });

  });

};


var changePassword = function(data, callback) {

  mongoConnect.execute(function(err, db) {

    var collection = db.collection("users");
    collection.find({user: data.user}).toArray(function(err, items){
      
      if(err) return callback(err);

      if(items.length == 0) {
        callback('Error: User [' + data.user + '] not found');
        return;
      }

      if(items.length > 1) {
        callback('Error: More users than expected found for [' + data.user + ']');
        return;
      }

      var dbPassword = items[0].password;
      if(dbPassword != data.oldPassword) {
        callback('Error: Incorrect old password received for [' + data.user + ']');
        return;
      }

      // Update the password
      // TODO: We should blow away all sessions but the current one
      var id = items[0]._id;
      var newData = {$set: {password : data.newPassword}};  // sessions : []
      collection.update({_id:id}, newData,  function(err){
        callback(err);
      });

    });

  });

};


var getNewSession = function(days) {
  // Calculate a date X days in the future
  // http://stackoverflow.com/a/11338904/446681
  var now = new Date();
  var oneDay = 1000 * 60 * 60 * 24;
  var totalMs = days * oneDay; 
  var expires = new Date(now.getTime() + totalMs);
    
  // A random number
  // http://stackoverflow.com/a/8084248/446681
  var token = Math.random().toString(36).substring(7);
  
  var session = {token: token, expires: expires}
  return session;
};


var login = function(login, callback) {

  mongoConnect.execute(function(err, db) {

    var collection = db.collection("users");
    collection.find({user: login.user}).toArray(function(err, items){
      
      if(err) return callback(err);

      if(items.length == 0) {
        callback('Error: User [' + login.user + '] not found');
        return;
      }

      if(items.length > 1) {
        callback('Error: More users than expected found for [' + login.user + ']');
        return;
      }

      if(login.password !== items[0].password) {
        callback('Error: Invalid password for [' + login.user + ']');
        return;
      }

      // Save the session
      var days = 30;
      var session = getNewSession(days);
      var query = {user: login.user};
      var action = {'$addToSet': {sessions:session}};
      collection.update(query, action, function(err, count) {

        if(err) {
          callback(err);
          return 
        }

        callback(null, session.token);

      });

    });

  });

};


var validateSession = function(session, callback) {
  var i, now;

  mongoConnect.execute(function(err, db) {

    var collection = db.collection("users");
    collection.find({user: session.user}).toArray(function(err, items){
      
      if(err) return callback(err);

      if(items.length == 0) {
        callback('Error: User [' + session.user + '] not found');
        return;
      }

      if(items.length > 1) {
        callback('Error: More users than expected found for [' + session.user + ']');
        return;
      }

      var sessions = items[0].sessions;
      if(sessions) {

        now = new Date();
        for(i = 0; i < sessions.length; i++) {
        
          if(sessions[i].token === session.token) {

            if(sessions[i].expires <= now) {
              return callback('Error: Session has expired for user [' + session.user + ']');
            }
            else {
              // Good session
              return callback(null);
            }

          }

        }

        return callback('Error: Session is not valid for user [' + session.user + ']');

      }

      callback('Error: No sessions were found for user [' + session.user + ']');

    });

  });

};


var killSession = function(session, callback) {
  var i, now;

  mongoConnect.execute(function(err, db) {

    var collection = db.collection("users");
    var query = {'user': session.user};
    var action = {'$pull': {'sessions': {'token': session.token}}}
    collection.update(query, action, function(err, count){
      
      if(err) return callback(err);

      if(count < 1 ) {
        callback('Error: Session was not found [' + session.user + '/' + session.token + ']');
        return;
      }

      if(count > 1 ) {
        callback('Error: More than one session was found [' + session.user + '/' + session.token + ']');
        return;
      }

      // just what we wanted.
      callback(null);

    });

  });

};

module.exports = {
  setup: setup,
  initialize: initialize,
  changePassword: changePassword,
  login: login,
  validateSession: validateSession,
  killSession: killSession
};

