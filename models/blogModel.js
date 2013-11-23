var db = require('./blogDb');
var util = require('./encodeUtil');
var dbUrl = null;


var _searchText = function(text) {
  if (typeof(text) === 'string') {
    // http://stackoverflow.com/a/9364527/446681
    return text.replace(/\W/g, ' ').toLowerCase().trim();
  }
  return '';
};


var _prepareForSave = function(rawData) {

  // encode data first
  data = {}
  data.key = rawData.key;
  data.title = util.encodeText(rawData.title);
  data.summary = rawData.summary; // util.encodeText(rawData.summary);
  data.text = rawData.text;       // util.encodeText(rawData.text);

  // calculate a few fields
  data.url = util.urlSafe(data.title);
  data.searchText = _searchText(data.title) + ' ' + 
    _searchText(data.text) + ' '
    _searchText(data.summary);

  return data;
};


var _decodeForEdit = function(data) {
  data.title = util.decodeText(data.title);
  data.summary = data.summary;    // util.decodeText(data.summary);
  data.text = data.text;          // util.decodeText(data.text);
  return data;
};


var getAll = function(includeDrafts, cb) {
  db.setup(dbUrl);
  db.fetchAll(includeDrafts, function(err, documents) {
    cb(err, documents);
  });
};


var getOne = function(key, decode, cb) {
  db.setup(dbUrl);
  db.fetchOne(key, function(err, document) {
    if(decode) {
      document = _decodeForEdit(document);
    }
    // console.dir(document);
    cb(err, document);
  });
};


var getOneByUrl = function(url, decode, cb) {
  db.setup(dbUrl);
  db.fetchOneByUrl(url, function(err, document) {
    if(decode) {
      document = _decodeForEdit(document);
    }
    // console.dir(document);
    cb(err, document);
  });
};


var addNew = function(cb) {

  db.setup(dbUrl);
  db.getNewId(function(err, id) {

    if(err) return cb(err);
    
    var data = {
      key: id,
      title: 'New blog post',
      text: '',
      summary: ''
    };
    data = _prepareForSave(data);

    db.addOne(data, function(err, savedDoc) {
      cb(err, savedDoc);
    });

  });

};


var updateOne = function(data, cb) {
 
  db.setup(dbUrl);
  data = _prepareForSave(data);
  db.updateOne(data, function(err, savedDoc) {
    if (err) return cb(err);
    cb(null, savedDoc);
  });

};


var markAsDraft = function(key, cb) {

  db.setup(dbUrl);
  db.markAsDraft(key, function(err) {
    if (err) return cb(err);
    cb(null);
  });

};


var markAsPosted = function(key, cb) {

  db.setup(dbUrl);
  db.markAsPosted(key, function(err, postedOn) {
    if (err) return cb(err);
    cb(null, postedOn);
  });

};


var publicApi = {
  getAll: getAll,
  getOne: getOne,
  getOneByUrl: getOneByUrl,
  addNew: addNew,
  updateOne: updateOne,
  markAsDraft: markAsDraft,
  markAsPosted: markAsPosted
};


module.exports.blog = function(dbConnString) {
  dbUrl = dbConnString;
  return publicApi;
};

