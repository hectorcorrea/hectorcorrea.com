# Tests for topicData class.
#
# Notice that these tests must be run in order
# in order so that the first ones add data that 
# the next ones can use.

fs = require 'fs'
{TopicData}  = require './topicData'
{TestUtil}  = require '../util/testUtil'

verbose = true
test = new TestUtil("topicDataTest", verbose)
dataOptions = { 
  dataPath: __dirname + "/../data_test"
  createDataFileIfNotFound: true,
  showDrafts: true
}


# Delete the current data file (if any)
dataFile = dataOptions.dataPath + '/blogs.json'
if fs.existsSync dataFile
  fs.unlinkSync dataFile


data = new TopicData(dataOptions)


addTopicTest = ->

  newTopic = {
    title: 'new title'
    url: 'new-title'
    summary: 'new summary'
    createdOn: new Date()
    updatedOn: new Date()
    postedOn: new Date()
  }

  data.addNew newTopic, "new content", (err, newTopic)->
    if err
      test.fail "addNew (#{err})"
    else
      data.findMeta newTopic.meta.id, (err, topic) ->
        if err 
          test.fail "addNew (#{err})"
        else
          test.passIf topic.title is "new title", "addNew"
    getAllTest()


getAllTest = ->
  data.getAll (err, topics) ->
    test.passIf err is null and topics.length > 0, "getAll"
    getRecentTest()


getRecentTest = ->
  data.getRecent (err, topics) ->
    test.passIf err is null and topics.length > 0, "getRecent"
    findValidIdMetaTest()


findValidIdMetaTest = ->
  data.findMeta 1, (err, topic) ->
    test.passIf topic isnt null, "findMeta valid id"
    findInvalidIdMetaTest()


findInvalidIdMetaTest = ->
  data.findMeta -9, (err, topic) ->
    test.passIf err isnt null, "findMeta invalid id"
    findValidUrlTest()


findValidUrlTest = ->
  data.findMetaByUrl "new-title", (err, topic) ->
    test.passIf topic.url is "new-title", "findMetaByUrl valid url"
    findInvalidUrlTest()


findInvalidUrlTest = ->
  data.findMetaByUrl 'topic-not-existing', (err, topic) ->
    test.passIf topic is null, "findMetaByUrl invalid url"
    loadValidContentTest()


loadValidContentTest = ->
  data.loadContent {id: 1}, (err, data) -> 
    test.passIf err is null, "loadContent valid id"
    loadInvalidContentTest()


loadInvalidContentTest = ->
  data.loadContent {id: -9}, (err, data) -> 
    test.passIf err isnt null, "loadContent invalid id"
    updateValidTopicTest()


updateValidTopicTest = ->
  newMeta = {
    title: "updated title 1"
    url: "updated-title-1"
    summary: "updated summary 1"
  }

  data.updateMeta 1, newMeta, (err, topic) ->
    test.passIf topic.title is newMeta.title, "updateMeta valid id" 
    updateInvalidTopicTest()


updateInvalidTopicTest = ->
  newMeta = {}
  data.updateMeta -9, newMeta, (err, topic) ->
    test.passIf topic is null, "updateMeta invalid id" 
    updateValidContentTest()


updateValidContentTest = ->
  data.updateContent {id: 1}, "new content 1", (err, data) ->
    test.passIf err is null, "updateContent valid id"

  # Notice that we don't test updateContent with an 
  # invalid id because that will just create an orphan 
  # content file because by design topicData does NOT 
  # perform any kind of validation.


# Kick off the first test
addTopicTest()




