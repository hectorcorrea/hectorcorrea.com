var db = require('./blogDb');
var util = require('./encodeUtil');


var searchText = function(text) {
  if (typeof(text) === 'string') {
    // http://stackoverflow.com/a/9364527/446681
    return text.replace(/\W/g, ' ').toLowerCase().trim();
  }
  return '';
};


var prepareForSave = function(rawData) {
  data = {}
  data.key = rawData.key;
  data.title = rawData.title;
  data.summary = rawData.summary;
  data.text = rawData.text;

  // calculate a few fields
  data.url = util.urlSafe(data.title);
  data.searchText = searchText(data.title) + ' ' + 
    searchText(data.text) + ' '
    searchText(data.summary);

  return data;
};


exports.getAll = function(includeDrafts, cb) {
  db.fetchAll(includeDrafts, function(err, documents) {
    cb(err, documents);
  });
};


exports.getOne = function(key, _notUsed, cb) {
  db.fetchOne(key, function(err, document) {
    cb(err, document);
  });
};


exports.getOneByUrl = function(url, _notUsed, cb) {
  db.fetchOneByUrl(url, function(err, document) {
    cb(err, document);
  });
};


exports.addNew = function(cb) {

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
  data = prepareForSave(data);
  db.updateOne(data, function(err, savedDoc) {
    if (err) return cb(err);
    cb(null, savedDoc);
  });
};


exports.markAsDraft = function(key, cb) {
  db.markAsDraft(key, function(err) {
    if (err) return cb(err);
    cb(null);
  });
};


exports.markAsPosted = function(key, cb) {
  db.markAsPosted(key, function(err, postedOn) {
    if (err) return cb(err);
    cb(null, postedOn);
  });
};


