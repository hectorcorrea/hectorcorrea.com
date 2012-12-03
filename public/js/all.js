/*! hectorcorrea.com 2012-12-02
* http://hectorcorrea.com/
* Copyright (c) 2012 Hector Correa; Licensed MIT */

/*
 * jQuery Reveal Plugin 1.1
 * www.ZURB.com
 * Copyright 2010, ZURB
 * Free to use under the MIT license.
 * http://www.opensource.org/licenses/mit-license.php
*/
/*globals jQuery */

(function ($) {
  'use strict';
  var modalQueued = false;
  
  $('a[data-reveal-id]').live('click', function (event) {
    event.preventDefault();
    var modalLocation = $(this).attr('data-reveal-id');
    $('#' + modalLocation).reveal($(this).data());
  });

  $.fn.reveal = function (options) {
    var defaults = {
      animation: 'fadeAndPop',                // fade, fadeAndPop, none
      animationSpeed: 300,                    // how fast animtions are
      closeOnBackgroundClick: true,           // if you click background will modal close?
      dismissModalClass: 'close-reveal-modal', // the class of a button or element that will close an open modal
      open: $.noop,
      opened: $.noop,
      close: $.noop,
      closed: $.noop
    };
    options = $.extend({}, defaults, options);

    return this.each(function () {
      var modal    = $(this),
        topMeasure = parseInt(modal.css('top'), 10),
        topOffset  = modal.height() + topMeasure,
        locked     = false,
        modalBg    = $('.reveal-modal-bg'),
        closeButton;

      if (modalBg.length === 0) {
        modalBg = $('<div class="reveal-modal-bg" />').insertAfter(modal);
        modalBg.fadeTo('fast', 0.8);
      }
      
      function unlockModal() {
        locked = false;
      }

      function lockModal() {
        locked = true;
      }
      
      function closeOpenModals(modal) {
        
        var openModals = $(".reveal-modal.open");
        if (openModals.length === 1) {
          modalQueued = true;
          $(".reveal-modal.open").trigger("reveal:close");
        }
      }

      function openAnimation() {
        if (!locked) {
          lockModal();
          closeOpenModals(modal);
          modal.addClass("open");
          if (options.animation === "fadeAndPop") {
            modal.css({'top': $(document).scrollTop() - topOffset, 'opacity': 0, 'visibility': 'visible', 'display' : 'block'});
            modalBg.fadeIn(options.animationSpeed / 2);
            modal.delay(options.animationSpeed / 2).animate({
              "top": $(document).scrollTop() + topMeasure + 'px',
              "opacity": 1
            }, options.animationSpeed, function () {
              modal.trigger('reveal:opened');
            });

          }
          if (options.animation === "fade") {
            modal.css({'opacity': 0, 'visibility': 'visible', 'display' : 'block', 'top': $(document).scrollTop() + topMeasure});
            modalBg.fadeIn(options.animationSpeed / 2);
            modal.delay(options.animationSpeed / 2).animate({
              "opacity": 1
            }, options.animationSpeed, function () {
              modal.trigger('reveal:opened');
            });

          }
          if (options.animation === "none") {
            modal.css({'visibility': 'visible', 'display' : 'block', 'top': $(document).scrollTop() + topMeasure});
            modalBg.css({"display": "block"});
            modal.trigger('reveal:opened');
          }
        }
      }
      modal.bind('reveal:open.reveal', openAnimation);

      function closeAnimation() {
        if (!locked) {
          lockModal();
          modal.removeClass("open");
          if (options.animation === "fadeAndPop") {
            modal.animate({
              "top":  $(document).scrollTop() - topOffset + 'px',
              "opacity": 0
            }, options.animationSpeed / 2, function () {
              modal.css({'top': topMeasure, 'opacity': 1, 'visibility': 'hidden', 'display': 'none'});
            });
            if (!modalQueued) {
              modalBg.delay(options.animationSpeed).fadeOut(options.animationSpeed, function () {
                modal.trigger('reveal:closed');
              });             
            } else {
              modal.trigger('reveal:closed');
            }
            modalQueued = false;
          }
          if (options.animation === "fade") {
            modal.animate({
              "opacity" : 0
            }, options.animationSpeed, function () {
              modal.css({'opacity': 1, 'visibility': 'hidden', 'display' : 'none', 'top': topMeasure});
            });
            if (!modalQueued) {
              modalBg.delay(options.animationSpeed).fadeOut(options.animationSpeed, function () {
                modal.trigger('reveal:closed');
              });
            } else {
              modal.trigger('reveal:closed');
            }
          }
          if (options.animation === "none") {
            modal.css({'visibility': 'hidden', 'display' : 'block', 'top': topMeasure});
            if (!modalQueued) {
              modalBg.css({'display': 'none'});
            }
            modal.trigger('reveal:closed');
          }
        }
      }

      function destroy() {
        modal.unbind('.reveal');
        modalBg.unbind('.reveal');
        $('.' + options.dismissModalClass).unbind('.reveal');
        $('body').unbind('.reveal');
      }

      modal.bind('reveal:close.reveal', closeAnimation);
      modal.bind('reveal:opened.reveal reveal:closed.reveal', unlockModal);
      modal.bind('reveal:closed.reveal', destroy);
      
      modal.bind('reveal:open.reveal', options.open);
      modal.bind('reveal:opened.reveal', options.opened);
      modal.bind('reveal:close.reveal', options.close);
      modal.bind('reveal:closed.reveal', options.closed);
      
      modal.trigger('reveal:open');

      closeButton = $('.' + options.dismissModalClass).bind('click.reveal', function () {
        modal.trigger('reveal:close');
      });

      if (options.closeOnBackgroundClick) {
        modalBg.css({"cursor": "pointer"});
        modalBg.bind('click.reveal', function () {
          modal.trigger('reveal:close');
        });
      }

      $('body').bind('keyup.reveal', function (event) {
        if (event.which === 27) { // 27 is the keycode for the Escape key
          modal.trigger('reveal:close');
        }
      });
    });
  };
} (jQuery));
/*
 * jQuery Orbit Plugin 1.4.0
 * www.ZURB.com/playground
 * Copyright 2010, ZURB
 * Free to use under the MIT license.
 * http://www.opensource.org/licenses/mit-license.php
*/


