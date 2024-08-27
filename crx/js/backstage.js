(function () {
  'use strict';

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

  requirejs('vendors/backstage-vendors.js');
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
      chrome.runtime.sendMessage(message, _ref => {
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
    });
  }
  window.addEventListener('message', async _ref2 => {
    let {
      /* type, source, origin, */data
    } = _ref2;
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
      console.debug("%c%s", "background-color: #f0f9ff", '\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
      try {
        const resolved = await callBackground(call);
        callbackFrontstage(callback, resolved, null);
      } catch (rejected) {
        callbackFrontstage(callback, null, rejected);
      }
    } else if (type === 'CALLBACK' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
      console.debug("%c%s", "background-color: #f0f9ff", '\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
      $callback(callback, call.resolved, call.rejected);
    } else {
      console.debug("%c%s", "color: #727272", '\t\t[BACKSTAGE] window.addEventListener.message', data);
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
      console.debug("%c%s", "background-color: #f0f9ff", '\t\t[BACKSTAGE] chrome.runtime.onMessage', from, to, type);
      callFrontstage(call).then(resolved => {
        console.debug("%c%s", "background-color: #f0f9ff", '\t\t[BACKSTAGE] chrome.runtime.onMessage.respond');
        respond({
          resolved,
          rejected: null
        });
      }).catch(rejected => {
        console.debug("%c%s", "background-color: #f0f9ff", '\t\t[BACKSTAGE] chrome.runtime.onMessage.respond');
        respond({
          resolved: null,
          rejected
        });
      });
    } else {
      console.debug("%c%s", "color: #727272", '\t\t[BACKSTAGE] chrome.runtime.onMessage', message);
    }
    return true;
  });

})();
