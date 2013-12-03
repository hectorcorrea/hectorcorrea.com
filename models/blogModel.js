var db = require('./blogDb');
var util = require('./encodeUtil');
var dbUrl = null;


var searchText = function(text) {
  if (typeof(text) === 'string') {
    // http://stackoverflow.com/a/9364527/446681
    return text.replace(/\W/g, ' ').toLowerCase().trim();
  }
  return '';
};


var prepareForSave = function(rawData) {

  // encode data first
  data = {}
  data.key = rawData.key;
  data.title = util.encodeText(rawData.title);
  data.summary = rawData.summary; // util.encodeText(rawData.summary);
  data.text = rawData.text;       // util.encodeText(rawData.text);

  // calculate a few fields
  data.url = util.urlSafe(data.title);
  data.searchText = searchText(data.title) + ' ' + 
    searchText(data.text) + ' '
    searchText(data.summary);

  return data;
};


var decodeForEdit = function(data) {
  data.title = util.decodeText(data.title);
  data.summary = data.summary;    // util.decodeText(data.summary);
  data.text = data.text;          // util.decodeText(data.text);
  return data;
};


exports.getAll = function(includeDrafts, cb) {
  // db.setup(dbUrl);
  db.fetchAll(includeDrafts, function(err, documents) {
    cb(err, documents);
  });
};


exports.getOne = function(key, decode, cb) {
  // db.setup(dbUrl);
  db.fetchOne(key, function(err, document) {
    if(decode) {
      document = decodeForEdit(document);
    }
    cb(err, document);
  });
};


exports.getOneByUrl = function(url, decode, cb) {
  // db.setup(dbUrl);
  db.fetchOneByUrl(url, function(err, document) {
    if(decode) {
      document = decodeForEdit(document);
    }
    cb(err, document);
  });
};


exports.addNew = function(cb) {

  // db.setup(dbUrl);
  db.getNewId(function(err, id) {

    if(err) return cb(err);
    
    var data = {
      key: id,
      title: 'New blog post',
      text: '',
      summary: ''
    };
    data = prepareForSave(data);

    db.addOne(data, function(err, savedDoc) {
      cb(err, savedDoc);
    });

  });

};


exports.updateOne = function(data, cb) {
 
  // db.setup(dbUrl);
  data = prepareForSave(data);
  db.updateOne(data, function(err, savedDoc) {
    if (err) return cb(err);
    cb(null, savedDoc);
  });

};


exports.markAsDraft = function(key, cb) {

  // db.setup(dbUrl);
  db.markAsDraft(key, function(err) {
    if (err) return cb(err);
    cb(null);
  });

};


exports.markAsPosted = function(key, cb) {

  // db.setup(dbUrl);
  db.markAsPosted(key, function(err, postedOn) {
    if (err) return cb(err);
    cb(null, postedOn);
  });

};


// exports.setup = function(dbConnString) {
//   dbUrl = dbConnString;
// };