(function($) {
  'use strict';
  $.fn.findFirstImage = function () {
    return this.first()
            .find('img')
            .andSelf().filter('img')
            .first();
  };

  var ORBIT = {

    defaults: {
      animation: 'horizontal-push',     // fade, horizontal-slide, vertical-slide, horizontal-push, vertical-push
      animationSpeed: 600,        // how fast animtions are
      timer: true,            // true or false to have the timer
      advanceSpeed: 4000,         // if timer is enabled, time between transitions
      pauseOnHover: false,        // if you hover pauses the slider
      startClockOnMouseOut: false,    // if clock should start on MouseOut
      startClockOnMouseOutAfter: 1000,  // how long after MouseOut should the timer start again
      directionalNav: true,         // manual advancing directional navs
      directionalNavRightText: 'Right', // text of right directional element for accessibility
      directionalNavLeftText: 'Left', // text of left directional element for accessibility
      captions: true,           // do you want captions?
      captionAnimation: 'fade',       // fade, slideOpen, none
      captionAnimationSpeed: 600,     // if so how quickly should they animate in
      resetTimerOnClick: false,      // true resets the timer instead of pausing slideshow progress on manual navigation
      bullets: false,           // true or false to activate the bullet navigation
      bulletThumbs: false,        // thumbnails for the bullets
      bulletThumbLocation: '',      // location from this file where thumbs will be
      afterSlideChange: $.noop,   // empty function
      afterLoadComplete: $.noop, //callback to execute after everything has been loaded
      fluid: true,
      centerBullets: true    // center bullet nav with js, turn this off if you want to position the bullet nav manually
    },

    activeSlide: 0,
    numberSlides: 0,
    orbitWidth: null,
    orbitHeight: null,
    locked: null,
    timerRunning: null,
    degrees: 0,
    wrapperHTML: '<div class="orbit-wrapper" />',
    timerHTML: '<div class="timer"><span class="mask"><span class="rotator"></span></span><span class="pause"></span></div>',
    captionHTML: '<div class="orbit-caption"></div>',
    directionalNavHTML: '<div class="slider-nav"><span class="right"></span><span class="left"></span></div>',
    bulletHTML: '<ul class="orbit-bullets"></ul>',

    init: function (element, options) {
      var $imageSlides,
          imagesLoadedCount = 0,
          self = this;

      // Bind functions to correct context
      this.clickTimer = $.proxy(this.clickTimer, this);
      this.addBullet = $.proxy(this.addBullet, this);
      this.resetAndUnlock = $.proxy(this.resetAndUnlock, this);
      this.stopClock = $.proxy(this.stopClock, this);
      this.startTimerAfterMouseLeave = $.proxy(this.startTimerAfterMouseLeave, this);
      this.clearClockMouseLeaveTimer = $.proxy(this.clearClockMouseLeaveTimer, this);
      this.rotateTimer = $.proxy(this.rotateTimer, this);

      this.options = $.extend({}, this.defaults, options);
      if (this.options.timer === 'false') this.options.timer = false;
      if (this.options.captions === 'false') this.options.captions = false;
      if (this.options.directionalNav === 'false') this.options.directionalNav = false;

      this.$element = $(element);
      this.$wrapper = this.$element.wrap(this.wrapperHTML).parent();
      this.$slides = this.$element.children('img, a, div');

      this.$element.bind('orbit.next', function () {
        self.shift('next');
      });

      this.$element.bind('orbit.prev', function () {
        self.shift('prev');
      });

      this.$element.bind('orbit.goto', function (event, index) {
        self.shift(index);
      });

      this.$element.bind('orbit.start', function (event, index) {
        self.startClock();
      });

      this.$element.bind('orbit.stop', function (event, index) {
        self.stopClock();
      });

      $imageSlides = this.$slides.filter('img');

      if ($imageSlides.length === 0) {
        this.loaded();
      } else {
        $imageSlides.bind('imageready', function () {
          imagesLoadedCount += 1;
          if (imagesLoadedCount === $imageSlides.length) {
            self.loaded();
          }
        });
      }
    },

    loaded: function () {
      this.$element
        .addClass('orbit')
        .css({width: '1px', height: '1px'});

      this.$slides.addClass('orbit-slide');

      this.setDimentionsFromLargestSlide();
      this.updateOptionsIfOnlyOneSlide();
      this.setupFirstSlide();

      if (this.options.timer) {
        this.setupTimer();
        this.startClock();
      }

      if (this.options.captions) {
        this.setupCaptions();
      }

      if (this.options.directionalNav) {
        this.setupDirectionalNav();
      }

      if (this.options.bullets) {
        this.setupBulletNav();
        this.setActiveBullet();
      }

      this.options.afterLoadComplete.call(this);
    },

    currentSlide: function () {
      return this.$slides.eq(this.activeSlide);
    },

    setDimentionsFromLargestSlide: function () {
      //Collect all slides and set slider size of largest image
      var self = this,
          $fluidPlaceholder;

      self.$element.add(self.$wrapper).width(this.$slides.first().outerWidth());
      self.$element.add(self.$wrapper).height(this.$slides.first().height());
      self.orbitWidth = this.$slides.first().outerWidth();
      self.orbitHeight = this.$slides.first().height();
      $fluidPlaceholder = this.$slides.first().findFirstImage().clone();


      this.$slides.each(function () {
        var slide = $(this),
            slideWidth = slide.outerWidth(),
            slideHeight = slide.height();

        if (slideWidth > self.$element.outerWidth()) {
          self.$element.add(self.$wrapper).width(slideWidth);
          self.orbitWidth = self.$element.outerWidth();
        }
        if (slideHeight > self.$element.height()) {
          self.$element.add(self.$wrapper).height(slideHeight);
          self.orbitHeight = self.$element.height();
          $fluidPlaceholder = $(this).findFirstImage().clone();
        }
        self.numberSlides += 1;
      });

      if (this.options.fluid) {
        if (typeof this.options.fluid === "string") {
          $fluidPlaceholder = $('<img src="http://placehold.it/' + this.options.fluid + '" />')
        }

        self.$element.prepend($fluidPlaceholder);
        $fluidPlaceholder.addClass('fluid-placeholder');
        self.$element.add(self.$wrapper).css({width: 'inherit'});
        self.$element.add(self.$wrapper).css({height: 'inherit'});

        $(window).bind('resize', function () {
          self.orbitWidth = self.$element.outerWidth();
          self.orbitHeight = self.$element.height();
        });
      }
    },

    //Animation locking functions
    lock: function () {
      this.locked = true;
    },

    unlock: function () {
      this.locked = false;
    },

    updateOptionsIfOnlyOneSlide: function () {
      if(this.$slides.length === 1) {
        this.options.directionalNav = false;
        this.options.timer = false;
        this.options.bullets = false;
      }
    },

    setupFirstSlide: function () {
      //Set initial front photo z-index and fades it in
      var self = this;
      this.$slides.first()
        .css({"z-index" : 3})
        .fadeIn(function() {
          //brings in all other slides IF css declares a display: none
          self.$slides.css({"display":"block"})
      });
    },

    startClock: function () {
      var self = this;

      if(!this.options.timer) {
        return false;
      }

      if (this.$timer.is(':hidden')) {
        this.clock = setInterval(function () {
          self.$element.trigger('orbit.next');
        }, this.options.advanceSpeed);
      } else {
        this.timerRunning = true;
        this.$pause.removeClass('active')
        this.clock = setInterval(this.rotateTimer, this.options.advanceSpeed / 180);
      }
    },

    rotateTimer: function (reset) {
      var degreeCSS = "rotate(" + this.degrees + "deg)"
      this.degrees += 2;
      this.$rotator.css({
        "-webkit-transform": degreeCSS,
        "-moz-transform": degreeCSS,
        "-o-transform": degreeCSS
      });
      if(this.degrees > 180) {
        this.$rotator.addClass('move');
        this.$mask.addClass('move');
      }
      if(this.degrees > 360 || reset) {
        this.$rotator.removeClass('move');
        this.$mask.removeClass('move');
        this.degrees = 0;
        this.$element.trigger('orbit.next');
      }
    },

    stopClock: function () {
      if (!this.options.timer) {
        return false;
      } else {
        this.timerRunning = false;
        clearInterval(this.clock);
        this.$pause.addClass('active');
      }
    },

    setupTimer: function () {
      this.$timer = $(this.timerHTML);
      this.$wrapper.append(this.$timer);

      this.$rotator = this.$timer.find('.rotator');
      this.$mask = this.$timer.find('.mask');
      this.$pause = this.$timer.find('.pause');

      this.$timer.click(this.clickTimer);

      if (this.options.startClockOnMouseOut) {
        this.$wrapper.mouseleave(this.startTimerAfterMouseLeave);
        this.$wrapper.mouseenter(this.clearClockMouseLeaveTimer);
      }

      if (this.options.pauseOnHover) {
        this.$wrapper.mouseenter(this.stopClock);
      }
    },

    startTimerAfterMouseLeave: function () {
      var self = this;

      this.outTimer = setTimeout(function() {
        if(!self.timerRunning){
          self.startClock();
        }
      }, this.options.startClockOnMouseOutAfter)
    },

    clearClockMouseLeaveTimer: function () {
      clearTimeout(this.outTimer);
    },

    clickTimer: function () {
      if(!this.timerRunning) {
          this.startClock();
      } else {
          this.stopClock();
      }
    },

    setupCaptions: function () {
      this.$caption = $(this.captionHTML);
      this.$wrapper.append(this.$caption);
      this.setCaption();
    },

    setCaption: function () {
      var captionLocation = this.currentSlide().attr('data-caption'),
          captionHTML;

      if (!this.options.captions) {
        return false;
      }

      //Set HTML for the caption if it exists
      if (captionLocation) {
        //if caption text is blank, don't show captions
        if ($.trim($(captionLocation).text()).length < 1){
          return false;
        }
        captionHTML = $(captionLocation).html(); //get HTML from the matching HTML entity
        this.$caption
          .attr('id', captionLocation) // Add ID caption TODO why is the id being set?
          .html(captionHTML); // Change HTML in Caption
          //Animations for Caption entrances
        switch (this.options.captionAnimation) {
          case 'none':
            this.$caption.show();
            break;
          case 'fade':
            this.$caption.fadeIn(this.options.captionAnimationSpeed);
            break;
          case 'slideOpen':
            this.$caption.slideDown(this.options.captionAnimationSpeed);
            break;
        }
      } else {
        //Animations for Caption exits
        switch (this.options.captionAnimation) {
          case 'none':
            this.$caption.hide();
            break;
          case 'fade':
            this.$caption.fadeOut(this.options.captionAnimationSpeed);
            break;
          case 'slideOpen':
            this.$caption.slideUp(this.options.captionAnimationSpeed);
            break;
        }
      }
    },

    setupDirectionalNav: function () {
      var self = this,
          $directionalNav = $(this.directionalNavHTML);

      $directionalNav.find('.right').html(this.options.directionalNavRightText);
      $directionalNav.find('.left').html(this.options.directionalNavLeftText);

      this.$wrapper.append($directionalNav);

      this.$wrapper.find('.left').click(function () {
        self.stopClock();
        if (self.options.resetTimerOnClick) {
          self.rotateTimer(true);
          self.startClock();
        }
        self.$element.trigger('orbit.prev');
      });

      this.$wrapper.find('.right').click(function () {
        self.stopClock();
        if (self.options.resetTimerOnClick) {
          self.rotateTimer(true);
          self.startClock();
        }
        self.$element.trigger('orbit.next');
      });
    },

    setupBulletNav: function () {
      this.$bullets = $(this.bulletHTML);
      this.$wrapper.append(this.$bullets);
      this.$slides.each(this.addBullet);
      this.$element.addClass('with-bullets');
      if (this.options.centerBullets) this.$bullets.css('margin-left', -this.$bullets.outerWidth() / 2);
    },

    addBullet: function (index, slide) {
      var position = index + 1,
          $li = $('<li>' + (position) + '</li>'),
          thumbName,
          self = this;

      if (this.options.bulletThumbs) {
        thumbName = $(slide).attr('data-thumb');
        if (thumbName) {
          $li
            .addClass('has-thumb')
            .css({background: "url(" + this.options.bulletThumbLocation + thumbName + ") no-repeat"});;
        }
      }
      this.$bullets.append($li);
      $li.data('index', index);
      $li.click(function () {
        self.stopClock();
        if (self.options.resetTimerOnClick) {
          self.rotateTimer(true);
          self.startClock();
        }
        self.$element.trigger('orbit.goto', [$li.data('index')])
      });
    },

    setActiveBullet: function () {
      if(!this.options.bullets) { return false; } else {
        this.$bullets.find('li')
          .removeClass('active')
          .eq(this.activeSlide)
          .addClass('active');
      }
    },

    resetAndUnlock: function () {
      this.$slides
        .eq(this.prevActiveSlide)
        .css({"z-index" : 1});
      this.unlock();
      this.options.afterSlideChange.call(this, this.$slides.eq(this.prevActiveSlide), this.$slides.eq(this.activeSlide));
    },

    shift: function (direction) {
      var slideDirection = direction;

      //remember previous activeSlide
      this.prevActiveSlide = this.activeSlide;

      //exit function if bullet clicked is same as the current image
      if (this.prevActiveSlide == slideDirection) { return false; }

      if (this.$slides.length == "1") { return false; }
      if (!this.locked) {
        this.lock();
        //deduce the proper activeImage
        if (direction == "next") {
          this.activeSlide++;
          if (this.activeSlide == this.numberSlides) {
              this.activeSlide = 0;
          }
        } else if (direction == "prev") {
          this.activeSlide--
          if (this.activeSlide < 0) {
            this.activeSlide = this.numberSlides - 1;
          }
        } else {
          this.activeSlide = direction;
          if (this.prevActiveSlide < this.activeSlide) {
            slideDirection = "next";
          } else if (this.prevActiveSlide > this.activeSlide) {
            slideDirection = "prev"
          }
        }

        //set to correct bullet
        this.setActiveBullet();

        //set previous slide z-index to one below what new activeSlide will be
        this.$slides
          .eq(this.prevActiveSlide)
          .css({"z-index" : 2});

        //fade
        if (this.options.animation == "fade") {
          this.$slides
            .eq(this.activeSlide)
            .css({"opacity" : 0, "z-index" : 3})
            .animate({"opacity" : 1}, this.options.animationSpeed, this.resetAndUnlock);
        }

        //horizontal-slide
        if (this.options.animation == "horizontal-slide") {
          if (slideDirection == "next") {
            this.$slides
              .eq(this.activeSlide)
              .css({"left": this.orbitWidth, "z-index" : 3})
              .animate({"left" : 0}, this.options.animationSpeed, this.resetAndUnlock);
          }
          if (slideDirection == "prev") {
            this.$slides
              .eq(this.activeSlide)
              .css({"left": -this.orbitWidth, "z-index" : 3})
              .animate({"left" : 0}, this.options.animationSpeed, this.resetAndUnlock);
          }
        }

        //vertical-slide
        if (this.options.animation == "vertical-slide") {
          if (slideDirection == "prev") {
            this.$slides
              .eq(this.activeSlide)
              .css({"top": this.orbitHeight, "z-index" : 3})
              .animate({"top" : 0}, this.options.animationSpeed, this.resetAndUnlock);
          }
          if (slideDirection == "next") {
            this.$slides
              .eq(this.activeSlide)
              .css({"top": -this.orbitHeight, "z-index" : 3})
              .animate({"top" : 0}, this.options.animationSpeed, this.resetAndUnlock);
          }
        }

        //horizontal-push
        if (this.options.animation == "horizontal-push") {
          if (slideDirection == "next") {
            this.$slides
              .eq(this.activeSlide)
              .css({"left": this.orbitWidth, "z-index" : 3})
              .animate({"left" : 0}, this.options.animationSpeed, this.resetAndUnlock);
            this.$slides
              .eq(this.prevActiveSlide)
              .animate({"left" : -this.orbitWidth}, this.options.animationSpeed);
          }
          if (slideDirection == "prev") {
            this.$slides
              .eq(this.activeSlide)
              .css({"left": -this.orbitWidth, "z-index" : 3})
              .animate({"left" : 0}, this.options.animationSpeed, this.resetAndUnlock);
            this.$slides
              .eq(this.prevActiveSlide)
              .animate({"left" : this.orbitWidth}, this.options.animationSpeed);
          }
        }

        //vertical-push
        if (this.options.animation == "vertical-push") {
          if (slideDirection == "next") {
            this.$slides
              .eq(this.activeSlide)
              .css({top: -this.orbitHeight, "z-index" : 3})
              .animate({top : 0}, this.options.animationSpeed, this.resetAndUnlock);
            this.$slides
              .eq(this.prevActiveSlide)
              .animate({top : this.orbitHeight}, this.options.animationSpeed);
          }
          if (slideDirection == "prev") {
            this.$slides
              .eq(this.activeSlide)
              .css({top: this.orbitHeight, "z-index" : 3})
              .animate({top : 0}, this.options.animationSpeed, this.resetAndUnlock);
            this.$slides
              .eq(this.prevActiveSlide)
              .animate({top : -this.orbitHeight}, this.options.animationSpeed);
          }
        }

        this.setCaption();
      }
    }
  };

  $.fn.orbit = function (options) {
    return this.each(function () {
      var orbit = $.extend({}, ORBIT);
      orbit.init(this, options);
    });
  };

})(jQuery);

