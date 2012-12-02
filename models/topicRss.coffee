events = require 'events'
fs = require 'fs'
RSS = require 'rss'

class TopicRss

  emitter = null
  processedCount = 0
  totalCount = 0
  feed = null
  callback = null
  errors = []

  constructor: (@settings, @pathToFiles, @topics) ->
    # Create the RSS feed
    @feed = new RSS(@settings)

    # Wire up the event emitter
    @emitter = new events.EventEmitter()
    @emitter.on 'createRssFeed', @_onCreateRssFeed
    @emitter.on 'addFeedItem', @_onAddFeedItem
    @emitter.on 'done', @_onDone


  _onCreateRssFeed: () =>
    @errors = []
    @processedCount = 0
    @totalCount = @topics.length
    for topic in @topics
      @emitter.emit 'addFeedItem', topic


  _onAddFeedItem: (topic) =>
    @processedCount++
    isLastItem = if @processedCount is @totalCount then true else false
    fileName = @pathToFiles + "/blog.#{topic.id}.html"
    fs.readFile fileName, 'utf8', (err, text) =>

      if err
        errors.push err
      else
        entry = {
          title: topic.title
          description: text
          url: @feed.site_url + '/blog/' + topic.url
          date: topic.postedOn
        }
        @feed.item entry

      @emitter.emit 'done' if isLastItem


  _onDone: =>
    if errors.length > 0
      @callback errors
    else
      @callback null, @feed.xml()


  createRssFeed: (callback) =>
    @callback = callback
    @emitter.emit('createRssFeed')


exports.TopicRss = TopicRss
  
