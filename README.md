hectorcorrea.com
================
This is the source code of the site that powers my personal site [http://hectorcorrea.com](http://hectorcorrea.com)

In a nutshell, this site is a home grown mini-blog engine using Node.js, Express.js, Angular.js, and MongoDB (the so called MEAN stack)


Requirements
------------
To run this code you need to have **Node.js** and **Mongo DB** 
installed on your machine. 

Last but not least, you'll need to install **Express**, **EJS**, and a few other utilities by running the following command from the Terminal 
*from inside the folder where you downloaded the source code*

    cd ~/dev/hectorcorrea.com
    npm install 

Express is an MVC-like JavaScript framework that takes care of the boiler plate code to handle HTTP requests and responses. [More info](http://expressjs.com)

EJS is a template engine for Node.js used to generate HTML pages with dynamic content. [More info](https://github.com/visionmedia/ejs)


How to run the site
-------------------
Download the source code and install the requirements listed above.

Update the settings.dev.json file and make sure the **dbUrl** points to your MongoDB database. 

To kick off the application, just run the following command from the Terminal window: 

    node server

...and browse to your *http://localhost:3000* 


Structure of the source code
----------------------------
TBD

Questions, comments, thoughts?
------------------------------
This is a very rough work in progress as I learn and play with Node.js.

Feel free to contact me with questions or comments about this project.

You can see a running version version of this code here:

  [http://54.200.135.65](http://54.200.135.65/)

Keep in mind that you'll need to host to the site on your own in order to be able to add new topics or edit existing ones. 
