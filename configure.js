var express = require('express');
var path = require('path');
var ejs = require('ejs');
var logger = require('log-hanging-fruit').defaultLogger;
var settingsUtil = require('./settings');
var dbSetup = require('./models/dbSetup');
var app = express();

app.set('port', process.env.PORT || 3000);
app.set('views', path.join(__dirname, 'views'));

app.set('view engine', 'ejs');

app.use(express.bodyParser());
app.use(express.methodOverride());

app.use(express.cookieParser('drink more coffee'));
app.use(express.session({cookie: { httpOnly: false }}));

// static must appear before app.router!
var viewsPath = path.join(__dirname, 'public');
if(process.env.NODE_ENV == 'production') {
  var oneDay = 86400000; 
  app.use(express.static(viewsPath, {maxAge: oneDay})); 
}
else {
  app.use(express.static(viewsPath)); 
}

app.use(express.logger('dev'));
app.use(app.router);

// Global error handler
app.use( function(err, req, res, next) {
  logger.error('Global error handler. Error: ' + err);
  res.status(500);
  if(req.xhr) {
    res.send({error: err + ''});
  }
  else {
    res.render('index', {error: err});
  }
});


app.use(express.errorHandler({ dumpExceptions: true, showStack: true }));


// Helper function to set cache expire values in response object
// http://stackoverflow.com/a/16750445/446681
app.set('setCache', function(res, minutes) {

  var seconds = minutes * 60;
  var ms = seconds * 1000;

  if (!res.getHeader('Cache-Control') || !res.getHeader('Expires')) {
    res.setHeader('Cache-Control', 'public, max-age=' + seconds); 
    res.setHeader('Expires', new Date(Date.now() + ms).toUTCString());
  }

});


var devSettings = function() {
  
  var file = __dirname + '/settings.dev.json';
  logger.info('Loading settings from ' + file);

  var settings = settingsUtil.loadSync(file);
  var defaultUser = {
    user: settings.defaultUser, 
    password: settings.defaultPassword, 
    salt: process.env.BLOG_SALT
  };
  dbSetup.init(settings.dbUrl, defaultUser);
  app.set('config', settings);

};
app.configure('development', devSettings); 


var prodSettings = function() {

  var file = __dirname + '/settings.prod.json';
  logger.info('Loading settings from ' + file);

  var settings = settingsUtil.loadSync(file);
  var defaultUser = {
    user: process.env.BLOG_USER, 
    password: process.env.BLOG_PWD, 
    salt: process.env.BLOG_SALT
  };

  console.dir(defaultUser);
  if(process.env.DB_URL) {
    settings.dbUrl = process.env.DB_URL;
    dbSetup.init(settings.dbUrl, defaultUser);
  }
  else {
    logger.error('This is not good. No DB_URL environment variable was found.');
  }
  app.set('config', settings);

};
app.configure('production', prodSettings);


exports.app = app;
