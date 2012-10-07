fs = require 'fs'
{TopicMeta}  = require './topicMeta'

# This class handles saving and retrieving of topic
# data. Notice that this class does NOT perform any
# validation, it trusts that the caller passes only
# good data. 
#
# This class is akin to a database table. It fetches
# and saves data but performs no validation.

class TopicData

  constructor: (options) ->

    if typeof options isnt "object"
      console.dir options
      throw "Invalid options parameter received" 

    if typeof options.dataPath isnt "string"
      throw "Options object does not have string dataPath property"

    @dataPath = options.dataPath
    @blogListFilePath = "#{options.dataPath}/blogs.json"
    @nextId = null
    @showDrafts = false
    if options.showDrafts? 
      @showDrafts = options.showDrafts is true

    if options.createDataFileIfNotFound is true
      @_createDataFileIfNotAvailable()


  _createDataFileIfNotAvailable: =>
    if fs.existsSync(@blogListFilePath) is false
      @nextId = 1
      topics = []
      @_saveMetaToDisk topics 


  # TODO: make sort optional
  _loadAll: (callback) =>
    fs.readFile @blogListFilePath, 'utf8', (err, text) =>

      if err
        callback "Error reading topics: #{err}"
      else

        try
          data = JSON.parse text
          @nextId = data.nextId

          topics = []
          for topic in data.blogs
            meta = new TopicMeta()
            meta.id = topic.id
            meta.title = topic.title
            meta.url = topic.url
            meta.summary = topic.summary
            meta.createdOn = new Date(topic.createdOn)
            meta.updatedOn = new Date(topic.updatedOn)
            
            isPosted = topic.postedOn?
            if isPosted
              meta.postedOn = new Date(topic.postedOn)
            
            if @showDrafts
              topics.push meta
            else if isPosted
              topics.push meta 

          # Sort most recently published topics on top
          # Force un-published (draft) topics to the top
          topics.sort (x, y) ->
            xDate = if x.postedOn is null then new Date() else x.postedOn
            yDate = if y.postedOn is null then new Date() else y.postedOn
            return -1 if xDate > yDate
            return 1 if xDate < yDate
            return 0 

          callback null, topics

        catch error
          callback "Error parsing data file. Error: #{error}"


  # This method is sync on purpose. We don't want anyone
  # else to read/write the file while we are updating it.
  _saveMetaToDisk: (topics) =>
    jsonText = JSON.stringify topics, null, "\t"
    jsonText = '{ "nextId": ' + @nextId + ', "blogs":' + jsonText + '}'
    fs.writeFileSync @blogListFilePath, jsonText, 'utf8'      


  getAll: (callback) =>
    @_loadAll callback


  getRecent: (callback) =>
    howMany = 10
    @_loadAll (err, topics) =>
      if err 
        callback err
      else
        if topics.length > howMany
          topics = topics.slice(0, howMany)
        callback null, topics 


  getNew: =>
    meta = new TopicMeta()
    return {meta: meta, content: ""}


  findMeta: (id, callback) =>
    @_loadAll (err, topics) => 
      if err
        callback err
      else
        err = "Topic id #{id} not found"
        topic = null
        for t in topics
          if t.id is id
            err = null
            topic = t
            break
        callback err, topic


  findMetaByUrl: (url, callback) =>
    @_loadAll (err, topics) =>
      if err
        callback err
      else
        err = "Topic url #{url} not found"
        topic = null
        for t in topics
          if t.url is url
            err = null
            topic = t
            break
        callback err, topic


  loadContent: (meta, callback) => 
    filePath = @dataPath + '/blog.' + meta.id + '.html'  
    fs.readFile filePath, 'utf8', (err, text) =>
      if err 
        callback "Could not retrieve content for id #{meta.id}"
      else
        # console.log "LOAD: ------------------"
        # console.log text
        callback null, {meta: meta, content: text}


  updateMeta: (id, newMeta, callback) =>
    @_loadAll (err, topics) => 
      if err
        callback err
      else
        err = "Topid id #{id} not found"
        topic = null
        for i in [0..topics.length-1]
          if topics[i].id is id
            # Update the meta data...
            topics[i].title = newMeta.title
            topics[i].url = newMeta.url
            topics[i].summary = newMeta.summary
            topics[i].createdOn = newMeta.createdOn
            topics[i].updatedOn = newMeta.updatedOn
            topics[i].postedOn = newMeta.postedOn
            @_saveMetaToDisk topics
            # ...and return the updated topic
            err = null
            topic = topics[i]
            break
        callback err, topic


  updateContent: (meta, content, callback) =>
    filePath = @dataPath + '/blog.' + meta.id + '.html'  
    fs.writeFile filePath, content, 'utf8', (err) => 
      if err 
        callback "Content for topic #{meta.id} could not be saved. Error #{err}"
      else
        callback null, {meta: meta, content: content}


  addNew: (meta, content, callback) =>
    @_loadAll (err, topics) => 
      if err
        callback err
      else
        meta.id = @nextId 
        topics.push meta
        @nextId++
        @_saveMetaToDisk topics
        @updateContent meta, content, callback


exports.TopicData = TopicData
