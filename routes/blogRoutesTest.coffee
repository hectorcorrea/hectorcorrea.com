# Tests for blogRoutes class.
#
# Notice that these tests must be run in order
# in order so that the first ones add data that 
# the next ones can use.

fs = require 'fs'
blogRoutes = require './blogRoutes'
{TestUtil}  = require '../util/testUtil'
{Logger} = require '../util/logger'
settings = require '../util/settings'

verbose = true
test = new TestUtil("blogRoutesTest", verbose)
Logger.setLevel 'NONE'

dataOptions = settings.load 'settings.test.json', __dirname
dataOptions.showDrafts = true


getBasicApp = ->
  app = {
    settings: { 
      dataOptions: dataOptions
      isReadOnly: false 
    } 
  }
  return app


# Delete the current data file (if any)
dataFile = dataOptions.dataPath + '/blogs.json'
if fs.existsSync dataFile
  fs.unlinkSync dataFile


getBasicRequest = ->
  req = { 
    app: getBasicApp() 
    cookies: {
      session: {authToken: "unittest"}
    }
  }
  return req


getBasicResponse = ->
  res = {
    app: getBasicApp()
  }
  return res


normalizeTopicTitle = ->
  test.passIf blogRoutes._normalizeTopicTitle('a-TOPIC-name') is 'a-topic-name', "_normalizeTopicTitle - mixed case"
  test.passIf blogRoutes._normalizeTopicTitle('a-TOPIC-name.aspx') is 'a-topic-name', "_normalizeTopicTitle - .aspx"
  saveNew() # fire next test


saveNew = ->
  req = getBasicRequest()
  req.body = {title: "title one", summary: "new summary", content: "new content"}

  res = getBasicResponse()
  res.redirect = (url) ->
    console.log url
    test.passIf url is "/blog/title-one", "saveNew"
    viewOneValid() # fire next test

  blogRoutes.saveNew req, res


viewOneValid = ->
  req = getBasicRequest()
  req.params = { topicUrl: "title-one" }

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogOne", "viewOneValid"
    viewOneInvalid() # fire next test

  blogRoutes.viewOne req, res


viewOneInvalid = ->
  req = getBasicRequest()
  req.params = { topicUrl: "topic-99" }

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "404", "viewOneInvalid"
    viewRecent() # fire next test

  blogRoutes.viewOne req, res


viewRecent = ->
  req = getBasicRequest()
  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogRecent", "viewRecent"
    viewAll() # fire next test

  blogRoutes.viewRecent req, res


viewAll = ->
  req = getBasicRequest()
  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogAll", "viewAll"
    editNoUrl()

  blogRoutes.viewAll req, res


editNoUrl = ->
  req = getBasicRequest()
  req.params = {}
  res = getBasicResponse()
  res.redirect = (redirUrl) ->
    test.passIf redirUrl is "/blog", "editNoUrl"
    editBadUrl()

  blogRoutes.edit req, res


editBadUrl = ->
  req = getBasicRequest()
  req.params = {topicUrl: "topic-99"}

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "404", "editBadUrl"
    editGoodUrl() 

  blogRoutes.edit req, res


editGoodUrl = ->
  req = getBasicRequest()
  req.params = {topicUrl: "title-one"}

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogEdit", "editGoodUrl"
    saveNoId()

  blogRoutes.edit req, res


saveNoId = ->

  req = getBasicRequest()
  req.params = {}

  res = getBasicResponse()
  res.redirect = (redirUrl) ->
    test.passIf redirUrl is "/blog", "saveNoId"
    saveBadId()

  blogRoutes.save req, res


saveBadId = ->

  req = getBasicRequest()
  req.params = {id: "ABC"}
  req.body = {title: "t1", summary: "s1", content: "c1"}

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "500", "saveBadId"
    saveNonExistingId()

  blogRoutes.save req, res


saveNonExistingId = ->

  req = getBasicRequest()
  req.params = {id: 99}
  req.body = {title: "t1", summary: "s1", content: "c1"}

  res = getBasicResponse()
  res.redirect = (url) ->
    test.passIf url is "/blog", "saveNonExistingId"
    saveNoBody()

  blogRoutes.save req, res


saveNoBody = ->

  req = getBasicRequest()
  req.params = {id: 1}

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogEdit" and viewModel.topic.errors.emptyTitle, "saveNoBody"
    saveIncompleteData()

  blogRoutes.save req, res


saveIncompleteData = ->

  req = getBasicRequest()
  req.params = {id: 1}
  req.body = {title: "", summary: "s1", content: ""}

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogEdit" and 
      viewModel.topic.errors.emptyTitle and 
      viewModel.topic.errors.emptyContent, "saveIncompleteData"
      saveCompleteData()

  blogRoutes.save req, res


saveCompleteData = ->

  req = getBasicRequest()
  req.params = {id: 1}
  req.body = {title: "updated title 2", summary: "s1", content: "c2"}

  res = getBasicResponse()
  res.redirect = (url) ->
    test.passIf url is "/blog/updated-title-2", "saveCompleteData"
    editNew()

  blogRoutes.save req, res


editNew = ->
  req = getBasicRequest()

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogEdit", "editNew"
    saveNewWithErrors()

  blogRoutes.editNew req, res


saveNewWithErrors = ->
  req = getBasicRequest()
  req.body = {title: "", summary: "new summary", content: "new content"}

  res = getBasicResponse()
  res.render = (page, viewModel) ->
    test.passIf page is "blogEdit" and 
      viewModel.topic.errors.emptyTitle, "saveNewWithErrors"

  blogRoutes.saveNew req, res


# -------------------
# Kick off the tests
normalizeTopicTitle()

