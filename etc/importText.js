// Imports the content from old text/json files 
// to a MongoDB database.
var fs = require('fs');
var MongoClient = require('mongodb').MongoClient;

var sourceDir = '/Users/hector/temp/data/';
var mongoUrl = 'mongodb://localhost:27017/hectorcorrea';

var listTxt = fs.readFileSync(sourceDir + 'blogs.json');
var list = JSON.parse(listTxt);
var i, blog, data;

var importOne = function(db, data) {

  var collection = db.collection("blog");
  var newData = {$set: {text : data.text}};
  collection.update({key: data.key}, newData, function(err, count){

    if(err) {
      console.log('Error saving: ' + err);
      console.dir(data);
      console.log('---------------------')
      return;
    }

    console.log('Saved OK: %s', data.key);

  });

};


console.log('Connecting...');
MongoClient.connect(mongoUrl, {}, function(err, db) {

  if(err) {
    console.log('Error connecting: ' + err);
    return;
  }

  console.log('Connected!');

  for(i=0; i < list.blogs.length; i++) {

    blog = list.blogs[i];
    fileName = sourceDir + 'blog.'+ blog.id + '.html';
    text = fs.readFileSync(fileName).toString()
    data = {key : blog.id, text : text};
    importOne(db, data);

  }

});

