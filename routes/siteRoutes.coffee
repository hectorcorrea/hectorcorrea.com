{Logger} = require '../util/logger'

home = (req, res) ->  
  viewModel = {}
  Logger.info "siteRoutes:home"
  res.render 'home', viewModel


about = (req, res) -> 
  viewModel = { title: "About Hector"}
  Logger.info "siteRoutes:about"
  res.render 'about', viewModel


credits = (req, res) -> 
  viewModel = { title: "Credits"}
  Logger.info "siteRoutes:credits"
  res.render 'credits', viewModel


notFound = (req, res) ->
  url = 
  Logger.info "siteRoutes:notFound URL: #{req.originalUrl}"
  # console.dir req
  res.render '404.ejs', { status: 404, message: 'Page not found' }


blowUp = (req, res) ->
  Logger.info "siteRoutes:blowUp (intentional)"
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
