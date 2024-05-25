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
        source: SOURCE,
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
      source: SOURCE,
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
    if (type === 'CALLBACK' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
      debug('[FRONTSTAGE] window.addEventListener.message', from, to, type);
      $callback(callback, call.resolved, call.rejected);
    } else if (type === 'CALL' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
      debug('[FRONTSTAGE] window.addEventListener.message', from, to, type);
      const resolved = {
        code: 0,
        message: 'call : backstage -> frontstage : ok'
      };
      const rejected = null;
      callbackBackstage(callback, resolved, rejected);
    } else {
      ignore('[FRONTSTAGE] window.addEventListener.message', data);
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
