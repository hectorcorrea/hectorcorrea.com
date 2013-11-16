// Fixes bad encoding of "less than" and "greater than"
// characters that was messed up during the migration.

var fs = require('fs');
var MongoClient = require('mongodb').MongoClient;
var mongoUrl = 'mongodb://localhost:27017/hectorcorrea';

console.log('Connecting...');
MongoClient.connect(mongoUrl, {}, function(err, db) {

  if(err) {
    console.log('Error connecting: ' + err);
    return;
  }

  var blog = db.collection("blog");
  var query = {};
  blog.find(query).toArray(function(err, docs) {
    
    var i, doc, isBad, newText, queryOne, updatedData;
    for(i = 0; i < docs.length; i++) {

      doc = docs[i];
      isBad = doc.text.indexOf("&amp;lt;") > -1 ||
        doc.text.indexOf("&amp;gt;") > -1;

      if(isBad) {
        console.log("Doc %s / %s needs to be fixed", doc.key, doc.title);

        newText = doc.text.replace(/&amp;lt;/g, "&lt;");
        newText = newText.replace(/&amp;gt;/g, "&gt;");

        queryOne = {key: doc.key};
        updatedData = {'$set' : {text: newText}};
        blog.update(queryOne, updatedData, function(err){
          if(err) {
            console.log('ERROR: ' + err);
          }
          else {
            console.log('Updated OK');
          }
        });

        // console.log("ORIGINAL");
        // console.log("--------");
        // console.log(doc.text);
        // newText = doc.text.replace(/&amp;lt;/g, "&lt;");
        // newText = newText.replace(/&amp;gt;/g, "&gt;");
        // console.log("NEW TEXT");
        // console.log("--------");
        // console.log(newText);
      }
      else {
        // console.log("Doc %s / %s is OK", doc.key, doc.title);
      }

    }

  });
  console.log('Connected!');

});

