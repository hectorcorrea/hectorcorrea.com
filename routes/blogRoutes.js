var model = require('../models/blogModel');
var logger = require('log-hanging-fruit').defaultLogger;


var notFound = function(req, res, key) {
  logger.warn('Blog entry not found. Key [' + key + ']');
  req.app.settings.setCache(res, 5);
  res.status(404).send({message: 'Blog entry not found' });
};


var notAuthenticated = function(req, res, method) {
  logger.error(method + ' User is not authenticated');
  res.status(401).send('User is not authenticated.');
};


var error = function(req, res, title, err) {
  logger.error(title + ' ' + err);
  res.status(500).send({message: title, details: err});
};


var docsToJson = function(documents) {
  var json = [];
  var i, blog, doc; 
  for(i=0; i<documents.length; i++) {
    doc = documents[i];
    // we don't include the text on purpose
    blog = {
      key: doc.key,
      title: doc.title,
      url: doc.url,
      summary: doc.summary,
      postedOn: doc.postedOn
    }
    json.push(blog);
  }
  return json;
}


var docToJson = function(doc) {
  var json = {
    key: doc.key,
    title: doc.title,
    url: doc.url,
    text: doc.text,
    summary: doc.summary,
    createdOn: doc.createdOn,
    updatedOn: doc.updatedOn,
    postedOn: doc.postedOn
  };
  return json
};


exports.all = function(req, res) {

  logger.info('blog.all');

  var m = model.blog(req.app.settings.config.dbUrl);
  var includeDrafts = req.isAuth;
  m.getAll(includeDrafts, function(err, documents){

    if(err) {
      return error(req, res, "Cannot retrieve all blog entries", err);
    }

    var blogs = docsToJson(documents);
    req.app.settings.setCache(res, 5);
    res.send(blogs);
  });

};


exports.one = function(req, res) {

  var key = parseInt(req.params.key)
  var url = req.params.url;
  var decode = req.query.decode === "true";
  var m = model.blog(req.app.settings.config.dbUrl);

  logger.info('blog.one (' + key + ', ' + url + ')');
  m.getOne(key, decode, function(err, doc){

    if(err) {
      return error(req, res, 'Error fetching blog [' + key + ']', err);
    }

    if(doc === null) {
      return notFound(req, res, key);
    }

    var blog = docToJson(doc);
    req.app.settings.setCache(res, 5);    
    res.send(blog);
  });

};


exports.draft = function(req, res) {

  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.draft');
  }

  var key = parseInt(req.params.key)
  var url = req.params.url;
  var decode = false;
  var m = model.blog(req.app.settings.config.dbUrl);

  logger.info('blog.draft (' + key + ', ' + url + ')');
  m.markAsDraft(key, function(err){

    if(err) {
      return error(req, res, 'Error marking as draft blog [' + key + ']', err);
    }

    var blog = docToJson({key: key});
    res.send(blog);
  });

};


exports.post = function(req, res) {

  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.post');
  }

  var key = parseInt(req.params.key)
  var url = req.params.url;
  var decode = false;
  var m = model.blog(req.app.settings.config.dbUrl);

  logger.info('blog.post (' + key + ', ' + url + ')');
  m.markAsPosted(key, function(err, postedOn){

    if(err) {
      return error(req, res, 'Error marking as posted blog [' + key + ']', err);
    }

    var blog = docToJson({key: key, postedOn: postedOn});
    res.send(blog);
  });

};


exports.newOne = function(req, res) {

  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.newOne');
  }

  var m = model.blog(req.app.settings.config.dbUrl);

  logger.info('blog.new');
  m.addNew(function(err, newDoc){

    if(err) {
      return error(req, res, 'Error adding new blog', err);
    }

    var blog = docToJson(newDoc);
    res.send(blog);
  });

};


exports.save = function(req, res) {

  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.save');
  }

  logger.info('blog.save');

  var m = model.blog(req.app.settings.config.dbUrl);

  var data = {
    key: parseInt(req.params.key, 10),
    title: req.body.title,
    summary: req.body.summary,
    text: req.body.text,
  };  

  logger.info('blog.save (' + data.key + ')');

  if(data.title === '') {
    return error(req, res, 'Blog title cannot be empty', 'key: ' + data.key);
  }

  if(data.text === '') {
    return error(req, res, 'Blog text cannot be empty', 'key: ' + data.key);
  }
  
  if(data.summary === '') {
    return error(req, res, 'Blog summary cannot be empty', 'key: ' + data.key);
  }

  m.updateOne(data, function(err, savedDoc){

    if(err) {
      return error(req, res, 'Error saving blog [' + data.key + ']', err);
    }

    var blog = docToJson(savedDoc);
    res.send(blog);
  });

};

