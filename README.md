hectorcorrea.com
================
This is the source code of the site that powers my personal site at http://hectorcorrea.com. 

In a nutshell, this site is a home grown mini-blog engine using Node.js, Express.js, and MongoDB.


Requirements
------------
To run this code you need to have **Node.js** and **Mongo DB** 
installed on your machine. 

Last but not least, you'll need to install a few modules that are used in the code. You can do this by executing the following commands in the Terminal 
*from inside the folder where you downloaded the source code*

    cd ~/dev/hectorcorrea.com
    npm install 


How to run the site
-------------------
Download the source code and install the requirements listed above.

Update the settings.dev.json file and make sure the **dbUrl** points to your MongoDB database (e.g. "mongodb://localhost:27017/hectorcorrea"). 

To kick off the application, just run the following command from the Terminal window: 

    node server

...and browse to your *http://localhost:3000* 

When the server connects to the database, *if there are no other users in the database, it will automatically create a default user with the parameters indicated in the settings.dev.json configuration file*. You can login with this user by browsing to *http://localhost:3000/login*


Structure of the source code
----------------------------

Server-side Code

* **server.js** is the main program. 
* **models/** models and database access code.
* **routes/** controllers.
* **views/** contains the views. 


Running the site in production
------------------------------
When you run the site in production you need to pass the connection URL to the database and the information for the default user somehow because the program will not read them from settings.prod.json. Although reading these values from settings.prod.json would have been easier to program I decided against this approach to prevent me (and others) from accidentally pushing these values to GitHub. Instead, the program expect these values in environment variables when it runs in production. 

The first time you run the site in production you can do something like this:

    NODE_ENV=production DB_URL=the_url BLOG_USER=u BLOG_PWD=p node server.js

You only need to pass BLOG_USER and BLOG_PWD the first time since once the default user has been created these values are not needed anymore. Therefore, after the first time you only need to run something like this:

    NODE_ENV=production DB_URL=the_url node server.js

Another way of achieving this is by setting environment variables in your init.d script. For example, I have something similar to this in my production server:

    export NODE_ENV=production
    export PORT=3000
    export DB_URL=[define-value-here]
    #export BLOG_USER=[define-value-here]
    #export BLOG_PASSWORD=[define-value-here]
    export BLOG_SALT=[define-value-here]
    node server.js


Questions, comments, thoughts?
------------------------------
This is a very rough work in progress as I learn and play with Node.js.

Feel free to contact me with questions or comments about this project.

You can see a running version version of this code here:

  [http://hectorcorrea.com](http://hectorcorrea.com)

Keep in mind that you'll need to host to the site on your own in order to be able to add new topics or edit existing ones. 
