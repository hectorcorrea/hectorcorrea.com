home = (req, res) ->  
	viewModel = {}
	res.render 'home', viewModel


about = (req, res) -> 
	viewModel = { title: "About Hector"}
	res.render 'about', viewModel


notFound = (req, res) ->
  res.render '404.ejs', { status: 404, message: 'Page not found' }


module.exports = {
  home: home
  about: about
  notFound: notFound
}
