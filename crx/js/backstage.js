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

  const SOURCE = '@head/backstage'; // eslint-disable-line import/prefer-default-export

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
        source: SOURCE,
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
      source: SOURCE,
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
    if (data.source !== SOURCE) {
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
      debug('\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
      const {
        resolved,
        rejected
      } = await callBackground(call);
      callbackFrontstage(callback, resolved, rejected);
    } else if (type === 'CALLBACK' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
      debug('\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
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
      debug('\t\t[BACKSTAGE] chrome.runtime.onMessage', from, to, type);
      callFrontstage(call).then((resolved, rejected) => {
        debug('\t\t[BACKSTAGE] chrome.runtime.onMessage.respond');
        respond(resolved, rejected);
      });
    } else {
      ignore('\t\t[BACKSTAGE] chrome.runtime.onMessage', message);
    }
    return true;
  });

})();
