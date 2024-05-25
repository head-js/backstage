(function () {
  'use strict';

  function ignore(message, extra1) {
    console.log('%c%s', 'color: #727272', message, extra1);
  }
  function debug(message) {
    console.log('%c%s', 'background-color: #f0f9ff', message);
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

  debug('frontstage.js');
  const $callback = new Callback();
  function callBackstage(call) {
    return $callback(id => {
      const message = {
        type: 'CALL',
        from: 'FRONTSTAGE',
        to: 'BACKSTAGE',
        call,
        callback: id
      };
      window.postMessage(message);
    });
  }
  function callbackBackstage(callback, resolved, rejected) {
    const message = {
      type: 'CALLBACK',
      from: 'FRONTSTAGE',
      to: 'BACKSTAGE',
      call: {
        resolved,
        rejected
      },
      callback
    };
    window.postMessage(message);
  }
  window.addEventListener('message', _ref => {
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
    if (type === 'CALLBACK' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
      debug('[FRONTSTAGE] window.addEventListner.message');
      $callback(callback, call.resolved, call.rejected);
    } else if (type === 'CALL' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
      debug('[FRONTSTAGE] window.addEventListner.message');
      const resolved = {
        code: 0,
        message: 'call : backstage -> frontstage : ok'
      };
      const rejected = null;
      callbackBackstage(callback, resolved, rejected);
    } else {
      ignore('[FRONTSTAGE] window.addEventListner.message', data);
    }
  });
  const backstage = window.backstage || {};
  window.backstage = backstage;
  backstage.invoke = function invoke(svc, method, ep) {
    let search = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
    let form = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : {};
    const call = {
      svc,
      method,
      ep,
      search,
      form
    };
    return callBackstage(call);
  };

})();
