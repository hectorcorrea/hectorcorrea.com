// ========================================================
// Services
// ========================================================
var services = angular.module('hc.services', ['ngResource']);

services.factory('Blog', ['$resource', 
  function($resource) {
    var url = '/api/blog/:url/:key';
    var params = {key: '@key', url: '@url'};
    
    var methods = {};
    
    methods.all = {
      method:'GET', 
      params:{url:'all', key:''}, 
      isArray:true
    };

    methods.createNew = {
      method:'POST', 
      params:{url:'new'}, 
      isArray:false
    };    

    return $resource(url, params, methods);
  }
]);


services.factory('SingleBlog', ['Blog', '$route', '$q',
  function(Blog, $route, $q) {
    return function(decode) {

      var delay = $q.defer();
      var query = {
        url: $route.current.params.url,
        key: $route.current.params.key
      };

      if(decode) {
        query.decode = true;
      }

      var ok = function(blog) {
        blog.editUrl = "/blog/" + blog.url + "/" + blog.key + "/edit";
        delay.resolve(blog);
      };

      var error = function() {
        delay.reject('Unable to fetch blog ' + query.key);
      };

      Blog.get(query, ok, error);
      return delay.promise;
    }
  }
]);


services.factory('ListBlogs', ['Blog', '$route', '$q',
  function(Blog, $route, $q) {
    return function() {

      var delay = $q.defer();

      var ok = function(entries) { 
        delay.resolve(entries); 
      };

      var error = function() { 
        delay.reject('Unable to fetch blog entries'); 
      };

      Blog.all(ok, error);
      return delay.promise;
    }
  }
]);


services.factory('Security', ['$cookies', 
  function($cookies) { 
    return {
      isAuth: function() { return ($cookies.authToken !== undefined); } 
    };
  }
]);


// ========================================================
// App Definition
// ========================================================

var hcApp = angular.module('hcApp', ['ngCookies', 'hc.services']);

var routesConfig = function($routeProvider, $locationProvider) {

  // $locationProvider.html5Mode(true);

  $routeProvider.
  when('/', {
    templateUrl: '/partials/home.html'   
  }).
  when('/blog', {
    controller: 'ListController',
    resolve: {
      entries: function(ListBlogs) { return ListBlogs(); }
    },
    templateUrl: '/partials/blogList.html' 
  }).
  when('/blog/:url/:key/edit', {
    controller: 'EditController',
    resolve: {
      blog: function(SingleBlog) { return SingleBlog(); }
    },
    templateUrl: '/partials/blogEdit.html' 
  }).
  when('/blog/:url/:key', {
    controller: 'ViewController',
    resolve: {
      blog: function(SingleBlog) { return SingleBlog(); }
    },
    templateUrl: '/partials/blogView.html' 
  }).
  when('/about', {
    templateUrl: '/partials/about.html' 
  }).
  when('/credits', {
    templateUrl: '/partials/credits.html' 
  }).
  when('/user/changePassword', {
    controller: 'ChangePasswordController',
    templateUrl: '/partials/changePassword.html' 
  }).

  when('/login/init', {
    controller: 'LoginController',
    templateUrl: '/partials/loginInit.html' 
  }).
  when('/logout', {
    controller: 'LoginController',
    templateUrl: '/partials/logout.html' 
  }).
  when('/login', {
    controller: 'LoginController',
    templateUrl: '/partials/login.html' 
  }).
  otherwise({
    templateUrl: '/partials/notFound.html' 
  });
}

hcApp.config(routesConfig);


// This global var is used to preserve the last search.
// I should be using $rootScope for this but I had a few
// timing issues setting the $rootScope in a controller  
// and then reading its values on the Service at a later
// time. For now use good old global vars.
var globalSearch = {text: null, data: null};


// ========================================================
// Controllers
// ========================================================

hcApp.controller('ListController', ['$scope', '$location', 'Security', 'Blog', 'entries', 
  function($scope, $location, Security, Blog, entries) {

    $scope.entries = entries;
    $scope.isAuth = Security.isAuth();

    $scope.new = function() {
      Blog.createNew(
        function(blog) {
          var editUrl = "/blog/" + blog.url + "/" + blog.key + "/edit";
          $location.url(editUrl);
        },
        function(e) {
          $scope.errorMsg = e.data.message;
        }
      );

    }

  }
]);


