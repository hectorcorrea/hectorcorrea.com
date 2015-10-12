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

var renderStaticPage = function(res, page) {
  logger.info(page);
  if(process.env.NODE_ENV == 'production') {
    res.app.settings.setCache(res, 5);
  }
  res.render(page)
}

// Blog Routes
app.get('/blog', authenticate, blogRoutes.viewAll)
app.post('/blog/:url/:key/edit', authenticate, blogRoutes.edit);
app.post('/blog/:url/:key/save', authenticate, blogRoutes.save);
app.get('/blog/:url/:key', authenticate, blogRoutes.viewOne);
app.get('/blog/rss', blogRoutes.rss);
app.get('/blog/:url', legacyRoutes.blogOne);
app.post('/blog/new', authenticate, blogRoutes.newBlog)
// app.post('/api/blog/:url/:key/draft', authenticate, blogRoutes.draft);
// app.post('/api/blog/:url/:key/post', authenticate, blogRoutes.post);

// Login/logout
app.get('/login', userRoutes.login)
app.post('/login', userRoutes.loginPost)
app.get('/logout', authenticate, userRoutes.logout)
// app.post('/api/user/changePassword', authenticate, userRoutes.changePassword);

// // Blog routes (for Angular.js client)
// app.get('/api/blog/all', authenticate, blogRoutes.all);
// app.get('/api/blog/:url/:key', authenticate, blogRoutes.one);
// app.get('/api/blog/:url/:key/edit', authenticate, blogRoutes.one);
// app.post('/api/blog/:url/:key/draft', authenticate, blogRoutes.draft);
// app.post('/api/blog/:url/:key/post', authenticate, blogRoutes.post);
// app.post('/api/blog/:url/:key', authenticate, blogRoutes.save);
// app.post('/api/blog/new', authenticate, blogRoutes.newOne);

// // Login and authentication (for Angular.js client)
// app.post('/api/user/changePassword', authenticate, userRoutes.changePassword);
// app.post('/api/login', userRoutes.login);
// app.post('/api/logout', authenticate, userRoutes.logout);

// Static pages
app.get('/credits', function(req, res) {
  renderStaticPage(res, 'credits')
});

app.get('/about', function(req, res) {
  renderStaticPage(res, 'about')
});

// Our humble home page (HTML)
app.get('/', function(req, res) {
  logger.info('home page for ' + req.url);
  if(process.env.NODE_ENV == 'production') {
    res.app.settings.setCache(res, 5);
  }
  renderStaticPage(res, 'home')
});

// All others get a horrible 404
app.get('*', function(req, res) {
  logger.error('Not found: ' + req.url);
  if(process.env.NODE_ENV == 'production') {
    res.app.settings.setCache(res, 5);
  }
  res.status(404).render('notFound');
});

// Fire up the web server!
var server = http.createServer(app);
var port = app.get('port');
server.listen(port, function() {
  var address = 'http://localhost:' + port;
  logger.info('Express listening at: ' + address);
});
