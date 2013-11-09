var mongoConnect = require("./mongoConnect");


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

      // Update the password
      // TODO: We should blow away all sessions but the current one
      var id = items[0]._id;
      var newData = {$set: {password : data.password}};  // sessions : []
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
      collection.update({user:login.user}, {'$addToSet':{sessions:session}}, function(err, count) {

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


module.exports = {
  setup: setup,
  initialize: initialize,
  changePassword: changePassword,
  login: login,
  validateSession: validateSession
};

