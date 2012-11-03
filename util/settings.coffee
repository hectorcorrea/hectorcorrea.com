fs = require('fs')

_readSettings = (fileName) ->
  if fs.existsSync(fileName) is false
    throw "Settings file #{fileName} was not found"

  text = fs.readFileSync fileName, 'utf8'
  return JSON.parse text


load = (fileName) ->
  settings = _readSettings fileName 
  settings.dataPath = settings.dataPath.replace '#{__dirname}', __dirname 
  return settings    


module.exports = {
  load: load
}