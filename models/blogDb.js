var logger = require('log-hanging-fruit').defaultLogger;
var dbCollection = "blog";
var mongoConnect = require("mongoconnect");


exports.getNewId = function(callback) {

  mongoConnect.execute(function(err, db) {

    if(err) return callback(err);

    var counters = db.collection('counters');
    var query = {'name': 'blogId'};
    var order = [['_id','asc']];
    var inc = {$inc:{'next':1}};
    var options = {new: true, upsert: true};
    counters.findAndModify(query, order, inc, options, function(err, doc) {

      if(err) {
        callback(err);
        return;
      }      

      var id = doc.next;
      callback(null, id);
    });

  });

};


exports.fetchAll = function(includeDrafts, callback) {
  var query = {};
  if(!includeDrafts) {
    query = {postedOn:{$ne:null}}
  }
  fetchList(query, callback);
};


var fetchList = function(query, callback) {

  mongoConnect.execute(function(err, db) {

    if(err) {
      logger.error("fetchList - connect error");
      db = null;
      return callback(err);
    }

    logger.debug("fetchList - connected ok");
    var collection = db.collection(dbCollection);
    var fields = {key: 1, title: 1, url: 1, summary: 1, postedOn: 1};
    var cursor = collection.find(query, fields).sort({postedOn:-1});
    cursor.toArray(function(err, items){
      if(err) {
        logger.error("fetchList - error reading");
        db = null;
        return callback(err);
      }
      logger.debug("fetchList - everything is OK");
      callback(null, items);
    });

  });
  
};


exports.fetchOne = fetchOne = function(key, callback) {

  mongoConnect.execute(function(err, db) {

    if(err) return callback(err);

    var collection = db.collection(dbCollection);
    var query = {key: key};
    collection.find(query).toArray(function(err, items){
      
      if(err) return callback(err);

      if(items.length === 1) {
        // just what we want
        callback(null, items[0]);
        return;
      }

      if(items.length > 1) {
        // oops! how come we got more than one?
        callback("Error: more than one record found for key [" + key + "]");
        return;
      }

      // no record found
      callback(null, null);

    });

  });

};


exports.fetchOneByUrl = function(url, callback) {

  mongoConnect.execute(function(err, db) {

    if(err) return callback(err);

    var collection = db.collection(dbCollection);
    var query = {url: url};
    collection.find(query).toArray(function(err, items){
      
      if(err) return callback(err);

      if(items.length === 1) {
        // just what we want
        callback(null, items[0]);
        return;
      }

      if(items.length > 1) {
        // oops! how come we got more than one?
        callback("Error: more than one record found for url [" + url + "]");
        return;
      }

      // no record found
      callback(null, null);

    });

  });

};


exports.updateOne = function(data, callback) {

  mongoConnect.execute(function(err, db) {

    fetchOne(data.key, function(err, item) {

      if(err) return callback(err);
      if(item === null) return callback("Item to update was not found for key [" + data.key + "]");

      // set the _id to match the one already on the database 
      data._id = item._id;

      // calculate the date fields
      data.createdOn = item.createdOn;
      data.updatedOn = new Date();
      data.postedOn = item.postedOn ? item.postedOn : null;

      var collection = db.collection(dbCollection);
      collection.save(data, function(err, savedCount){

        if(err) return callback(err);
        if(savedCount == 0) return callback("No document was updated");
        if(savedCount > 1) return callback("More than one document was updated");

        fetchOne(data.key, function(err, item) {
          if(err) return callback(err);
          callback(null, item);
        });

      });

    });

  });
  
};


exports.addOne = function(data, callback) {

  mongoConnect.execute(function(err, db) {

    fetchOne(data.key, function(err, item) {

      if(err) return callback(err);
      if(item !== null) return callback("An item with the same key already exists [" + data.key + "]");

      // automatically calculate the date fields
      data.createdOn = new Date();
      data.updatedOn = null;
      data.postedOn = null;

      var collection = db.collection(dbCollection);
      collection.save(data, function(err, savedCount){

        if(err) return callback(err);
        callback(null, savedCount);

      });

    });

  });

};


exports.markAsDraft = function(key, callback) {

  mongoConnect.execute(function(err, db) {

    fetchOne(key, function(err, item) {

      if(err) return callback(err);
      if(item === null) return callback("Item to mark as draft was not found for key [" + key + "]");

      var query = {key: key};
      var data = {'$set' : {postedOn: null}};
   
      var collection = db.collection(dbCollection);
      collection.update(query, data, function(err){
        if(err) return callback(err);
        callback(null);
      });

    });

  });

};


exports.markAsPosted = function(key, callback) {

  mongoConnect.execute(function(err, db) {

    fetchOne(key, function(err, item) {

      if(err) return callback(err);
      if(item === null) return callback("Item to mark as posted was not found for key [" + key + "]");

      var query = {key: key};
      var postedOn = new Date();
      var data = {'$set' : {postedOn: postedOn}};
   
      var collection = db.collection(dbCollection);
      collection.update(query, data, function(err){
        if(err) return callback(err);
        callback(null, postedOn);
      });

    });

  });

};

