fs = require 'fs'
{Logger} = require '../util/logger'


_isValidLogDate = (text) -> 

  # Date must be in format YYYY-MM-DD and only years 20YY are valid.
  regEx = /20\d\d-\d\d-\d\d/
  match = text.match(regEx)
  return false if match is null

  # Date must match input string so that we don't have any extra characters
  # on the string. For example 2012-09-23XX is rejected.
  date = match[0]
  return false if date isnt text

  date = new Date(text)
  if date instanceof Date && isFinite(date)
    return true # yay!

  return false


#TODO: This code should use streams rather than reading the
# entire log file to memory. 
viewCurrent = (req, res) ->  

  logFile = Logger.currentLogFile()
  Logger.info 'logRoutes.viewCurrent'
  fs.readFile logFile, (err, text) ->
    if err
      Logger.error "logRoutes.viewSpecific: Error reading log file #{logFile}\r\n#{err}"
      res.render '500', {status: 500, message: "Could not read log file: #{logFile}"}
    else
      text = text.toString().replace(/\r\n/g, '<br/>')
      res.send text


viewSpecific = (req, res) ->  

  logDate = req.params.logDate
  if _isValidLogDate logDate
    Logger.info "logRoutes.viewSpecific #{logDate}"
    logDate = logDate.replace(/-/g, '_')
    logFile = "#{res.app.settings.dataOptions.logPath}/#{logDate}.txt"

    fs.readFile logFile, (err, text) ->
      if err
        Logger.error "logRoutes.viewSpecific: Error reading log file #{logFile}\r\n#{err}"
        res.render '500', {status: 500, message: "Could not read log file: #{logFile}"}
      else
        text = text.toString().replace(/\r\n/g, '<br/>')
        res.send text
  else
    Logger.error "logRoutes.viewSpecific: Invalid log date received #{logDate}"
    res.redirect '/'


list = (req, res) ->
  # todo: implement code to view a list of log files

module.exports = {
  viewCurrent: viewCurrent
  viewSpecific: viewSpecific
}
