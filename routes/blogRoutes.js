var RSS = require('rss');

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
