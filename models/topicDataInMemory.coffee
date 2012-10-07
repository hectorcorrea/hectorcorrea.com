{TopicMeta}  = require './topicMeta'

# This version of TopicData stores data in memory and
# it's only good as a mock for unit tests. This class
# is NOT meant to be used in production.
#
# This class handles saving and retrieving of topic
# data. Notice that this class does NOT perform any
# validation, it trusts that the caller passes only
# good data. 

class TopicDataInMemory

  @topics = null
  @contents = null
  @nextId = null

  constructor: ->
    @topics = []
    @contents = []
    @nextId = null
  

  _loadAll: =>
    console.log "-- In memory version --"
    if @topics.length is 0
      for i in [1..3]
        topic = new TopicMeta()
        topic.id = i
        topic.title = "topic #{i}"
        topic.url = "topic-#{i}"
        topic.summary = "topic #{i} summary"
        topic.createdOn = new Date("Jan #{i}, 2012")
        topic.updatedOn = new Date("Feb #{i}, 2012")
        topic.postedOn = new Date("March #{i}, 2012")
        @topics.push topic
        @contents.push "content for topic #{i}"
      @nextId = 4
    @topics


  getAll: =>
    return @_loadAll()


  getRecent: =>
    return @_loadAll().slice(0,2)


  getNew: =>
    meta = new TopicMeta()
    return {meta: meta, content: ""}


  findMeta: (id) =>
    topics = @_loadAll()
    for topic in topics
      if topic.id is id
        return topic
    return null


  findMetaByUrl: (url) =>
    topics = @_loadAll()
    for topic in topics
      if topic.url is url
        return topic
    return null


  loadContent: (meta, callback) => 
    process.nextTick =>
      content = @contents[meta.id-1]
      if content
        callback null, {meta: meta, content: content}
      else 
        callback "Invalid ID #{meta.id}"


  updateMeta: (id, newMeta) =>
    topic = @findMeta(id)
    return null if topic is null 
    topic.title = newMeta.title
    topic.url = newMeta.url
    topic.summary = newMeta.summary
    topic.createdOn = newMeta.createdOn
    topic.updatedOn = newMeta.updatedOn
    topic.postedOn = newMeta.postedOn
    return topic


  updateContent: (meta, content, callback) =>
    maxId = @nextId-1
    if meta.id? and meta.id in [1..maxId]
      process.nextTick =>
        @contents[meta.id] = content
        callback null, {meta: meta, content: content}
    else
      process.nextTick =>
        callback "Invalid id #{meta.id}"


  addNew: (meta, content, callback) =>
    @_loadAll() # Forces @nextId to be initialized

    meta.id = @nextId 
    @topics.push meta
    process.nextTick =>
      @contents.push content
      @nextId++
      callback null, {meta: meta, content: content}

exports.TopicData = TopicDataInMemory
