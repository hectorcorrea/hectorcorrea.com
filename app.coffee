# Module dependencies.
express = require 'express'
path = require 'path'
ejs = require 'ejs'
http = require 'http'
os = require 'os'
fs = require 'fs'

{Logger} = require './util/logger'
settings = require './util/settings'
siteRoutes = require './routes/siteRoutes'
blogRoutes = require './routes/blogRoutes'
logRoutes = require './routes/logRoutes'
authRoutes = require './routes/authRoutes'


# Configuration
app = express()
app.configure -> 
  app.set 'port', process.env.PORT || 3000
  app.set 'views', path.join(__dirname, 'views')
  app.set 'datapath', path.join(__dirname, 'data')

  # View engine options
  ejs.open = '{{'
  ejs.close = '}}'
  app.set 'view engine', 'ejs'

  app.use express.favicon()
  app.use express.logger('dev')

  app.use express.bodyParser()
  app.use express.methodOverride()

  app.use express.cookieParser('your secret here')
  app.use express.session()

  # Static handler must come before app.router!
  app.use express.static path.join(__dirname, 'public')   

  app.use app.router


  # Global error handler
  app.use (err, req, res, next) ->
    Logger.error "Global error handler. Error #{err}"
    res.status 500
    res.render '500', {message: err}


# Development settings
app.configure 'development', -> 
  app.use express.errorHandler({ dumpExceptions: true, showStack: true })
  app.set "dataOptions", settings.load 'settings.dev.json', __dirname


# Production settings
app.configure 'production', ->
  app.use express.errorHandler()
  app.set "dataOptions", settings.load 'settings.prod.json', __dirname


# Routes
app.get '/', siteRoutes.home
app.get '/about', siteRoutes.about
app.get '/credits', siteRoutes.credits

app.get '/blog/new', blogRoutes.editNew
app.post '/blog/new', blogRoutes.saveNew

app.get '/blog/edit/:topicUrl', blogRoutes.edit
app.post '/blog/save/:id', blogRoutes.save

app.get '/blog/list', blogRoutes.viewAll

app.get '/blog/rss', blogRoutes.rssList

# Switch to viewRecent when blog list gets too long
app.get '/blog', blogRoutes.viewAll

app.get '/blog/:topicUrl', blogRoutes.viewOne

app.get '/logs/current', logRoutes.viewCurrent
app.get '/logs/:logDate', logRoutes.viewSpecific
app.get '/logs/', logRoutes.viewCurrent
 
app.get '/login/:key', authRoutes.loginConfirm
app.get '/login', authRoutes.loginGet
app.post '/login', authRoutes.loginPost

app.get '/logout', authRoutes.logout

# Test route to force an error
app.get '/blowup', siteRoutes.blowUp 

app.get '*', siteRoutes.notFound


# Fire it up!
server = http.createServer(app)
port = app.get('port')
server.listen port, ->
  address = "http://localhost:#{port}"
  machine = os.hostname()
  environment = app.settings.env
  Logger.info "Express server listening on #{address} in #{environment} mode (#{machine})"
