fs = require 'fs'

class Logger


  @_logLevel = 'INFO'
  @_logPath = './logs/'   # set up null to prevent log to text file
  @_logFile = ''          # calculated in _getTimestamp()


  @_pad: (number, zeroes) =>
    return ('000000' + number).slice(-zeroes)


  @_getTimestamp: =>
    now = new Date()

    day = now.getDate()
    month = now.getMonth() + 1
    date = now.getFullYear() + '-' + @_pad(month, 2) + '-' + @_pad(day, 2) 

    if @_logPath isnt null
      # Make sure the log file matches the current date
      fileName = now.getFullYear() + '_' + @_pad(month, 2) + '_' + @_pad(day, 2) + '.txt'
      if @_logFile isnt fileName
        @_logFile = fileName
    
    hours = now.getHours()
    minutes = now.getMinutes()
    seconds = now.getMinutes()
    milliseconds = now.getMilliseconds()
    time = @_pad(hours, 2) + ':' + @_pad(minutes, 2) + ':' + @_pad(seconds, 2) + '.' + @_pad(milliseconds, 3)

    timestamp = date + ' ' + time
    return timestamp


  @_doLog: (level, text) =>

    textToLog = "#{@_getTimestamp()} #{level}: #{text}"

    # Log to console
    console.log textToLog
    
    if @_logPath isnt null
      # Log to disk
      fullLogFileName = @_logPath + @_logFile
      fs.appendFile fullLogFileName, textToLog + '\r\n', (err) ->
        if err
          console.log "ERROR writting to log file #{fullLogFileName}. Error: [#{err}]"


  @setLevel: (level) =>
    if level.toUpperCase() in ['INFO', 'WARN', 'ERROR', 'NONE']
      @_logLevel = level.toUpperCase()


  @info: (text) =>
    if @_logLevel is 'INFO'
      @_doLog 'INFO', text


  @warn: (text) =>
    if @_logLevel in ['INFO', 'WARN']
      @_doLog 'WARN', text


  @error: (text) =>
    if @_logLevel in ['INFO', 'WARN', 'ERROR']
      @_doLog 'ERROR', text


  @error: (text, exception = null) =>
    if @_logLevel in ['INFO', 'WARN', 'ERROR']
      text = text + "\r\n#{exception}" if exception?
      @_doLog 'ERROR', "#{text}"


  @currentLogFile: => 
    return null if @_logPath is null
    # Make sure @_logFile has been set
    @_getTimestamp() 
    return @_logPath + @_logFile


  @setPath: (path) =>
    @_logPath = path


exports.Logger = Logger

