hectorcorrea.com
================
This is the source code of the site that powers my personal site at http://hectorcorrea.com.

This branch (go_version) is my work in progress as I transition the code from
Node.js/Express.js/MongoDB to Go/MySQL and **it is not in production yet**

Notice I am pretty new to Go and this might not follow Go recommended
practices but I am using it as my sandbox.

How to run the site
------------

```
# Get the code
git clone <this repo>
cd hectorcorrea.com
git checkout go_version

# Create the MySQL database
mysql -u root < misc/createdb.sql

# Compile and run it
go build  
go get github.com/go-sql-driver/mysql
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

Please ignore the `*.js` files, they are left overs as I transition
the code from Node.js to Go.


The database
--------------
The code will connect to a MySQL database with the parameters indicated in the
following environment variables. If you don't set these environment variables
the code will assume the value indicated inside the parenthesis.

* DB_USER (root)
* DB_PASSWORD ()
* DB_NAME (blogdb)

You can see where these values are used in `models/db.go`

Questions, comments, thoughts?
------------------------------
This is a very rough work in progress as I learn and play with Go.

Feel free to contact me with questions or comments about this project.
