// Find relative links (e.g. somepath/somefile.ext)
var fs = require('fs');
var MongoClient = require('mongodb').MongoClient;
var mongoUrl = 'mongodb://localhost:27017/hectorcorrea';

// http://stackoverflow.com/a/3410557/446681
function getIndicesOf(searchStr, str, caseSensitive) {
  var startIndex = 0, searchStrLen = searchStr.length;
  var index, indices = [];
  if (!caseSensitive) {
      str = str.toLowerCase();
      searchStr = searchStr.toLowerCase();
  }
  while ((index = str.indexOf(searchStr, startIndex)) > -1) {
      indices.push(index);
      startIndex = index + searchStrLen;
  }
  return indices;
}

console.log('Connecting...');
MongoClient.connect(mongoUrl, {}, function(err, db) {

  if(err) {
    console.log('Error connecting: ' + err);
    return;
  }

  var blog = db.collection("blog");
  var query = {};
  blog.find(query).toArray(function(err, docs) {

    var i, doc, text, indices, j, matches, pos, link;
    for(i = 0; i < docs.length; i++) {

      doc = docs[i];
      text = doc.text;
      indices = getIndicesOf("href", text);
      matches = 0;
      for(j = 0; j < indices.length; j++) {
                
        pos = indices[j];
        link = text.substring(pos, pos+10);
        if (link == "href=\"http") {
          // console.log('Valid at %s found %s', pos, text.substring(pos, pos+20));
        }
        else {
          matches++;
          if(matches == 1) {
            console.log('Blog: %s %s', doc.key, doc.title);
          }
          console.log('\tInvalid link at %s: %s', pos, text.substring(pos, pos+60));          
        } 

      }

    }

  });
  console.log('Connected!');

});

