var logger = require('log-hanging-fruit').defaultLogger;
var model = require('../models/blogModel');
var util = require('../models/encodeUtil');

var about = function(req, res) {
  logger.warn('legacy.about');
  res.status(301).redirect('/#/about');
};


var blogAll = function(req, res) {
  logger.warn('legacy.blogAll');
  res.status(301).redirect('/#/blog');
};


var blogOne = function(req, res) {

  var url = util.urlSafe(req.params.url);
  var legacyExt = /-aspx$/.test(url)
  if(legacyExt) {
    url = url.substring(0, url.length-5);
  }
  var decode = req.query.decode === "true";
  var m = model.blog(req.app.settings.config.dbUrl);

  logger.info('legacy.blogOne (' + url + ')');
  m.getOneByUrl(url, decode, function(err, doc){

    if(err) {
      logger.error('Error fetching legacy blog [' + url + '] ' + err);
      res.status(500).send('Error fetching blog topic requested');
      return;
    }

    if(doc === null) {
      logger.error('Legacy blog not found: [' + url + ']');
      res.status(404).send('Blog topic requested was not found');
      return;
    }

    res.status(301).redirect('/#/blog/' + doc.url + '/' + doc.key);
  });

};


var rss = function(req, res) {
  res.send('<h1>RSS</h1>');
}

module.exports = {
  about: about,
  blogAll: blogAll,
  blogOne: blogOne,
  rss: rss
}
