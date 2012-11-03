fs = require 'fs'
path = require 'path'
dateUtil = require '../util/dateUtil'
email = require('emailjs') 
{Logger} = require '../util/logger'
randomUtil = require '../util/randomUtil'

class AuthModel

  constructor: (settings) ->

    # Copy the settings into instance variables
    @dataFile = path.join settings.dataPath, 'auth.json'

    @adminUserEmail = settings.adminUserEmail
    @rootUrl = settings.rootUrl
    @loginKeyValidHours = settings.loginKeyValidHours

    @emailFrom = settings.email.from

    @smtpCredentials = {
      user: settings.email.user
      password: settings.email.password
      host: settings.email.host
      ssl: settings.email.ssl
    }


  _emailLoginKey: (key) => 

    emailText = "Click the link below to authenticate to the site\r\n
      \r\n
      #{@rootUrl}/login/#{key}\r\n
      \r\n
      This link is only valid for #{@loginKeyValidHours} hours."

    message = {
      from: @emailFrom
      to: "admin user <#{@adminUserEmail}>"
      subject: "Login Information"
      text: emailText
    }

    if @smtpCredentials.password is "bogus"

      # Output this to the console, but don't use
      # the logger so that it's not recorded.
      console.log '----------------------------------'
      console.log 'Fake e-mail text below'
      console.log emailText

    else

      mailServer  = email.server.connect(@smtpCredentials)
      mailServer.send message, (err, msg) -> 
        if err
          Logger.error "authModel: Error sending e-mail", err
        else
          Logger.info "authModel: e-mail sent successfully"
      Logger.info 'authModel: sending e-mail...'


  _saveAuthData: (data) =>
    text = JSON.stringify data, null, '\t'
    fs.writeFileSync @dataFile, text, 'utf8'
    return data


  _initializeDataFile: =>
    data = {
      user: @adminUserEmail
    }
    @_saveAuthData data
    return data


  _readAuthData: =>
    text = fs.readFileSync @dataFile, 'utf8'
    data = JSON.parse text


  generateLoginKey: =>
    key = randomUtil.randomString(10)
    today = new Date()
    expire = dateUtil.addHours today, @loginKeyValidHours
    data = {
      user: @adminUserEmail
      loginKey: key
      loginKeyExpire: expire
    }
    @_saveAuthData data
    @_emailLoginKey key
    return data


  generateAuthToken: =>
    token = randomUtil.randomString(10)
    data = {
      user: @adminUserEmail
      authToken: token
    }
    @_saveAuthData data
    return data


  clearAuthData: =>
    return @_initializeDataFile()


  loadAuthData: =>
    if fs.existsSync(@dataFile)
      @_readAuthData()
    else
      @_initializeDataFile()
  

  isAuthenticated: (req) =>
    authenticated = false
    cookie = req.cookies.session

    authTokenInCookie = cookie?.authToken? 
    if authTokenInCookie
      authData = @loadAuthData()
      if authData.authToken? 
        authenticated = true if authData.authToken is cookie.authToken 

    #console.log "authenticated? #{authenticated}"
    return authenticated


exports.AuthModel = AuthModel
