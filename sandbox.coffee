_getUserInfoFromRequest = (req) ->
  userInfo = {isAuth: false}
  if req.cookies?.session?
    cookie = req.cookies?.session
    if cookie.authToken?
      userInfo.authToken = cookie.authToken
  return userInfo


_isValidAuthToken = (data, authToken) ->
  if authToken? and data.authToken is authToken
    return true
  return false 


_isValidLoginKey = (data, loginKey) ->
  if data.loginKey? and data.loginKey is loginKey
    # TODO: validate that key is within certain timeframe
    return true 
  return false  


_generateLoginKey = ->
  return "abc" 


_generateAuthToken = ->
  return "xyz"


_setUserInfoInResponse = (res, name, authToken) ->
  userInfo = {authToken: authToken}
  
  oneHr = 1000 * 60 * 60
  oneDay = oneHr * 24
  oneMonth = oneDay * 30
  res.cookie "session", userInfo, {maxAge: oneMonth} 


_isAuthenticatedRequest = (req, data) ->
  userInfo = getUserInfoFromRequest(req)
  return isValidAuthToken(data, userInfo.authToken)

