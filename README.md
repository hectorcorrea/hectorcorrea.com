hectorcorrea.com
================
This is the source code of the site that powers my personal site at http://hectorcorrea.com. 

In a nutshell, this site is a home grown mini-blog engine using Node.js, Express.js, Angular.js, and MongoDB.


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


Structure of the source code
----------------------------

Server-side Code

* **server.js** is the main server-side program. 
* **models/** contains the server-side models and database access code.
* **routes/** contains the server-side controllers.
* **views/** contains the server-side views. There is really only on server-side view (index.ejs) since the rest of the content is loaded via Angular views.

Client-side Code

* **public/js/app.js** is the main client-side program. This is the where Angular is configured. I am also storing in here the Angular client-side controllers. Ideally they should be on their own JavaScript file but I have not split them.
* **public/js/partials/** contains the views loaded (client-side) by Angular

You can also take a look at the diagrams in the [docs folder](https://github.com/hectorcorrea/hectorcorrea.com/tree/master/docs) to get an idea on how the different components work together.


Questions, comments, thoughts?
------------------------------
This is a very rough work in progress as I learn and play with Node.js.

Feel free to contact me with questions or comments about this project.

You can see a running version version of this code here:

  [http://hectorcorrea.com](http://hectorcorrea.com)

Keep in mind that you'll need to host to the site on your own in order to be able to add new topics or edit existing ones. 