/*!
 * jQuery imageready Plugin
 * http://www.zurb.com/playground/
 *
 * Copyright 2011, ZURB
 * Released under the MIT License
 */
(function ($) {

  var options = {};

  $.event.special.imageready = {

    setup: function (data, namespaces, eventHandle) {
      options = data || options;
    },

    add: function (handleObj) {
      var $this = $(this),
          src;

      if ( this.nodeType === 1 && this.tagName.toLowerCase() === 'img' && this.src !== '' ) {
        if (options.forceLoad) {
          src = $this.attr('src');
          $this.attr('src', '');
          bindToLoad(this, handleObj.handler);
          $this.attr('src', src);
        } else if ( this.complete || this.readyState === 4 ) {
          handleObj.handler.apply(this, arguments);
        } else {
          bindToLoad(this, handleObj.handler);
        }
      }
    },

    teardown: function (namespaces) {
      $(this).unbind('.imageready');
    }
  };

  function bindToLoad(element, callback) {
    var $this = $(element);

    $this.bind('load.imageready', function () {
       callback.apply(element, arguments);
       $this.unbind('load.imageready');
     });
  }

}(jQuery));
/*
 * jQuery Custom Forms Plugin 1.0
 * www.ZURB.com
 * Copyright 2010, ZURB
 * Free to use under the MIT license.
 * http://www.opensource.org/licenses/mit-license.php
*/

