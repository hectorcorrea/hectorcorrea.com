hectorcorrea.com
================
This *used to be* the source code of the site that powers my personal site at http://hectorcorrea.com but it is not anymore. The site is now powered by the code in https://github.com/hectorcorrea/hectorcorrea.github.io

Notice I am pretty new to Go and this might not follow Go recommended
practices but I am using it as my sandbox and it is running in production.

How to run the site
------------

```
# Get the code
git clone https://github.com/hectorcorrea/hectorcorrea.com.git
cd hectorcorrea.com

# Compile the code
go build  

# and run it with the default sample configuration
source .env_sample
./hectorcorrea.com

# browse to localhost:9001
```

Once the site is running
--------
Browse to http://localhost:9001/auth/login to login. Use user `user1` password
`welcome1`

Then go to http://localhost:9001/blog and click `new` to add a new blog.



Structure of the source code
----------------------------
* **main.go** launches the web server
* **web/** routes requests to the proper models.
* **models/** connect to the database.
* **views/** contains the views.


The database
--------------
Data will be stored in text files via the [textodb](https://github.com/hectorcorrea/textodb) package.

By default data is stored in `./data/` but you can configure this via the environment variable `DB_ROOT_DIR`

The first time the server is run it will automatically create a `user.xml` file with the user name and password indicated in the following environment variables, the values in parenthesis are the default values if you don't specify them.

* BLOG_USR (user1)
* BLOG_PASS (welcome1)
* BLOG_SALT ()

You can see where these values are used in `models/user.go`


Questions, comments, thoughts?
------------------------------
This is a very rough work in progress as I learn and play with Go.

Feel free to contact me with questions or comments about this project.
