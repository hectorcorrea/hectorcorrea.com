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
app.use(express.static(path.join(__dirname, 'public'))); 
app.use(express.logger('dev'));
app.use(app.router);

// Global error handler
app.use( function(err, req, res, next) {
  logger.error("Global error handler. Error: " + err);
  res.status(500);
  if(req.xhr) {
    res.send({error: err + ""});
  }
  else {
    res.render('index', {error: err});
  }
});


app.use(express.errorHandler({ dumpExceptions: true, showStack: true }));

// Values used to set expiration cookies throughout the app.
app.set('oneMinInSec', 60);
app.set('oneMinInMs', 60 * 1000);
// app.set('fiveMinInSec', 5 * 60);
// app.set('fiveMinInMs', 5 * 60 * 1000)
// app.set('oneHrInSec', 60 * 60);
// app.set('oneHrInMs', 60 * 60 * 1000);

// Helper function to set cache expire values in response object
//http://stackoverflow.com/a/16750445/446681
app.set('setCache', function(res, minutes) {

  var seconds = res.app.settings.oneMinInSec * minutes;
  var ms = res.app.settings.oneMinInMs * minutes;

  if (!res.getHeader('Cache-Control') || !res.getHeader('Expires')) {
    console.log('Cache set to ' + minutes + ' minutes');
    res.setHeader("Cache-Control", "public, max-age=" + seconds); 
    res.setHeader("Expires", new Date(Date.now() + ms).toUTCString());
  }
  else {
    console.log('Cache cannot be overwritten');
  }

});

var devSettings = function() {
  
  var file = __dirname + "/settings.dev.json";
  logger.info('Loading settings from ' + file);

  var settings = settingsUtil.loadSync(file);
  dbSetup.init(settings.dbUrl);
  app.set("config", settings);

};
app.configure('development', devSettings); 


var prodSettings = function() {

  var file = __dirname + "/settings.prod.json";
  logger.info('Loading settings from ' + file);

  var settings = settingsUtil.loadSync(file);
  if(process.env.DB_URL) {
    settings.dbUrl = process.env.DB_URL;
    dbSetup.init(settings.dbUrl);
  }
  else {
    logger.error("This is not good. No DB_URL environment variable was found.");
  }
  app.set("config", settings);

};
app.configure('production', prodSettings);


exports.app = app;
