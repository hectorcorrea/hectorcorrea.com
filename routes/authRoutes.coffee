{Logger} = require '../util/logger'
authModel = require '../models/authModel'

_setAuthCookie = (res, authToken) ->
  userInfo = {authToken: authToken}
  
  oneHr = 60 * 60 # in seconds
  oneDay = oneHr * 24
  oneMonth = oneDay * 30
  res.cookie 'session', userInfo, {maxAge: oneMonth} 


_clearAuthCookie = (res) ->
  res.clearCookie 'session'


loginGet = (req, res) ->
  Logger.info 'authRoutes:loginGet'

  dataPath = res.app.settings.dataOptions.dataPath
  authData = authModel.readAuthData dataPath
  if authData.user? is false
    Logger.info "***************************************"
    Logger.info "No user name has been configured in the"
    Logger.info "auth.json data file. All login attempts"
    Logger.info "will fail."
    Logger.info "***************************************"
  res.render 'login'


loginPost = (req, res) ->
  Logger.info 'authRoutes:loginPost'

  email = req.body?.email ? ''
  dataPath = res.app.settings.dataOptions.dataPath
  authData = authModel.readAuthData dataPath

  if email is authData.user
    loginKey = authModel.getRandomKey()
    #TODO: e-mail loginKey and remove from Logging!!!
    console.log "#{loginKey}"
    authModel.saveLoginKey dataPath, loginKey
    Logger.info "e-mail with loginKey has been sent"
    res.render 'loginPost'
  else
    Logger.warn "Invalid login attempt. E-mail received [#{email}]"
    res.render 'login', {errorMsg: 'The e-mail indicated is not valid.'}


loginConfirm = (req, res) ->
  Logger.info 'authRoutes:loginConfirm'
  dataPath = res.app.settings.dataOptions.dataPath
  loginKey = req.params.key
  
  if loginKey
    authData = authModel.readAuthData dataPath
    if authData.loginKey? and authData.loginKey is loginKey
      authToken = authModel.getRandomKey()
      _setAuthCookie res, authToken
      authModel.saveAuthToken dataPath, authToken 
      Logger.info 'Woo-hoo! your are in.'
      res.redirect '/'
    else
      Logger.warn "loginKey received [#{loginKey}] is not valid"
      res.redirect '/'
  else
    Logger.warn 'No loginKey received'
    res.redirect '/'


logout = (req, res) ->
  Logger.info 'authRoutes:logout'
  dataPath = res.app.settings.dataOptions.dataPath
  authModel.clearAuthData dataPath
  _clearAuthCookie res
  res.redirect '/'


module.exports = {
  loginGet: loginGet
  loginPost: loginPost
  loginConfirm: loginConfirm
  logout: logout
}

