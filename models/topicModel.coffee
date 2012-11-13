{TopicData}  = require './topicData'
{TopicRss}  = require './topicRss'

class TopicModel
  
  @data = null

  constructor: (dataOptions) ->
    @data = new TopicData(dataOptions)


  _getUrlFromTitle: (title) ->
    title = title.trim()
    title = title.toLowerCase()
    title = title.replace('c#', 'csharp')

    url = ""
    for i in [0..title.length-1]
      c = title[i]
      if c >= 'a' and c <= 'z'
        url += c
      else if c >= '0' and c <= '9'
        url += c
      else if c is '(' or c is ')' 
        url += c
      else
        url += '-'

    while url.indexOf('--') > -1
      url = url.replace('--', '-')

    return url    


  _isValidDate: (date) ->
    return date instanceof Date && isFinite(date)


  _validateTopic: (topic, callback) =>

    valid = true
    errors = {
      emptyTitle: false
      duplicateTitle: false
      emptySummary: false
      emptyContent: false
    }
    
    if topic?.meta?.title? is false or topic.meta.title is ""
      valid = false
      errors.emptyTitle = true
    if topic?.meta?.summary? is false or topic.meta.summary is ""
      valid = false
      errors.emptySummary = true
    if topic?.content? is false or topic.content is ""
      valid = false
      errors.emptyContent = true

    if errors.emptyTitle
      # No need to validate for a duplicate title
      callback null, errors
    else
      @_isDuplicateTitle topic, (err, isDuplicate) =>
        if err
          callback err
        else
          if isDuplicate
            valid = false
            errors.duplicateTitle = true
          callback null, if valid then null else errors


  _isDuplicateTitle: (topic, callback) =>

    valid = true
    url = @_getUrlFromTitle(topic.meta.title)
    isNewTopic = isNaN(topic.meta.id)
    
    if isNewTopic
      # When validating a new topic, finding any 
      # topic with the same title is enough to 
      # consider it a duplicate
      @data.findMetaByUrl url, (err, topicFound) =>
        isDuplicate = topicFound isnt null
        callback null, isDuplicate
    else
      # When updating an existing topic we consider 
      # a duplicate only if the topic found is NOT 
      # the same as the one we are trying to save.
      @data.findMetaByUrl url, (err, topicFound) =>
        if topicFound is null
          callback null, false
        else
          isDuplicate = topicFound.id isnt topic.meta.id
          callback null, isDuplicate  


  getAll: (callback) =>
    @data.getAll callback


  getRecent: (callback) =>
    @data.getRecent callback


  getOne: (id, callback) =>
    @data.findMeta id, (err, meta) =>
      if err 
        callback err
      else
        @data.loadContent meta, callback


  getOneByUrl: (url, callback) =>
    @data.findMetaByUrl url, (err, meta) =>
      if err 
        callback err
      else
        @data.loadContent meta, callback


  getNew: =>
    return @data.getNew()


  getRssList: (callback) =>
    @data.getAll (err, topics) =>
      if err
        callback err
      else
        rss = new TopicRss()
        rss.toRss topics, callback

  # topic must be in the form 
  # {meta: {id: i, title: t, summary: s, ...}, content: c}
  # notice that we need an id
  save: (topic, callback) =>

    # Load the topic from the DB...
    @data.findMeta topic.meta.id, (err, meta) => 

      if err
        callback err
      else
        # ...merge the topic that we received with the one 
        # on the DB
        topic.meta.createdOn = meta.createdOn
        topic.meta.updatedOn = new Date()
        topic.meta.url = @_getUrlFromTitle(topic.meta.title)
        if topic.meta.postedOn is null
          # preserve the original postedOn date
          topic.meta.postedOn = meta.postedOn

        # ...make sure the topic is valid
        @_validateTopic topic, (err, validationErrors) =>
          if err
            callback err
          else
            isTopicValid = validationErrors is null
            if isTopicValid
              # ...update the meta data
              @data.updateMeta topic.meta.id, topic.meta, (err, updatedMeta) =>
                if err
                  callback err
                else  
                  # ... and the content
                  @data.updateContent updatedMeta, topic.content, callback
            else
              # topic has {meta: X, content: Y, errors: Z}
              topic.errors = validationErrors
              callback null, topic


  # topic must be in the form 
  # {meta: {title: t, summary: s, ...}, content: c}
  # notice that we don't need an id
  saveNew: (topic, callback) => 
    # Fill in values required for new topics
    topic.meta.createdOn = new Date()
    topic.meta.updatedOn = new Date()
    #topic.meta.postedOn = null if @_isValidDate(topic.meta.postedOn) is false
    topic.meta.url = @_getUrlFromTitle(topic.meta.title)

    # ...make sure the topic is valid
    @_validateTopic topic, (err, validationErrors) =>
      if err
        callback err
      else
        isTopicValid = validationErrors is null
        if isTopicValid
          # Add topic to the database (meta+content)
          @data.addNew topic.meta, topic.content, callback
        else
            # topic has {meta: X, content: Y, errors: Z}
          topic.errors = validationErrors
          callback null, topic


exports.TopicModel = TopicModel
