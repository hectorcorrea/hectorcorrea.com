var dbCollection = "users";
var mongoConnect = require("./mongoConnect");


var setup = function(dbConnString) {
  mongoConnect.setup(dbConnString);
};


var initialize = function(data, callback) {

  mongoConnect.execute(function(err, db) {

    var collection = db.collection(dbCollection);
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

    var collection = db.collection(dbCollection);
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

      // Force an update
      data._id = items[0]._id;

      collection.save(data, function(err){
        callback(err);
      });

    });

  });

};


var login = function(login, callback) {

  mongoConnect.execute(function(err, db) {

    var collection = db.collection(dbCollection);
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

      // TODO: Generate an authkey and save it to the DB
      callback(null, 'ok');

    });

  });

};


module.exports = {
  setup: setup,
  initialize: initialize,
  changePassword: changePassword,
  login: login
};

