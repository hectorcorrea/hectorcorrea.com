/*global module:false*/
module.exports = function(grunt) {

  // Project configuration.
  grunt.initConfig({
    meta: {
      version: '0.1.0',
      banner: '/*! hectorcorrea.com ' +
        '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
        '* http://hectorcorrea.com/\n' +
        '* Copyright (c) <%= grunt.template.today("yyyy") %> ' +
        'Hector Correa; Licensed MIT */'
    },
    lint: {
      files: ['grunt.js', 'lib/**/*.js', 'test/**/*.js']
    },
    test: {
      files: ['test/**/*.js']
    },
    concat: {
      dist: {
        src: ['<banner:meta.banner>', 
        'public/js/jquery.reveal.js',
        'public/js/jquery.orbit-1.4.0.js',
        'public/js/jquery.customforms.js',
        'public/js/jquery.placeholder.min.js',
        'public/js/jquery.tooltips.js',
        'public/js/app.js',
        'public/js/jquery.offcanvas.js'],
        dest: 'public/js/all.js'
      }
    },
    min: {
      dist: {
        src: ['<banner:meta.banner>', '<config:concat.dist.dest>'],
        dest: 'public/js/all.min.js'
      }
    },
    watch: {
      files: '<config:lint.files>',
      tasks: 'lint test'
    },
    jshint: {
      options: {
        curly: true,
        eqeqeq: true,
        immed: true,
        latedef: true,
        newcap: true,
        noarg: true,
        sub: true,
        undef: true,
        boss: true,
        eqnull: true
      },
      globals: {}
    },
    uglify: {}
  });

  // Default task.
  grunt.registerTask('default', 'lint test concat min');

};