jQuery.foundation = jQuery.foundation || {};
jQuery.foundation.customForms = jQuery.foundation.customForms || {};

jQuery(document).ready(function ($) {
    
  $.foundation.customForms.appendCustomMarkup = function (options) {
    var defaults = {
      disable_class: "js-disable-custom"
    };
    var options = $.extend(defaults, options);
    
    function appendCustomMarkup(idx, sel) {
      var $this = $(sel).hide(),
          type  = $this.attr('type'),
          $span = $this.next('span.custom.' + type);

      if ($span.length === 0) {
        $span = $('<span class="custom ' + type + '"></span>').insertAfter($this);
      }

      $span.toggleClass('checked', $this.is(':checked'));
      $span.toggleClass('disabled', $this.is(':disabled'));
    }

    function appendCustomSelect(idx, sel) {
      var $this = $(sel),
          $customSelect = $this.next('div.custom.dropdown'),
          $options = $this.find('option'),
          $seloptions = $this.find('option:selected'),
          maxWidth = 0,
          $li;

      if ($this.hasClass('no-custom')) { return; }
      if ($customSelect.length === 0) {
        $customSelectSize = '';
        if ($(sel).hasClass('small')) {
          $customSelectSize = 'small';
        } else if ($(sel).hasClass('medium')) {
          $customSelectSize = 'medium';
        } else if ($(sel).hasClass('large')) {
          $customSelectSize = 'large';
        } else if ($(sel).hasClass('expand')) {
          $customSelectSize = 'expand';
        }
        $customSelect = $('<div class="custom dropdown ' + $customSelectSize + '"><a href="#" class="selector"></a><ul></ul></div>"');
        $options.each(function () {
          $li = $('<li>' + $(this).html() + '</li>');
          $customSelect.find('ul').append($li);
        });
        $customSelect.prepend('<a href="#" class="current">' + $seloptions.html() + '</a>');

        $this.after($customSelect);
        $this.hide();

      } else {
        // refresh the ul with options from the select in case the supplied markup doesn't match
        $customSelect.find('ul').html('');
        $options.each(function () {
          $li = $('<li>' + $(this).html() + '</li>');
          $customSelect.find('ul').append($li);
        });
      }

      $customSelect.toggleClass('disabled', $this.is(':disabled'));

      $options.each(function (index) {
        if (this.selected) {
          $customSelect.find('li').eq(index).addClass('selected');
          $customSelect.find('.current').html($(this).html());
        }
      });

      $customSelect.css('width', 'inherit');
      $customSelect.find('ul').css('width', 'inherit');

      $customSelect.find('li').each(function () {
        $customSelect.addClass('open');
        if ($(this).outerWidth() > maxWidth) {
          maxWidth = $(this).outerWidth();
        }
        $customSelect.removeClass('open');
      });

      if (!$customSelect.is('.small, .medium, .large, .expand')) {
        $customSelect.css('width', maxWidth + 18 + 'px');
        $customSelect.find('ul').css('width', maxWidth + 16 + 'px');
      }

    }
    
    $('form.custom input:radio[data-customforms!=disabled]').each(appendCustomMarkup);
    $('form.custom input:checkbox[data-customforms!=disabled]').each(appendCustomMarkup);
    $('form.custom select[data-customforms!=disabled]').each(appendCustomSelect);
  };

});

