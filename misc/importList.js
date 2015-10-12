// Imports a list of topics (without their content)
// from old text/json files to a MongoDB database.
var fs = require('fs');
var MongoClient = require('mongodb').MongoClient;

var sourceDir = '/Users/hector/temp/data/';
var mongoUrl = 'mongodb://localhost:27017/hectorcorrea';

var listTxt = fs.readFileSync(sourceDir + 'blogs.json');
var list = JSON.parse(listTxt);
var i, blog, data;

var importOne = function(db, data) {

  var collection = db.collection("blog");
  collection.save(data, function(err, count){

    if(err) {
      console.log('Error saving: ' + err);
      console.dir(data);
      console.log('---------------------')
      return;
    }

    console.log('Saved OK: %s %s', data.key, data.title);

  });

};


var updateNextId = function(db, nextId) {

  var counters = db.collection('counters');
  var query = {'name': 'blogId'};
  var order = [['_id','asc']];
  var inc = {$set:{'next':nextId}};
  var options = {new: true, upsert: true};
  counters.findAndModify(query, order, inc, options, function(err, doc) {

    if(err) {
      console.log('Error setting nextId counter: ' + err);
      return;
    }      

    console.log('Set nextId counter');

  });

};


console.log('Connecting...');
MongoClient.connect(mongoUrl, {}, function(err, db) {

  if(err) {
    console.log('Error connecting: ' + err);
    return;
  }

  console.log('Connected!');
  updateNextId(db, list.nextId);

  for(i=0; i < list.blogs.length; i++) {

    blog = list.blogs[i];
    data = {
      key : blog.id,
      title : blog.title,
      summary : blog.summary,
      text : '',
      url : blog.url,
      searchText : '',
      createdOn : blog.createdOn,
      updatedOn : blog.updatedOn,
      postedOn : blog.postedOn
    };

    importOne(db, data);

  }

});