hcApp.controller('ViewController', ['$scope', '$http', '$location', 'Security', 'blog',
  function($scope, $http, $location, Security, blog) {

    $scope.blog = blog;
    $scope.isAuth = Security.isAuth();

    $scope.edit = function() {
      $location.url($scope.blog.editUrl);
    };

    // TODO: Move this code to the service
    $scope.draft = function() {
      var draftUrl = '/api/blog/' + $scope.blog.url + '/' + $scope.blog.key + '/draft';
      $http.post(draftUrl, {}).
      success(function(data, status) {
        $scope.blog.postedOn = null;
      }).
      error(function(data, status) {
        console.log("ERROR: not marked as draft");
      });
    };

    // TODO: Move this code to the service
    $scope.post = function() {
      var postUrl = '/api/blog/' + $scope.blog.url + '/' + $scope.blog.key + '/post';
      $http.post(postUrl, {}).
      success(function(data, status) {
        $scope.blog.postedOn = data.postedOn;
      }).
      error(function(data, status) {
        console.log("ERROR: not marked as posted");
      });
    };

  }
]);


hcApp.controller('EditController', ['$scope', '$location', 'Security', 'Blog', 'blog', 
  function($scope, $location, Security, Blog, blog) {

    $scope.blog = blog;
    $scope.isAuth = Security.isAuth();

    $scope.submit = function() {
      var blog = new Blog($scope.blog);
      blog.$save(
        function(b) {
          var viewUrl = "/blog/" + b.url + "/"+ b.key;
          $location.url(viewUrl);
        },
        function(e) {
          $scope.errorMsg = e.data.message;
        }
      );
    }

  }
]);


hcApp.controller('RecipeSearchController', ['$scope', '$routeParams', 'Security', 'Recipe', 'recipes',
  function($scope, $routeParams, Security, Recipe, recipes) {

    $scope.recipes = recipes;
    $scope.searchText = globalSearch.text;
    $scope.message = "";
    $scope.errorMsg = null;
    $scope.isAuth = Security.isAuth();

    $scope.search = function() {

      Recipe.query(
        {text: $scope.searchText}, 
        function(recipes) {

          $scope.message = "";
          $scope.recipes = recipes;
          $scope.errorMsg = null;
          globalSearch.text = $scope.searchText;
          globalSearch.data = recipes;

          if(recipes.length == 0) {
            $scope.message = "No recipes were found"
          }
          else {
            // Give the focus to another element so that
            // the keyboard presented by phones and tables
            // disappears.
            // This should probably go as an Angular Directive
            // rather than manipulating the DOM here but
            // we'll leave that for another day.
            var btn = document.getElementById("btnSearch");
            if(btn) btn.focus();
          }

        }, 
        function(e) {

          globalSearch.text = null;
          globalSearch.data = null;

          $scope.errorMsg = e.message + "/" + e.details;
          console.log($scope.errorMsg);

        }
      );

    }

  }
]);


hcApp.controller('LoginController', ['$scope', '$http', '$location', 'Security', 
  function($scope, $http, $location, Security) {

    $scope.user = '';
    $scope.password = '';
    $scope.isAuth = Security.isAuth();

    $scope.init = function() {
      console.log('About to init login');
      $http.post('/api/login/initialize', {}).
      success(function(data, status) {
        $scope.errorMsg = 'Initialized';
      }).
      error(function(data, status) {
        debugger;
        $scope.errorMsg = data;
        console.log('ERROR in login init');
      });
    };

    $scope.logout = function() {
      console.log('About to logout');
      $http.post('/api/logout', {}).
      success(function(data, status) {
        $scope.errorMsg = 'Logged out';
      }).
      error(function(data, status) {
        debugger;
        $scope.errorMsg = data;
        console.log('ERROR logging out');
      });
    };

    $scope.login = function() {
      console.log('About to login');
      $http.post('/api/login', {user: $scope.user, password: $scope.password}).
      success(function(data, status) {
        $scope.errorMsg = 'Logged in OK';
        // TODO: redirect to home page?
      }).
      error(function(data, status) {
        debugger;
        $scope.errorMsg = data;
        console.log('ERROR in login');
      });
    };

  }
]);


hcApp.controller('ChangePasswordController', ['$scope', '$http', '$location', 'Security', 
  function($scope, $http, $location, Security) {

    $scope.user = '';
    $scope.password = '';
    $scope.isAuth = Security.isAuth();

    $scope.changePassword = function() {
      console.log('About to change password');
      $http.post('/api/user/changePassword', {user: $scope.user, password: $scope.password}).
      success(function(data, status) {
        $scope.errorMsg = 'Changed password OK';
        // TODO: redirect to home page?
      }).
      error(function(data, status) {
        debugger;
        $scope.errorMsg = data;
        console.log('ERROR changing password');
      });
    };

  }
]);
