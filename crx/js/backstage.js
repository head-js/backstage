(function () {
  'use strict';

  function debug(message) {
    console.log('%c%s', 'background-color: #f0f9ff', message);
  }

  function requirejs(src) {
    const s = document.createElement('script');
    s.src = chrome.runtime.getURL(src);
    // s.type = 'module';
    s.onload = () => s.remove();
    (document.head || document.documentElement).append(s);
  }

  debug('backstage.js');
  requirejs('js/frontstage.js');

})();