(function ($) {
  
  function refreshCustomSelect($select) {
    var maxWidth = 0,
        $customSelect = $select.next();
    $options = $select.find('option');
    $customSelect.find('ul').html('');
    
    $options.each(function () {
      $li = $('<li>' + $(this).html() + '</li>');
      $customSelect.find('ul').append($li);
    });
    
    // re-populate
    $options.each(function (index) {
      if (this.selected) {
        $customSelect.find('li').eq(index).addClass('selected');
        $customSelect.find('.current').html($(this).html());
      }
    });
    
    // fix width
    $customSelect.removeAttr('style')
      .find('ul').removeAttr('style');
    $customSelect.find('li').each(function () {
      $customSelect.addClass('open');
      if ($(this).outerWidth() > maxWidth) {
        maxWidth = $(this).outerWidth();
      }
      $customSelect.removeClass('open');
    });
    $customSelect.css('width', maxWidth + 18 + 'px');
    $customSelect.find('ul').css('width', maxWidth + 16 + 'px');
    
  }
  
  function toggleCheckbox($element) {
    var $input = $element.prev(),
        input = $input[0];

    if (false == $input.is(':disabled')) {
        input.checked = ((input.checked) ? false : true);
        $element.toggleClass('checked');

        $input.trigger('change');
    }
  }
  
  function toggleRadio($element) {
    var $input = $element.prev(),
        input = $input[0];

    if (false == $input.is(':disabled')) {
      $('input:radio[name="' + $input.attr('name') + '"]').each(function () {
        $(this).next().removeClass('checked');
      });
      input.checked = ((input.checked) ? false : true);
      $element.toggleClass('checked');
    
      $input.trigger('change');
    }
  }
  
  $('form.custom span.custom.checkbox').live('click', function (event) {
    event.preventDefault();
    event.stopPropagation();
    
    toggleCheckbox($(this));
  });
  
  $('form.custom span.custom.radio').live('click', function (event) {
    event.preventDefault();
    event.stopPropagation();
    
    toggleRadio($(this));
  });
  
  $('form.custom select').live('change', function (event) {
    refreshCustomSelect($(this));
  });
  
  $('form.custom label').live('click', function (event) {
    var $associatedElement = $('#' + $(this).attr('for')),
        $customCheckbox,
        $customRadio;
    if ($associatedElement.length !== 0) {
      if ($associatedElement.attr('type') === 'checkbox') {
        event.preventDefault();
        $customCheckbox = $(this).find('span.custom.checkbox');
        toggleCheckbox($customCheckbox);
      } else if ($associatedElement.attr('type') === 'radio') {
        event.preventDefault();
        $customRadio = $(this).find('span.custom.radio');
        toggleRadio($customRadio);
      }
    }
  });

  $('form.custom div.custom.dropdown a.current, form.custom div.custom.dropdown a.selector').live('click', function (event) {
    var $this = $(this),
        $dropdown = $this.closest('div.custom.dropdown'),
        $select = $dropdown.prev();
    
    event.preventDefault();
    $('div.dropdown').removeClass('open');

    if (false == $select.is(':disabled')) {
        $dropdown.toggleClass('open');

        if ($dropdown.hasClass('open')) {
          $(document).bind('click.customdropdown', function (event) {
            $dropdown.removeClass('open');
            $(document).unbind('.customdropdown');
          });
        } else {
          $(document).unbind('.customdropdown');
        }
        return false;
    }
  });
  
  $('form.custom div.custom.dropdown li').live('click', function (event) {
    var $this = $(this),
        $customDropdown = $this.closest('div.custom.dropdown'),
        $select = $customDropdown.prev(),
        selectedIndex = 0;
        
    event.preventDefault();
    event.stopPropagation();
    $('div.dropdown').removeClass('open');
    
    $this
      .closest('ul')
      .find('li')
      .removeClass('selected');
    $this.addClass('selected');
    
    $customDropdown
      .removeClass('open')
      .find('a.current')
      .html($this.html());
    
    $this.closest('ul').find('li').each(function (index) {
      if ($this[0] == this) {
        selectedIndex = index;
      }
      
    });
    $select[0].selectedIndex = selectedIndex;
    
    $select.trigger('change');
  });
})(jQuery);
/*! http://mths.be/placeholder v2.0.7 by @mathias */
;(function(f,h,$){var a='placeholder' in h.createElement('input'),d='placeholder' in h.createElement('textarea'),i=$.fn,c=$.valHooks,k,j;if(a&&d){j=i.placeholder=function(){return this};j.input=j.textarea=true}else{j=i.placeholder=function(){var l=this;l.filter((a?'textarea':':input')+'[placeholder]').not('.placeholder').bind({'focus.placeholder':b,'blur.placeholder':e}).data('placeholder-enabled',true).trigger('blur.placeholder');return l};j.input=a;j.textarea=d;k={get:function(m){var l=$(m);return l.data('placeholder-enabled')&&l.hasClass('placeholder')?'':m.value},set:function(m,n){var l=$(m);if(!l.data('placeholder-enabled')){return m.value=n}if(n==''){m.value=n;if(m!=h.activeElement){e.call(m)}}else{if(l.hasClass('placeholder')){b.call(m,true,n)||(m.value=n)}else{m.value=n}}return l}};a||(c.input=k);d||(c.textarea=k);$(function(){$(h).delegate('form','submit.placeholder',function(){var l=$('.placeholder',this).each(b);setTimeout(function(){l.each(e)},10)})});$(f).bind('beforeunload.placeholder',function(){$('.placeholder').each(function(){this.value=''})})}function g(m){var l={},n=/^jQuery\d+$/;$.each(m.attributes,function(p,o){if(o.specified&&!n.test(o.name)){l[o.name]=o.value}});return l}function b(m,n){var l=this,o=$(l);if(l.value==o.attr('placeholder')&&o.hasClass('placeholder')){if(o.data('placeholder-password')){o=o.hide().next().show().attr('id',o.removeAttr('id').data('placeholder-id'));if(m===true){return o[0].value=n}o.focus()}else{l.value='';o.removeClass('placeholder');l==h.activeElement&&l.select()}}}function e(){var q,l=this,p=$(l),m=p,o=this.id;if(l.value==''){if(l.type=='password'){if(!p.data('placeholder-textinput')){try{q=p.clone().attr({type:'text'})}catch(n){q=$('<input>').attr($.extend(g(this),{type:'text'}))}q.removeAttr('name').data({'placeholder-password':true,'placeholder-id':o}).bind('focus.placeholder',b);p.data({'placeholder-textinput':q,'placeholder-id':o}).before(q)}p=p.removeAttr('id').hide().prev().attr('id',o).show()}p.addClass('placeholder');p[0].value=p.attr('placeholder')}else{p.removeClass('placeholder')}}}(this,document,jQuery));
/*
 * jQuery Foundation Tooltip Plugin 2.0.1
 * http://foundation.zurb.com
 * Copyright 2012, ZURB
 * Free to use under the MIT license.
 * http://www.opensource.org/licenses/mit-license.php
*/

