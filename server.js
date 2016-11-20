
// Blog Routes
app.get('/blog', authenticate, blogRoutes.viewAll)
app.post('/blog/:url/:key/edit', authenticate, blogRoutes.edit);
app.post('/blog/:url/:key/save', authenticate, blogRoutes.save);
app.post('/blog/:url/:key/post', authenticate, blogRoutes.post);
app.post('/blog/:url/:key/draft', authenticate, blogRoutes.draft);
app.get('/blog/:url/:key', authenticate, blogRoutes.viewOne);
app.get('/blog/rss', blogRoutes.rss);
app.get('/blog/:url', legacyRoutes.blogOne);
app.post('/blog/new', authenticate, blogRoutes.newBlog)

app.get('/changePassword', authenticate, userRoutes.changePassword)
app.post('/changePassword', authenticate, userRoutes.changePasswordPost)
