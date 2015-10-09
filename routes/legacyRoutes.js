var logger = require('log-hanging-fruit').defaultLogger;
var model = require('../models/blogModel');
var util = require('../models/encodeUtil');

exports.blogOne = function(req, res) {

  var url = util.urlSafe(req.params.url);
  var legacyExt = /-aspx$/.test(url)
  if(legacyExt) {
    url = url.substring(0, url.length-5);
  }
  var decode = req.query.decode === "true";

  logger.info('legacy.blogOne (' + url + ')');
  model.getOneByUrl(url, decode, function(err, doc){

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

    res.redirect(301, '/blog/' + doc.url + '/' + doc.key);
  });

};

