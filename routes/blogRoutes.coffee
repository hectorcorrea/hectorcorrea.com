fs = require 'fs'
{Logger} = require '../util/logger'
{TopicModel} = require '../models/topicModel'
authModel = require '../models/authModel'

_normalizeTopicTitle = (title) ->
	title = title.trim().toLowerCase()
	title = title.replace('.aspx', '')
	title


# encode &lt; and &gt; as &amp;lt and &amp;gt;
_encodeContent = (content) -> 
  content.replace(/&gt;/g, '&amp;gt;').replace(/&lt;/g, '&amp;lt;')


# restore &amp;lt and &amp;gt; to &lt; and &gt; 
_decodeContent = (content) -> 
  content.replace(/&amp;gt;/g, '&gt;').replace(/&amp;lt;/g, '&lt;')


renderNotFound = (res, error) -> 
	Logger.error "renderNotFound #{error}"
	res.render '404', {status: 404, message: error}


renderError = (res, error) ->
	Logger.error "renderError #{error}"
	res.render '500', {status: 500, message: error}


viewModelForTopic = (topic, isReadOnly) ->
	{ 
		topic: topic
		page:
			title: topic?.meta?.title 
			isReadOnly: isReadOnly
	}


viewModelForTopics = (topics, title, isReadOnly) ->
	{
		matrix: topicsToMatrix(topics, 3),
		page:
			title: title 
			isReadOnly: isReadOnly
	}


topicsToMatrix = (topics, cols) ->
	rows = Math.ceil(topics.length / cols)
	matrix = new Array(rows)
	for r in [0..rows-1]
		matrix[r] = new Array(cols)
		for c in [0..cols-1]
			i = (r * cols) + c
			if i < topics.length
				matrix[r][c] = topics[i]
	matrix


requestToTopic = (req, id) ->
	topic = {
		meta: {
			id: parseInt(id)
			title: req.body?.title ? ""
			summary: req.body?.summary ? ""
		}
		content: req.body?.content ? ""
	}
	return topic
#postedOn: new Date(req.body?.postedOn + ' ' + req.body?.postedAt)


viewOne = (req, res) -> 

	dataPath = req.app.get('datapath')
	dataOptions = res.app.settings.dataOptions
	model = new TopicModel dataOptions 
	url = req.params.topicUrl
	Logger.info "blogRoutes:viewOne #{url}"

	if url
		model.getOneByUrl url, (err, topic) ->  
			if err
				normalizedTitle = _normalizeTopicTitle(url)
				model.getOneByUrl normalizedTitle, (err2, topic2) ->
					if err2
						renderNotFound res, err
					else
						Logger.info "Redirecting to #{normalizedTitle}"
						res.redirect '/blog/' + normalizedTitle, 301
			else
				topic.content = _decodeContent(topic.content)
				isReadOnly = authModel.isAuthenticated(req, dataPath) is false
				viewModel = viewModelForTopic(topic, isReadOnly)
				# console.dir viewModel
				res.render 'blogOne', viewModel
	else
		# we shouldn't get here
		Logger.warn "viewOne without a URL was detected"
		viewRecent()


viewRecent = (req, res) -> 
	Logger.info "blogRoutes:viewRecent"

	dataPath = req.app.get('datapath')
	dataOptions = res.app.settings.dataOptions
	model = new TopicModel dataOptions 

	model.getRecent (err, topics) -> 
		if err 
			Logger.error err
			renderError res, "Error getting recent topics"
		else
			isReadOnly = authModel.isAuthenticated(req, dataPath) is false
			viewModel = viewModelForTopics topics, "Recent Blog Posts", isReadOnly
			res.render 'blogRecent', viewModel


viewAll = (req, res) -> 
	Logger.info "blogRoutes:viewAll"

	dataPath = req.app.get('datapath')
	dataOptions = res.app.settings.dataOptions
	model = new TopicModel dataOptions 

	model.getAll (err, topics) -> 
		if err 
			Logger.error err
			renderError res, "Error getting topics"
		else
			isReadOnly = authModel.isAuthenticated(req, dataPath) is false
			viewModel = viewModelForTopics topics, "All Blog Posts", isReadOnly
			res.render 'blogAll', viewModel


rssList = (req, res) -> 
	Logger.info "blogRoutes:rssList"

	dataOptions = res.app.settings.dataOptions
	model = new TopicModel dataOptions 

	model.getRssList (err, xml) -> 
		if err 
			Logger.error err
			renderError res, "Error getting topics"
		else
			res.send xml, { 'Content-Type': 'application/atom+xml' }, 200


