var home = function(req, res) {
  res.send({title:'home'});
}

var about = function(req, res) {
  res.send({title:'about'});
}

module.exports = {
  home: home,
  about: about
}
