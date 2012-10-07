class TestUtil

  constructor: (@namespace = "", @verbose = false, @prefix = "t>") ->
    if namespace isnt ""
      @namespace = namespace + "."

  _logPass: (message) =>
    if @verbose is true
      console.log "#{@prefix} OK #{message}"
    else
      # nothing to log
      

  _logFail: (message) =>
    console.log "#{@prefix} ** Failed ** #{message}"


  passIf: (condition, testName) => 
    if condition is true
      @_logPass "#{@namespace}#{testName}"
    else
      @_logFail "#{@namespace}#{testName}"


  failIf: (condition, testName) =>
    if condition is true
      @_logFail "#{@namespace}#{testName}"
    else
      @_logPass "#{@namespace}#{testName}"


  fail: (testName) =>
    @_logFail "#{@namespace}#{testName}"


  pass: (testName = "") => 
    @_logPass "#{@namespace}#{testName}"
  

exports.TestUtil = TestUtil
