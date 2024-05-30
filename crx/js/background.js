(function () {
  'use strict';

  function ignore(message) {
    let extra1 = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : '';
    console.log('%c%s', 'color: #727272', message, extra1);
  }
  function debug(message) {
    let extra1 = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : '';
    let extra2 = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : '';
    let extra3 = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : '';
    console.log('%c%s', 'background-color: #f0f9ff', message, extra1, extra2, extra3);
  }

  // logger.debug('background.js');

  function callBackstage(call) {
    return new Promise((resolve, reject) => {
      chrome.tabs.query({
        url: call.svc
      }, tabs => {
        if (tabs.length > 0) {
          const message = {
            type: 'CALL',
            from: 'BACKGROUND',
            to: 'BACKSTAGE',
            call
          };
          chrome.tabs.sendMessage(tabs[0].id, message, _ref => {
            let {
              resolved,
              rejected
            } = _ref;
            if (rejected) {
              reject(rejected);
            } else {
              resolve(resolved);
            }
          });
        }
      });
    });
  }
  chrome.runtime.onMessage.addListener((message, sender, respond) => {
    const {
      from,
      to,
      type,
      call
    } = message;
    if (type === 'CALL' && from === 'BACKSTAGE' && to === 'BACKGROUND') {
      debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage', from, to, type);
      callBackstage(call).then(resolved => {
        debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage.respond');
        respond({
          resolved,
          rejected: null
        });
      }).catch(rejected => {
        debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage.respond');
        respond({
          resolved: null,
          rejected
        });
      });
    } else {
      ignore('\t\t\t\t[BACKGROUND] chrome.runtime.onMessage', message);
    }
    return true;
  });

})();
