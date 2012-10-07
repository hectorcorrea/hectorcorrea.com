class TopicMeta
  constructor: ->
    @id = NaN
    @title = ""
    @createdOn = new Date()
    @updatedOn = new Date()
    @postedOn = null
    @url = ""
    @summary = ""

exports.TopicMeta = TopicMeta
