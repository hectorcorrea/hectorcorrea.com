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


services.factory('FlashMsg', function($rootScope) {
  var message = '';

  return {
    set: function(newMessage) {
      message = newMessage;
    },
    get: function() {
      var oldMessage = message;
      message = '';
      return oldMessage;
    }
  };
});


// ========================================================
// App Definition
// ========================================================

var hcApp = angular.module('hcApp', ['ngCookies', 'hc.services']);

var routesConfig = function($routeProvider, $locationProvider) {

  // $locationProvider.html5Mode(true);

  $routeProvider.
  when('/', {
    controller: 'EmptyController',
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
    controller: 'EmptyController',
    templateUrl: '/partials/about.html' 
  }).
  when('/credits', {
    controller: 'EmptyController',
    templateUrl: '/partials/credits.html' 
  }).
  when('/user/changePassword', {
    controller: 'ChangePasswordController',
    templateUrl: '/partials/changePassword.html' 
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
    controller: 'EmptyController',
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

hcApp.controller('EmptyController', ['$scope', '$location', '$window', 'FlashMsg',
  function($scope, $location, $window, FlashMsg) {

    $scope.flash = FlashMsg.get();

    $scope.$on('$viewContentLoaded', function(event) {
      // Track the page in Google Analytics
      $window._gaq.push(['_trackPageview', $location.path()]);
    });

  }
]);

hcApp.controller('ListController', ['$scope', '$location', '$window', 'Security', 'Blog', 'entries', 
  function($scope, $location, $window, Security, Blog, entries) {

    $scope.entries = entries;
    $scope.isAuth = Security.isAuth();

    $scope.$on('$viewContentLoaded', function(event) {
      $window._gaq.push(['_trackPageview', $location.path()]);
    });

    $scope.new = function() {
      Blog.createNew(
        function(blog) {
          var editUrl = "/blog/" + blog.url + "/" + blog.key + "/edit";
          $location.path(editUrl);
        },
        function(e) {
          $scope.errorMsg = e.data.message;
        }
      );

    }

  }
]);


hcApp.controller('ViewController', ['$scope', '$http', '$location', '$window', 'Security', 'blog',
  function($scope, $http, $location, $window, Security, blog) {

    $scope.blog = blog;
    $scope.isAuth = Security.isAuth();

    // http://stackoverflow.com/a/10713709/446681
    $scope.$on('$viewContentLoaded', function(event) {
      $window._gaq.push(['_trackPageview', $location.path()]);
    });

    $scope.edit = function() {
      $location.path($scope.blog.editUrl);
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


hcApp.controller('EditController', ['$scope', '$location', '$window', 'Security', 'Blog', 'blog', 
  function($scope, $location, $window, Security, Blog, blog) {

    $scope.blog = blog;
    $scope.isAuth = Security.isAuth();

    $scope.submit = function() {
      var blog = new Blog($scope.blog);
      blog.$save(
        function(b) {
          var viewUrl = "/#/blog/" + b.url + "/"+ b.key;
          // Use $window instead of $location to force a full reload
          // and pick up the updated content
          $window.location.href =  viewUrl;
        },
        function(e) {
          $scope.errorMsg = e.data.message;
        }
      );
    }

    $scope.saveWorkInProgress = function() {
      var blog = new Blog($scope.blog);
      blog.$save(
        function(b) {
          var timestamp = new Date();
          $scope.errorMsg = 'Work in progress saved at ' + timestamp.toTimeString();
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


hcApp.controller('LoginController', ['$scope', '$http', '$location', 'Security', 'FlashMsg',
  function($scope, $http, $location, Security, FlashMsg) {

    $scope.user = '';
    $scope.password = '';
    $scope.isAuth = Security.isAuth();
    $scope.errorMsg = '';
    $scope.flash = FlashMsg.get();

    $scope.logout = function() {
      console.log('About to logout');
      $http.post('/api/logout', {}).
      success(function(data, status) {
        FlashMsg.set("You've been logged out");
        $location.path('/');
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
        FlashMsg.set('Welcome back ' + $scope.user);
        $location.path('/');
      }).
      error(function(data, status) {
        debugger;
        $scope.errorMsg = data;
        console.log('ERROR in login');
      });
    };

  }
]);


hcApp.controller('ChangePasswordController', ['$scope', '$http', '$location', '$cookies', 'Security', 
  function($scope, $http, $location, $cookies, Security) {

    $scope.user = $cookies.user;
    $scope.oldPassword = '';
    $scope.newPassword = '';
    $scope.repeatPassword = '';
    $scope.isAuth = Security.isAuth();
    if (!$scope.isAuth) {
      // WTF, how did you get here buddy?
      $scope.errorMsg = 'Must be logged in first.';
    }

    $scope.changePassword = function() {

      if($scope.newPassword == $scope.repeatPassword) {

        console.log('About to change password');
        var formData = {
          user: $scope.user, 
          oldPassword: $scope.oldPassword,
          newPassword: $scope.newPassword
        }; 
        $http.post('/api/user/changePassword', formData).
        success(function(data, status) {
          $scope.errorMsg = 'Changed password OK';
          // TODO: redirect to home page?
        }).
        error(function(data, status) {
          debugger;
          $scope.errorMsg = data;
          console.log('ERROR changing password');
        });

      }
      else {
        $scope.errorMsg = 'Repeat password does not match with new password';
      }

    };

  }
]);
