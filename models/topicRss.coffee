RSS = require 'rss'

class TopicRss

  _getFeedSettings: =>
    # TODO: Make these settings configurable
    return {
      title: 'Hector Correa'
      description: "Hector Correa's blog"
      feed_url: 'http://simple-blog.jitsu.com/blog/rss.xml'
      site_url: 'http://simple-blog.jitsu.com'
      author: 'hector@hectorcorrea.com'
    }


  toRss: (topics, callback) =>
    settings = @_getFeedSettings()
    feed = new RSS(settings)
    for t in topics
      # TODO: Get the topic content rather than the summary
      blogEntry = {
        title: t.title
        description: t.summary
        url: settings.site_url + '/blog/' + t.url
        date: t.postedOn
      }
      feed.item blogEntry

    callback null, feed.xml()


exports.TopicRss = TopicRss
  
