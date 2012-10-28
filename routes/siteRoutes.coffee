home = (req, res) ->  
	viewModel = {}
	res.render 'home', viewModel


about = (req, res) -> 
	viewModel = { title: "About Hector"}
	res.render 'about', viewModel


credits = (req, res) -> 
  viewModel = { title: "Credits"}
  res.render 'credits', viewModel


notFound = (req, res) ->
  res.render '404.ejs', { status: 404, message: 'Page not found' }


blowUp = (req, res) ->
  a = b.x
  viewModel = { title: "About Hector"}
  res.render 'about', viewModel


module.exports = {
  home: home
  about: about
  notFound: notFound
  credits: credits
  blowUp: blowUp
}