;(function($) {
  'use strict';
  var settings = {
    bodyHeight : 0,
    targetClass : '.has-tip',
    tooltipClass : '.tooltip',
    tipTemplate : function (selector, content) {
      return '<span data-selector="' + selector + '" class="' + settings.tooltipClass.substring(1) + '">' + content + '<span class="nub"></span></span>';
    }
  },
  methods = {
    init : function (options) {
      return this.each(function () {
        var $body = $('body'),
          self = this;

        if (Modernizr.touch) {
          $body.on('click.tooltip touchstart.tooltip touchend.tooltip', settings.targetClass, function (e) {
            e.preventDefault();
            var $this = $(this);
            $(settings.tooltipClass).hide();
            methods.showOrCreateTip($this);
          });
          $(settings.tooltipClass).on('click.tooltip touchstart.tooltip touchend.tooltip', function (e) {
            e.preventDefault();
            $(this).fadeOut(150);
          });
        } else {
          $body.on('mouseover.tooltip mouseout.tooltip', settings.targetClass, function (e) {
            var $this = $(this);
            if (e.type === 'mouseover') {
              methods.showOrCreateTip($this);
            } else if (e.type === 'mouseout') {
              methods.hide($this);
            }
          });
        }

      });
    },
    showOrCreateTip : function ($target) {
      var $tip = methods.getTip($target);
      if ($tip && $tip.length > 0) {
        methods.show($target);
      } else {
        methods.create($target);
      }
    },
    getTip : function ($target) {
      var selector = methods.selector($target),
        tip = null;
      if (selector) tip = $('span[data-selector=' + selector + ']' + settings.tooltipClass);
      return (tip) ? tip : false;
    },
    selector : function ($target) {
      var id = $target.attr('id'),
        dataSelector = $target.data('selector');
      if (id === undefined && dataSelector === undefined) {
        dataSelector = 'tooltip' + Math.random().toString(36).substring(7);
        $target.attr('data-selector', dataSelector);
      }
      return (id) ? id : dataSelector;
    },
    create : function ($target) {
      var $tip = $(settings.tipTemplate(methods.selector($target), $target.attr('title'))),
          classes = methods.inheritable_classes($target);
      $tip.addClass(classes).appendTo('body');
      if (Modernizr.touch) $tip.append('<span class="tap-to-close">tap to close </span>');
      $target.removeAttr('title');
      methods.show($target);
    },
    reposition : function (target, tip, classes) {
      var width, nub, nubHeight, nubWidth, row, objPos;

      tip.css('visibility', 'hidden').show();

      width = target.data('width');
      nub = tip.children('.nub');
      nubHeight = nub.outerHeight();
      nubWidth = nub.outerWidth();

      objPos = function (obj, top, right, bottom, left, width) {
        return obj.css({
          'top' : top,
          'bottom' : bottom,
          'left' : left,
          'right' : right,
          'width' : (width) ? width : 'auto'
        }).end();
      };

      objPos(tip, (target.offset().top + target.outerHeight() + 10), 'auto', 'auto', target.offset().left, width);
      objPos(nub, -nubHeight, 'auto', 'auto', 10);

      if ($(window).width() < 767) {
        row = target.parents('.row');
        tip.width(row.outerWidth() - 20).css('left', row.offset().left).addClass('tip-override');
        objPos(nub, -nubHeight, 'auto', 'auto', target.offset().left);
      } else {
        if (classes.indexOf('tip-top') > -1) {
          objPos(tip, (target.offset().top - tip.outerHeight() - nubHeight), 'auto', 'auto', target.offset().left, width)
            .removeClass('tip-override');
          objPos(nub, 'auto', 'auto', -nubHeight, 'auto');
        } else if (classes.indexOf('tip-left') > -1) {
          objPos(tip, (target.offset().top + (target.outerHeight() / 2) - nubHeight), 'auto', 'auto', (target.offset().left - tip.outerWidth() - 10), width)
            .removeClass('tip-override');
          objPos(nub, (tip.outerHeight() / 2) - (nubHeight / 2), -nubHeight, 'auto', 'auto');
        } else if (classes.indexOf('tip-right') > -1) {
          objPos(tip, (target.offset().top + (target.outerHeight() / 2) - nubHeight), 'auto', 'auto', (target.offset().left + target.outerWidth() + 10), width)
            .removeClass('tip-override');
          objPos(nub, (tip.outerHeight() / 2) - (nubHeight / 2), 'auto', 'auto', -nubHeight);
        }
      }
      tip.css('visibility', 'visible').hide();
    },
    inheritable_classes : function (target) {
      var inheritables = ['tip-top', 'tip-left', 'tip-bottom', 'tip-right', 'noradius'],
      filtered = target.attr('class').split(' ').map(function (el, i) {
        if ($.inArray(el, inheritables) !== -1) {
          return el;
        }
      }).join(' ');
      return $.trim(filtered);
    },
    show : function ($target) {
      var $tip = methods.getTip($target);
      methods.reposition($target, $tip, $target.attr('class'));
      $tip.fadeIn(150);
    },
    hide : function ($target) {
      var $tip = methods.getTip($target);
      $tip.fadeOut(150);
    },
    reload : function () {
      var $self = $(this);
      return ($self.data('tooltips')) ? $self.tooltips('destroy').tooltips('init') : $self.tooltips('init');
    },
    destroy : function () {
      return this.each(function () {
        $(window).off('.tooltip');
        $(settings.targetClass).off('.tooltip');
        $(settings.tooltipClass).each(function(i) {
          $($(settings.targetClass).get(i)).attr('title', $(this).text());
        }).remove();
      });
    }
  };

  $.fn.tooltips = function (method) {
    if (methods[method]) {
      return methods[method].apply(this, Array.prototype.slice.call(arguments, 1));
    } else if (typeof method === 'object' || ! method) {
      return methods.init.apply(this, arguments);
    } else {
      $.error('Method ' +  method + ' does not exist on jQuery.tooltips');
    }
  };
})(jQuery);