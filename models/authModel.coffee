fs = require 'fs'
path = require 'path'
dateUtil = require '../util/dateUtil'
email = require('emailjs') 
{Logger} = require '../util/logger'


# TODO: Make these values configurable
emailFrom = "hector j corrrea <hectorjcorrea@gmail.com>" 
adminUserEmail = "hector@hectorcorrea.com"
smtpCredentials = {
  user: "bogus"
  password: "bogus" 
  host: "smtp.gmail.com" 
  ssl: true
}
loginKeyValidHours = 2


_emailLoginKey = (key) -> 

  emailText = "Click the link below to authenticate to the site\r\n
    \r\n
    http://hectorcorrea.com/login/#{key}\r\n
    \r\n
    This link is only valid for #{loginKeyValidHours} hours."

  message = {
    from: emailFrom
    to: "admin user <#{adminUserEmail}>"
    subject: "Login Information"
    text: emailText
  }

  mailServer  = email.server.connect(smtpCredentials)
  mailServer.send message, (err, msg) -> 
    if err
      Logger.error "authModel: Error sending e-mail", err
    else
      Logger.info "authModel: e-mail sent successfully"

  Logger.info 'authModel: sending e-mail...'


_randomNumber = (max, min) -> 
  # http://stackoverflow.com/a/1527834/446681
  r = Math.floor(Math.random() * (max - min + 1)) + min


_randomUpperChar = ->
  n = _randomNumber(65, 90); # A..Z
  String.fromCharCode(n)


_saveAuthData = (filePath, data) ->
  text = JSON.stringify data, null, '\t'
  fileName = path.join filePath, 'auth.json'
  fs.writeFileSync fileName, text, 'utf8'
  return


getRandomKey = (length = 10) ->
  key = ""
  for i in [1..length]
    key += _randomUpperChar()
  key


saveLoginKey = (filePath, key) ->
  today = new Date()
  expire = dateUtil.addHours today, loginKeyValidHours
  data = {
    user: adminUserEmail
    loginKey: key
    loginKeyExpire: expire
  }
  _saveAuthData filePath, data
  _emailLoginKey key


saveAuthToken = (filePath, token) ->
  data = {
    user: adminUserEmail
    authToken: token
  }
  _saveAuthData filePath, data


clearAuthData = (filePath) ->
  data = {
    user: adminUserEmail
  }
  _saveAuthData filePath, data


readAuthData = (filePath) ->
  fileName = path.join filePath, 'auth.json'
  if fs.existsSync(fileName)
    text = fs.readFileSync fileName, 'utf8'
    data = JSON.parse text
  else
    # TODO: initialize with admin user name.
    # data = { "user": "somebody@gmail.com" }
    data = {}
  return data


isAuthenticated = (req, filePath) ->
  authenticated = false
  cookie = req.cookies.session
  if cookie?.authToken? 
    console.log "cookie value found"
    authData = readAuthData filePath
    if authData.authToken? 
      console.log "data value found"
      authenticated = true if authData.authToken is cookie.authToken 
  console.log "authenticated? #{authenticated}"
  return authenticated


module.exports = {
  getRandomKey: getRandomKey
  saveLoginKey: saveLoginKey
  saveAuthToken: saveAuthToken
  clearAuthData: clearAuthData
  readAuthData: readAuthData
  isAuthenticated: isAuthenticated
}
