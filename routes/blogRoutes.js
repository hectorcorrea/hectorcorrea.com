var RSS = require('rss');
var logger = require('log-hanging-fruit').defaultLogger;
var model = require('../models/blogModel');

var notFound = function(req, res, key) {
  logger.warn('Blog entry not found. Key [' + key + ']');
  req.app.settings.setCache(res, 5);
  res.status(404).render('notFound');
};


var notAuthenticated = function(req, res, method) {
  logger.error(method + ' User is not authenticated');
  res.status(401).render('error', {error: 'User is not authenticated'});
};


var redirectToView = function(req, res, blog) {
  logger.error('Redirecting to ' + blog.url + '/' + blog.key);
  res.redirect(301, '/blog/' + blog.url + '/' + blog.key);
}


var error = function(req, res, title, err) {
  logger.error(title + ' ' + err);
  res.status(500).render('error', {title: title, error: err});
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
      postedOn: doc.postedOn,
      isDraft: (doc.postedOn == null)
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
    postedOn: doc.postedOn,
    isDraft: (doc.postedOn == null)
  };
  return json
};


exports.viewAll = function(req, res) {

  logger.info('blog.all');
  var includeDrafts = req.isAuth;
  model.getAll(includeDrafts, function(err, documents){

    if(err) {
      return error(req, res, "Cannot retrieve all blog entries", err);
    }

    var blogs = docsToJson(documents);
    req.app.settings.setCache(res, 5);
    res.render('blogList', {blogs: blogs, isAuth: req.isAuth})
  });

};


exports.viewOne = function(req, res) {

  var key = parseInt(req.params.key)
  var url = req.params.url;

  logger.info('blog.blogView (' + key + ', ' + url + ')');
  model.getOne(key, null, function(err, doc){

    if(err) {
      return error(req, res, 'Error fetching blog [' + key + ']', err);
    }

    if(doc === null) {
      return notFound(req, res, key);
    }

    var blog = docToJson(doc);
    req.app.settings.setCache(res, 5);
    res.render('blogView', {blog: blog, isAuth: req.isAuth})
  });

};


exports.edit = function(req, res) {

  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.edit');
  }

  var key = parseInt(req.params.key)
  var url = req.params.url;

  logger.info('blog.edit (' + key + ', ' + url + ')');
  model.getOne(key, null, function(err, doc){

    if(err) {
      return error(req, res, 'Error fetching blog [' + key + ']', err);
    }

    if(doc === null) {
      return notFound(req, res, key);
    }

    var blog = docToJson(doc);
    console.log(blog)
    req.app.settings.setCache(res, 5);
    res.render('blogEdit', {blog: blog, isAuth: req.isAuth})
  });

};


exports.post = function(req, res) {
  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.post');
  }

  var key = parseInt(req.params.key)
  var url = req.params.url;

  logger.info('blog.post (' + key + ', ' + url + ')');
  model.markAsPosted(key, function(err){

    if(err) {
      return error(req, res, 'Error posting blog [' + key + ']', err);
    }

    var blog = docToJson({key: key, url: url});
    req.app.settings.setCache(res, 5);
    redirectToView(req, res, blog);
  });
}

exports.draft = function(req, res) {
  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.draft');
  }

  var key = parseInt(req.params.key)
  var url = req.params.url;

  logger.info('blog.draft (' + key + ', ' + url + ')');
  model.markAsDraft(key, function(err){

    if(err) {
      return error(req, res, 'Error marking as draft blog [' + key + ']', err);
    }

    var blog = docToJson({key: key, url: url});
    req.app.settings.setCache(res, 5);
    redirectToView(req, res, blog);
  });
};


exports.newBlog = function(req, res) {

  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.newOne');
  }

  logger.info('blog.new');
  model.addNew(function(err, newDoc){

    if(err) {
      return error(req, res, 'Error adding new blog', err);
    }

    var blog = docToJson(newDoc);
    res.render('blogEdit', {blog: blog, isAuth: true})
  });

};


exports.save = function(req, res) {

  if(!req.isAuth) {
    return notAuthenticated(req, res, 'blog.save');
  }

  logger.info('blog.save');

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

  model.updateOne(data, function(err, savedDoc){

    if(err) {
      return error(req, res, 'Error saving blog [' + data.key + ']', err);
    }

    var blog = docToJson(savedDoc);
    console.log(blog);
    redirectToView(req, res, blog);
  });

};


exports.rss = function(req, res) {

  logger.info('blog.rss');

  var includeDrafts = false;
  model.getAll(includeDrafts, function(err, documents){

    if(err) {
      return error(req, res, "Cannot retrieve all blog entries for RSS feed", err);
    }

    var options = req.app.settings.config.rss;
    var rootUrl = options.site_url;
    var feed = new RSS(options);
    var i, entry, doc;
    for(i=0; i<documents.length; i++) {
      doc = documents[i];
      entry = {
        title: doc.title,
        description: doc.summary,
        url: rootUrl + '/blog/' + doc.url + '/' + doc.key,
        date: doc.postedOn
      };
      feed.item(entry);
    }

    req.app.settings.setCache(res, 60);
    res.send(feed.xml());
  });

};


