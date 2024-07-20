import Application from '@head/edge';
import { SOURCE } from './utils/constants';
import Callback from './utils/callback';


const $callback = new Callback();


const { router, client } = new Application();


router.verb('GET', '/readiness', async (ctx, next) => {
  ctx.body = { code: 0, message: 'ok' };
  next();
});


function callBackstage(call) {
  return $callback((id) => {
    const message = { source: SOURCE, type: 'CALL', from: 'FRONTSTAGE', to: 'BACKSTAGE', call, callback: id };
    window.postMessage(message);
  });
}


function callbackBackstage(callback, resolved, rejected) {
  const message = { source: SOURCE, type: 'CALLBACK', from: 'FRONTSTAGE', to: 'BACKSTAGE', call: { resolved, rejected }, callback };
  window.postMessage(message);
}


window.addEventListener('message', async ({ /* type, source, origin, */ data }) => {
  if (data.source !== SOURCE) {
    return false;
  }

  const { from, to, type, call, callback } = data;

  if (type === 'CALLBACK' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
    console.debug('[FRONTSTAGE] window.addEventListener.message', from, to, type);
    $callback(callback, call.resolved, call.rejected);
  } else if (type === 'CALL' && from === 'BACKSTAGE' && to === 'FRONTSTAGE') {
    console.debug('[FRONTSTAGE] window.addEventListener.message', from, to, type);
    // const resolved = { code: 0, message: 'call : backstage -> frontstage : ok' };

    try {
      const resolved = await client.verb(call.method, call.ep, call.search, call.form);
      callbackBackstage(callback, resolved, null);
    } catch (rejected) {
      callbackBackstage(callback, null, rejected);
    }
  } else {
    console.verbose('[FRONTSTAGE] window.addEventListener.message', data);
  }
});


const backstage = window.backstage || {};
window.backstage = backstage;


backstage.invoke = function invoke(svc, method, ep, search = {}, form = {}) {
  const call = { svc, method, ep, search, form };
  return callBackstage(call);
};


backstage.route = router.verb.bind(router);


backstage.verb = client.verb.bind(client);


window.postMessage({ source: SOURCE, type: 'READY', from: 'FRONTSTAGE' });
