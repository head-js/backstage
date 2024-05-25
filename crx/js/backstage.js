(function () {
  'use strict';

  function ignore(message, extra1) {
    console.log('%c%s', 'color: #727272', message, extra1);
  }
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

  function seq() {
    return Date.now() + '-' + Math.random().toString(36).substring(2);
  }
  function Callback() {
    const promised = {};
    return function callback(fn, resolved, rejected) {
      if (typeof fn === 'function') {
        const id = seq();
        return new Promise((resolve, reject) => {
          promised[id] = {
            resolve,
            reject
          };
          fn(id);
        });
      } else {
        // typeof fn === 'string'
        const {
          resolve,
          reject
        } = promised[fn];
        if (rejected) {
          reject(rejected);
        } else {
          resolve(resolved);
        }
        delete promised[fn];
      }
    };
  }

  debug('backstage.js');
  requirejs('js/frontstage.js');
  const $callback = new Callback();
  function callFrontstage(call) {
    return $callback(id => {
      const message = {
        type: 'CALL',
        from: 'BACKSTAGE',
        to: 'FRONTSTAGE',
        call,
        callback: id
      };
      window.postMessage(message);
    });
  }
  function callbackFrontstage(callback, resolved, rejected) {
    const message = {
      type: 'CALLBACK',
      from: 'BACKSTAGE',
      to: 'FRONTSTAGE',
      call: {
        resolved,
        rejected
      },
      callback
    };
    window.postMessage(message);
  }
  async function callBackground(call) {
    const message = {
      type: 'CALL',
      from: 'BACKSTAGE',
      to: 'BACKGROUND',
      call
    };
    return new Promise((resolve, reject) => {
      // eslint-disable-line no-unused-vars
      chrome.runtime.sendMessage(message, (resolved, rejected) => {
        resolve({
          resolved,
          rejected
        });
      });
    });
  }
  window.addEventListener('message', async _ref => {
    let {
      /* type, source, origin, */data
    } = _ref;
    if (data.source === 'react-devtools-content-script') {
      return false;
    }
    const {
      from,
      to,
      type,
      call,
      callback
    } = data;
    if (type === 'CALL' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
      debug('\t\t[BACKSTAGE] window.addEventListner.message');
      const {
        resolved,
        rejected
      } = await callBackground(call);
      callbackFrontstage(callback, resolved, rejected);
    } else if (type === 'CALLBACK' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
      debug('\t\t[BACKSTAGE] window.addEventListner.message');
      $callback(callback, call.resolved, call.rejected);
    } else {
      ignore('\t\t[BACKSTAGE] window.addEventListener.message', data);
    }
  });
  chrome.runtime.onMessage.addListener((message, sender, respond) => {
    const {
      from,
      to,
      type,
      call
    } = message;
    if (type === 'CALL' && from === 'BACKGROUND' && to === 'BACKSTAGE') {
      debug('\t\t[BACKSTAGE] chrome.runtime.onMessage');
      callFrontstage(call).then((resolved, rejected) => {
        respond(resolved, rejected);
      });
    } else {
      ignore('\t\t[BACKSTAGE] chrome.runtime.onMessage', message);
    }
    return true;
  });

})();
