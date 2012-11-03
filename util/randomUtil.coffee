randomNumber = (max, min) -> 
  # http://stackoverflow.com/a/1527834/446681
  r = Math.floor(Math.random() * (max - min + 1)) + min


randomUpperChar = ->
  n = randomNumber(65, 90); # A..Z
  String.fromCharCode(n)


randomString = (length = 10) ->
  string = ""
  for i in [1..length]
    string += randomUpperChar()
  return string

module.exports = {
  randomNumber: randomNumber
  randomUpperChar: randomUpperChar
  randomString: randomString
}