var path = require('path');
var http = require('http');
var logger = require('log-hanging-fruit').defaultLogger;
var legacyRoutes = require('./routes/legacyRoutes');
var siteRoutes = require('./routes/siteRoutes');
var blogRoutes = require('./routes/blogRoutes');
var userRoutes = require('./routes/userRoutes');

// Set the path for the log files 
var options = {filePath: path.join(__dirname, 'logs') };
logger.setup(options);

// Configure Express settings
var app = require('./configure').app;

// Legacy Routes
app.get('/about', legacyRoutes.about);
app.get('/blog', legacyRoutes.blogAll)
app.get('/blog/rss', legacyRoutes.rss);
app.get('/blog/:url', legacyRoutes.blogOne);

// New routes (for Angular.js client)
app.get('/api/blog/all', blogRoutes.all);
app.get('/api/blog/:url/:key', blogRoutes.one);
app.post('/api/blog/:url/:key/draft', blogRoutes.draft);
app.post('/api/blog/:url/:key/post', blogRoutes.post);
app.post('/api/blog/:url/:key', blogRoutes.save);
app.post('/api/blog/new', blogRoutes.newOne);

// Login and authentication
app.post('/login/initialize', userRoutes.initialize);
app.post('/login/changePassword', userRoutes.changePassword);
app.post('/login', userRoutes.login);

// Our humble home page (HTML)
app.get('/', function(req, res) {
  logger.info('home page');
  res.render('index')
});

// All others get a horrible 404
app.get('*', function(req, res) {
  logger.error('Not found: ' + req.url);
  res.status(404).render('index.ejs', { error: 'Page not found' });
});

// Fire up the web server! 
var server = http.createServer(app);
var port = app.get('port');
server.listen(port, function() {
  var address = 'http://localhost:' + port;
  logger.info('Express listening at: ' + address);
});
