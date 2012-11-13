# Tests for topicModel class.
#
# Notice that these tests must be run in order
# in order so that the first ones add data that 
# the next ones can use.

fs = require 'fs'
{TopicModel}  = require './topicModel'
{TestUtil}  = require '../util/testUtil'
settings = require '../util/settings'

verbose = true
dataOptions = settings.load 'settings.test.json', __dirname
dataOptions.showDrafts = true


# Delete the current data file (if any)
dataFile = dataOptions.dataPath + '/blogs.json'
if fs.existsSync dataFile
  fs.unlinkSync dataFile


model = new TopicModel(dataOptions)


testGetUrlFromTitle = ->
  test = new TestUtil("topicModelTest.testGetUrlFromTitle", verbose)
  test.passIf model._getUrlFromTitle("hello") is "hello", "basic test"
  test.passIf model._getUrlFromTitle("hello-World") is "hello-world", "lowercase test"
  test.passIf model._getUrlFromTitle("hello-World.aspx") is "hello-world-aspx", "dots test"
  test.passIf model._getUrlFromTitle("hello-c#-World.aspx") is "hello-csharp-world-aspx", "c# test"
  test.passIf model._getUrlFromTitle("this is #4") is "this-is-4", "pound (#) test"
  testValidateGoodTopic()


testValidateGoodTopic = ->
  test = new TestUtil("topicModelTest.testValidateGoodTopic", verbose)
  goodTopic = {
    meta: {
      id: 1
      title: "hello world title"
      summary: "hello world summary"
    }
    content: "hello world content"
  }

  model._validateTopic goodTopic, (err, validationErrors) ->
    test.passIf err is null and validationErrors is null, "good topic" 
    testValidateEmptyTopic()


testValidateEmptyTopic = ->
  test = new TestUtil("topicModelTest.testValidateEmptyTopic", verbose)
  emptyTopic = {}
  model._validateTopic emptyTopic, (err, validationErrors) ->
    test.passIf validationErrors isnt null, "empty topic" 
    testValidateEmptyTitle()


testValidateEmptyTitle = ->
  test = new TestUtil("topicModelTest.testValidate", verbose)
  emptyTitleTopic = { meta: { id: 1, summary: "s"}, content: "c" }
  model._validateTopic emptyTitleTopic, (err, validationErrors) ->
    test.passIf validationErrors.emptyTitle is true, "empty title"
    testSaveNewGoodTopic()


testSaveNewGoodTopic = ->
  test = new TestUtil("topicModelTest.testSaveNewGoodTopic", verbose)
  newTopic = {
    meta: {
      title: "new test topic",
      summary: "new summary for test topic"
    }
    content: "new content for test topic"
  }

  model.saveNew newTopic, (err, data) ->
    if err
      test.fail "unexpected error #{err}" 
    else 
      model.getOne data.meta.id, (err, data) ->
        if err isnt null
          test.fail "error retrieving new record id: #{data.meta.id} #{err}"
        else
          test.pass ""
        testIsDuplicateTitleNew()  


testIsDuplicateTitleNew = ->
  test = new TestUtil("topicModelTest.testIsDuplicateTitleNew", verbose)

  newTopic = {
    meta: {
      title: "new test topic",
      summary: "new summary for test topic 2"
    }
    content: "new content for test topic 2"
  }

  model._isDuplicateTitle newTopic, (err, isDuplicate) ->
    test.passIf isDuplicate, ""
    testIsDuplicateTitleExisting()  


testIsDuplicateTitleExisting = ->
  test = new TestUtil("topicModelTest.testIsDuplicateTitleExisting", verbose)

  model.getOne 1, (err, topic) ->
    if err
      test.fail err
    else
      model._isDuplicateTitle topic, (err, isDuplicate) ->
        test.failIf isDuplicate, ""
        testSaveNewBadTopic()  


testSaveNewBadTopic = ->
  test = new TestUtil("topicModelTest.testSaveNewBadTopic", verbose)

  badNewTopic = { meta: { title: "" }, content: "blah blah blah" }
  model.saveNew badNewTopic, (err, topic) ->
    if err
      test.fail "#{err}" 
    else 
      test.passIf topic.errors isnt null, "" 
    testGetAll()


testGetAll = ->
  test = new TestUtil("topicModelTest", verbose)
  model.getAll (err, topics) ->
    test.passIf err is null and topics.length > 0, "getAll"
    testGetRecent()


testGetRecent = ->
  test = new TestUtil("topicModelTest", verbose)
  model.getRecent (err, topics) ->
    test.passIf err is null and topics.length > 0, "getRecent"
    testOneValid()


testOneValid = ->
  test = new TestUtil("topicModelTest", verbose)
  model.getOne 1, (err, data) ->
    test.passIf data.meta.id is 1, "getOne valid id"
    testOneInvalid()


testOneInvalid = ->
  test = new TestUtil("topicModelTest", verbose)
  model.getOne 99, (err, data) ->
    test.passIf err isnt null, "getOne invalid id"
    testOneValidUrl()


testOneValidUrl = ->
  test = new TestUtil("topicModelTest", verbose)
  model.getOne 1, (err, topic1) ->
    if err
      test.fail "getOneByUrl valid url"
    else
      model.getOneByUrl topic1.meta.url, (err, topic) ->
        test.passIf topic.meta.url is topic1.meta.url, "getOneByUrl valid id"
        testOneInvalidUrl()


testOneInvalidUrl = ->
  test = new TestUtil("topicModelTest", verbose)
  model.getOneByUrl 'topic-99', (err, topic) ->
    test.passIf err isnt null, "getOneByUrl invalid id"
    testSaveGoodTopic()


testSaveGoodTopic = ->
  test = new TestUtil("topicModelTest", verbose)

  goodTopic = {
    meta: {
      id: 1
      title: "new title 1"
      summary: "new summary 1"
    }
    content: "updated content 1"
  }

  model.save goodTopic, (err, data) ->
    if err
      console.dir err
      test.fail "saveGoodTopic"
    else
      errors = []
      if data.meta.id isnt 1
        errors.push "Invalid id received after save"
      if data.meta.title isnt "new title 1"
        errors.push "Invalid title after save"
      if data.content isnt "updated content 1"
        errors.push "Invalid content after save"
      if data.meta.updatedOn? is false
        errors.push "Updated on not populated"
    
      test.passIf errors.length is 0, "saveGoodTopic"
      if errors.length > 0
        console.dir errors
      testSaveBadTopic()


testSaveBadTopic = ->
  test = new TestUtil("topicModelTest", verbose)
  badData = {meta: {id: 1, title: ""}}
  model.save badData, (err, data) ->
    test.passIf data.errors.emptyTitle, "testSaveBadTopic"
    testSaveNonExistingTopic()


testSaveNonExistingTopic = ->
  test = new TestUtil("topicModelTest", verbose)
  notExistingTopic = {meta: {id: 99}}
  model.save notExistingTopic, (err, data) ->
    test.passIf err isnt null, "testSaveNonExistingTopic"
    testRssList()


testRssList = ->
  test = new TestUtil("topicModelTest", verbose)
  model.getRssList (err, xml) ->
    test.passIf err is null, "testRssList"


# -------------------
# Kick off the tests
testGetUrlFromTitle()









