simple-blog
===========
A very simple blog web site using Node.js, CoffeeScript, and Express.js

This is a rough work in progress.

At this point there is code to handle a home page, an about page, a way to see a list of topics, and click on each topic to view the details of the post.

The data is kept in text files for now.


Requirements
------------
To run this code you need to have **Node.js** installed on your machine. If you don't have Node.js you can get it from [nodejs.org](http://nodejs.org)

In addition to Node.js you'll need **CoffeeScript**. Once you've installed Node.js you can easily install CoffeeScript from the Terminal with the following command: 

    npm install -g coffee-script

CoffeeScript is a language that compiles to JavaScript. It allow us to write the code with a cleaner syntax than raw JavaScript. [More info](http://coffeescript.org)

Last but not least, you'll need to install **Express** and **EJS** by running the following command from the Terminal *from inside the folder where you downloaded the source code*

    cd ~/dev/simple-blog
    npm install 

Express is MVC-like JavaScript framework that takes care of the boiler plate code to handle HTTP requests and responses. [More info](http://expressjs.com)

EJS is a template engine for Node.js used to generate HTML pages with dynamic content. [More info](https://github.com/visionmedia/ejs)


How to run the site
-------------------
Once you've downloaded the source code and installed the requirements listed above, just run the following command from the Terminal window: 

    coffee app 

...and browse to your *http://localhost:3000* You should see the rather anti-climactic web site with the beginnings of what will eventually be a blog engine. Enjoy it!


Structure of the source code
----------------------------
**app.coffee** is the main file. The rest of the code is organized as follow:

* **\models:** CoffeeScript files with the models
* **\routes:** CoffeeScript files with the controllers
* **\views:** The views of the project (HTML + EJS)
* **\public:** Static public files (client side JavaScript, CSS)
* **\util:** Several small JavaScript utilities used throught the code
* **\data:** Text files with the sample data
* **\data_test:** Text files created when running the unit tests
* **\logs:** Text files with the information logged when the app runs

The **unit tests** for the models and the controllers are on the same 
folder as the respective functionality (i.e. inside models and routes.) 
You can run the tests for each class individuall 
(e.g. coffee topicModelTest.coffee or coffee topicDataTest.coffee)

There is also a Bash shell script on the root folder (testSuite.sh) that 
runs all the unit tests at once. 


Questions, comments, thoughts?
------------------------------
This is a very rough work in progress as I learn and play with Node.js.

Feel free to contact me with questions or comments about this project.

You can see a running version version of this code here:

  [http://simple-blog.jitsu.com](http://simple-blog.jitsu.com)

Keep in mind that you'll need to host to the site on your own in order 
to be able to add new topics or edit existing ones. 