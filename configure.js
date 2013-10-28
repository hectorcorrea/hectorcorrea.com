var express = require('express');
var path = require('path');
var ejs = require('ejs');
var logger = require('log-hanging-fruit').defaultLogger;
var settingsUtil = require('./settings');
var dbSetup = require('./models/dbSetup');

var initialize = function() {

  app.set('port', process.env.PORT || 3000);
  app.set('views', path.join(__dirname, 'views'));

  app.set('view engine', 'ejs');

  app.use(express.bodyParser());
  app.use(express.methodOverride());

  app.use(express.cookieParser('your secret here'));
  app.use(express.session());

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

};


var devSettings = function() {
  
  var file = __dirname + "/settings.dev.json";
  logger.info('Loading settings from ' + file);

  var settings = settingsUtil.loadSync(file);
  dbSetup.init(settings.dbUrl);
  app.set("config", settings);

};


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


var app = express();
app.configure(initialize); 
app.configure('development', devSettings); 
app.configure('production', prodSettings);

exports.app = app;
