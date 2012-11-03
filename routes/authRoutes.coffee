{Logger} = require '../util/logger'
{AuthModel} = require '../models/authModel'

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

  authModel = new AuthModel(res.app.settings.dataOptions)
  authData = authModel.loadAuthData()

  if authData.user? is false
    Logger.info "***************************************"
    Logger.info "No user name has been configured in the"
    Logger.info "auth.json data file. All login attempts"
    Logger.info "will fail. Review your settings file.  "
    Logger.info "***************************************"
  res.render 'login'


loginPost = (req, res) ->
  Logger.info 'authRoutes:loginPost'

  email = req.body?.email ? ''
  authModel = new AuthModel(res.app.settings.dataOptions)
  authData = authModel.loadAuthData()

  if email is authData.user
    data = authModel.generateLoginKey()
    Logger.info "E-mail with loginKey has been sent"
    console.dir data # don't log this value
    res.render 'loginPost'
  else
    Logger.warn "Invalid login attempt. E-mail received [#{email}]"
    res.render 'login', {errorMsg: 'The e-mail indicated is not valid.'}


loginConfirm = (req, res) ->
  Logger.info 'authRoutes:loginConfirm'

  loginKey = req.params.key
  if loginKey

    authModel = new AuthModel(res.app.settings.dataOptions)
    authData = authModel.loadAuthData()
    if authData.loginKey? and authData.loginKey is loginKey

      authData = authModel.generateAuthToken()
      _setAuthCookie res, authData.authToken
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
  authModel = new AuthModel(res.app.settings.dataOptions)
  authModel.clearAuthData()
  _clearAuthCookie res
  res.redirect '/'


module.exports = {
  loginGet: loginGet
  loginPost: loginPost
  loginConfirm: loginConfirm
  logout: logout
}

