fs = require('fs')

_readSettings = (fileName) ->
  if fs.existsSync(fileName) is false
    throw "Settings file #{fileName} was not found"

  text = fs.readFileSync fileName, 'utf8'
  return JSON.parse text


load = (fileName, rootDir = "") ->
  settings = _readSettings fileName 
  settings.dataPath = settings.dataPath.replace '#{rootDir}', rootDir 
  settings.logPath = settings.logPath.replace '#{rootDir}', rootDir 
  return settings    


module.exports = {
  load: load
}