import { SOURCE } from './utils/constants';
import requirejs from './utils/requirejs';
import Callback from './utils/callback';


requirejs('js/frontstage.js');


const $callback = new Callback();


function callFrontstage(call) {
  return $callback((id) => {
    const message = { source: SOURCE, type: 'CALL', from: 'BACKSTAGE', to: 'FRONTSTAGE', call, callback: id };
    window.postMessage(message);
  });
}


function callbackFrontstage(callback, resolved, rejected) {
  const message = { source: SOURCE, type: 'CALLBACK', from: 'BACKSTAGE', to: 'FRONTSTAGE', call: { resolved, rejected }, callback };
  window.postMessage(message);
}


async function callBackground(call) {
  const message = { type: 'CALL', from: 'BACKSTAGE', to: 'BACKGROUND', call };
  return new Promise((resolve, reject) => {
    chrome.runtime.sendMessage(message, ({ resolved, rejected }) => {
      if (rejected) {
        reject(rejected);
      } else {
        resolve(resolved);
      }
    });
  });
}


window.addEventListener('message', async ({ /* type, source, origin, */ data }) => {
  if (data.source !== SOURCE) {
    return false;
  }

  const { from, to, type, call, callback } = data;

  if (type === 'CALL' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
    console.debug('\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
    try {
      const resolved = await callBackground(call);
      callbackFrontstage(callback, resolved, null);
    } catch (rejected) {
      callbackFrontstage(callback, null, rejected);
    }
  } else if (type === 'CALLBACK' && from === 'FRONTSTAGE' && to === 'BACKSTAGE') {
    console.debug('\t\t[BACKSTAGE] window.addEventListener.message', from, to, type);
    $callback(callback, call.resolved, call.rejected);
  } else {
    console.verbose('\t\t[BACKSTAGE] window.addEventListener.message', data);
  }
});


chrome.runtime.onMessage.addListener((message, sender, respond) => {
  const { from, to, type, call } = message;

  if (type === 'CALL' && from === 'BACKGROUND' && to === 'BACKSTAGE') {
    console.debug('\t\t[BACKSTAGE] chrome.runtime.onMessage', from, to, type);
    callFrontstage(call)
      .then((resolved) => {
        console.debug('\t\t[BACKSTAGE] chrome.runtime.onMessage.respond');
        respond({ resolved, rejected: null });
      })
      .catch((rejected) => {
        console.debug('\t\t[BACKSTAGE] chrome.runtime.onMessage.respond');
        respond({ resolved: null, rejected });
      });
  } else {
    console.verbose('\t\t[BACKSTAGE] chrome.runtime.onMessage', message);
  }

  return true;
});
