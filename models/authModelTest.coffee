{AuthModel} = require './authModel'
{TestUtil}  = require '../util/testUtil'
settings = require '../util/settings'


# Notice that these tests run synchronously

test = new TestUtil("autModelTest.EmptyRequest", true)
testSettings = settings.load 'settings.test.json', __dirname
model = new AuthModel(testSettings)
model.clearAuthData()

req = {
  cookies: {}
}

test.passIf model.isAuthenticated(req) is false, "empty req/empty data"


data = model.generateLoginKey()
test.passIf data.loginKey?, 'generate loginKey'


data = model.generateAuthToken()
test.passIf data.authToken?, 'generate authToken'


req = {
  cookies: {
    session: {
      authToken: data.authToken
    }
  }
}
test.passIf model.isAuthenticated(req), "authenticated request"

