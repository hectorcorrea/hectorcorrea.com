var path = require('path');
var http = require('http');
var logger = require('log-hanging-fruit').defaultLogger;
var legacyRoutes = require('./routes/legacyRoutes');
var siteRoutes = require('./routes/siteRoutes');
var blogRoutes = require('./routes/blogRoutes');
var userRoutes = require('./routes/userRoutes');

var testCounter = 0;

// Set the path for the log files 
var options = {filePath: path.join(__dirname, 'logs') };
logger.setup(options);

// Configure Express settings
var app = require('./configure').app;

// Authentication middleware
var authenticate = function(req, res, next) {

  req.isAuth = false;
  if(req.cookies && req.cookies.authToken) {
    // Make sure the session is valid.
    userRoutes.validateSession(req, res, next);
  }
  else {
    // Not authenticated, nothing else to do.
    next();
  }

}


// Legacy Routes (redirect to new URLs)
app.get('/about', legacyRoutes.about);
app.get('/blog', legacyRoutes.blogAll)
app.get('/blog/rss', blogRoutes.rss);
app.get('/blog/:url', legacyRoutes.blogOne);

// Blog routes (for Angular.js client)
app.get('/api/blog/all', authenticate, blogRoutes.all);
app.get('/api/blog/:url/:key', authenticate, blogRoutes.one);
app.get('/api/blog/:url/:key/edit', authenticate, blogRoutes.one);
app.post('/api/blog/:url/:key/draft', authenticate, blogRoutes.draft);
app.post('/api/blog/:url/:key/post', authenticate, blogRoutes.post);
app.post('/api/blog/:url/:key', authenticate, blogRoutes.save);
app.post('/api/blog/new', authenticate, blogRoutes.newOne);

// Login and authentication (for Angular.js client)
app.post('/api/login/initialize', userRoutes.initialize);
app.post('/api/user/changePassword', authenticate, userRoutes.changePassword);
app.post('/api/login', userRoutes.login);
app.post('/api/logout', authenticate, userRoutes.logout);

// Our humble home page (HTML)
app.get('/', function(req, res) {
  logger.info('home page for ' + req.url);
  if(process.env.NODE_ENV == 'production') {
    res.app.settings.setCache(res, 5);
  }
  res.render('index')
});

// All others get a horrible 404
app.get('*', function(req, res) {
  logger.error('Not found: ' + req.url);
  if(process.env.NODE_ENV == 'production') {
    res.app.settings.setCache(res, 5);
  }
  res.status(404).render('index.ejs', { error: 'Page not found' });
});

// Fire up the web server! 
var server = http.createServer(app);
var port = app.get('port');
server.listen(port, function() {
  var address = 'http://localhost:' + port;
  logger.info('Express listening at: ' + address);
});
