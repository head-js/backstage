import * as logger from './utils/logger';
import Callback from './utils/callback';


logger.debug('frontstage.js');


const $callback = new Callback();


function callBackstage(call) {
  return $callback((id) => {
    const message = { type: 'CALL', from: 'FRONTSTAGE', to: 'BACKSTAGE', call, callback: id };
    window.postMessage(message);
  });
}


function callbackBackstage(callback, resolved, rejected) {
  const message = { type: 'CALLBACK', from: 'FRONTSTAGE', to: 'BACKSTAGE', call: { resolved, rejected }, callback };
  window.postMessage(message);
}


window.addEventListener('message', ({ /* type, source, origin, */ data }) => {
  if (data.source === 'react-devtools-content-script') {
    return false;
  }

  const { from, to, type, call, callback } = data;

  if (type === 'CALLBACK' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
    logger.debug('[FRONTSTAGE] window.addEventListener.message', from, to, type);
    $callback(callback, call.resolved, call.rejected);
  } else if (type === 'CALL' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
    logger.debug('[FRONTSTAGE] window.addEventListener.message', from, to, type);
    const resolved = { code: 0, message: 'call : backstage -> frontstage : ok' };
    const rejected = null;

    callbackBackstage(callback, resolved, rejected);
  } else {
    logger.ignore('[FRONTSTAGE] window.addEventListener.message', data);
  }
});


const backstage = window.backstage || {};
window.backstage = backstage;


backstage.invoke = function invoke(svc, method, ep, search = {}, form = {}) {
  const call = { svc, method, ep, search, form };
  return callBackstage(call);
};
