import * as logger from './utils/logger';
import requirejs from './utils/requirejs';
import Callback from './utils/callback';


logger.debug('backstage.js');


requirejs('js/frontstage.js');


const $callback = new Callback();


function callFrontstage(call) {
  return $callback((id) => {
    const message = { type: 'CALL', from: 'BACKSTAGE', to: 'FRONTSTAGE', call, callback: id };
    window.postMessage(message);
  });
}


function callbackFrontstage(callback, resolved, rejected) {
  const message = { type: 'CALLBACK', from: 'BACKSTAGE', to: 'FRONTSTAGE', call: { resolved, rejected }, callback };
  window.postMessage(message);
}


async function callBackground(call) {
  const message = { type: 'CALL', from: 'BACKSTAGE', to: 'BACKGROUND', call };
  return new Promise((resolve, reject) => { // eslint-disable-line no-unused-vars
    chrome.runtime.sendMessage(message, (resolved, rejected) => {
      resolve({ resolved, rejected });
    });
  });
}


window.addEventListener('message', async ({ /* type, source, origin, */ data }) => {
  if (data.source === 'react-devtools-content-script') {
    return false;
  }

  const { from, to, type, call, callback } = data;

  if (type === 'CALL' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
    logger.debug('\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
    const { resolved, rejected } = await callBackground(call);
    callbackFrontstage(callback, resolved, rejected);
  } else if (type === 'CALLBACK' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
    logger.debug('\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
    $callback(callback, call.resolved, call.rejected);
  } else {
    logger.ignore('\t\t[BACKSTAGE] window.addEventListener.message', data);
  }
});


chrome.runtime.onMessage.addListener((message, sender, respond) => {
  const { from, to, type, call } = message;

  if (type === 'CALL' && from === 'BACKGROUND' && to === 'BACKSTAGE') {
    logger.debug('\t\t[BACKSTAGE] chrome.runtime.onMessage', from, to, type);
    callFrontstage(call).then((resolved, rejected) => {
      respond(resolved, rejected);
    });
  } else {
    logger.ignore('\t\t[BACKSTAGE] chrome.runtime.onMessage', message);
  }

  return true;
});