edit = (req, res) -> 

	url = req.params.topicUrl
	if url is undefined
		Logger.warn 'Edit without a URL was detected. Redirecting to blog list.'
		res.redirect '/blog'
		return

	dataPath = req.app.get('datapath')
	if authModel.isAuthenticated(req, dataPath) is false
		Logger.warn "Unauthenticated user attempted to edit topic #{url}"
		res.redirect '/blog'
		return

	Logger.info "blogRoutes:edit #{url}"

	dataOptions = res.app.settings.dataOptions
	model = new TopicModel dataOptions 

	model.getOneByUrl url, (err, topic) -> 
		if err
			Logger.error err
			renderNotFound res, err
		else 
			viewModel = viewModelForTopic(topic, false)
			# console.dir viewModel
			res.render 'blogEdit', viewModel


save = (req, res) -> 

	id = req.params.id
	if id is undefined
		Logger.warn 'Save without an Id was detected. Redirecting to blog list.'
		res.redirect '/blog'
		return
	
	dataPath = req.app.get('datapath')
	if authModel.isAuthenticated(req, dataPath) is false
		Logger.warn "Unauthenticated user attempted to save topic #{id}"
		res.redirect '/blog'
		return

	Logger.info "blogRoutes:save #{id}"

	topic = requestToTopic req, id
	if isNaN(topic.meta.id)
		renderError res, "Invalid id #{id} detected on save."
	else
		# console.dir topic
		isFinal = if req.body?.final then true else false
		topic.meta.postedOn = if isFinal then new Date() else null
		topic.content = _encodeContent(topic.content)
		dataOptions = res.app.settings.dataOptions
		model = new TopicModel dataOptions 
		model.save topic, (err, savedTopic) -> 
			if err
				# Unexpected error, send user to blogs main page
				Logger.error "Error while saving: #{err}"
				res.redirect '/blog'
			else if typeof savedTopic.errors isnt 'undefined'
				# Validation error, send user to edit this topic
				Logger.info "Validation errors detected"
				console.dir savedTopic
				res.render 'blogEdit', viewModelForTopic(savedTopic, false)
			else
				Logger.info "Saved, redirecting to /blog/#{savedTopic.meta.url}"
				res.redirect '/blog/'+ savedTopic.meta.url


editNew = (req, res) ->

	dataPath = req.app.get('datapath')
	if authModel.isAuthenticated(req, dataPath) is false
		Logger.warn "Unauthenticated user attempted to add new topic"
		res.redirect '/blog'
		return

	Logger.info "blogRoutes:editNew"
	dataOptions = res.app.settings.dataOptions
	model = new TopicModel dataOptions 
	topic = model.getNew()
	res.render 'blogEdit', viewModelForTopic(topic, false)


saveNew = (req, res) -> 

	dataPath = req.app.get('datapath')
	if authModel.isAuthenticated(req, dataPath) is false
		Logger.warn "Unauthenticated user attempted to save new topic"
		res.redirect '/blog'
		return

	Logger.info "blogRoutes:saveNew"
	id = null
	topic = requestToTopic req, id
	isFinal = if req.body?.final then true else false
	topic.meta.postedOn = if isFinal then new Date() else null
	topic.content = _encodeContent(topic.content)
	dataOptions = res.app.settings.dataOptions
	model = new TopicModel dataOptions 

	model.saveNew topic, (err, savedTopic) ->
		if err
			# Unexpected error, send user to blogs main page
			Logger.error "Error while saving #{err}"
			res.redirect '/blog'
		else if typeof savedTopic.errors isnt 'undefined'
			# Validation error, send user to edit this topic
			# savedTopic is in the form {meta: X, content: Y, errors: Z}
			Logger.info "Validation errors detected"
			console.dir savedTopic
			res.render 'blogEdit', viewModelForTopic(savedTopic, false)
		else
			Logger.info "New topic added, redirecting to /blog/#{savedTopic.meta.url}"
			res.redirect '/blog/'+ savedTopic.meta.url


module.exports = {
	viewOne: viewOne,
	viewRecent: viewRecent,
	viewAll: viewAll,
	edit: edit,
	save: save,
	editNew: editNew,
	saveNew: saveNew,
	rssList: rssList
	_normalizeTopicTitle: _normalizeTopicTitle
}

