import Application from '@head/edge';
import { SOURCE } from './utils/constants';
import Callback from './utils/callback';
import '@shoelace-style/shoelace/dist/components/alert/alert.js'; // eslint-disable-line import/extensions
import '@shoelace-style/shoelace/dist/components/drawer/drawer.js'; // eslint-disable-line import/extensions
// import '@shoelace-style/shoelace/dist/themes/light.css';


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


backstage.toast = function toast(message, variant = 'primary', duration = 2000) {
  // console.log('@head/backstage toast');

  function escapeHtml(html) {
    const div = document.createElement('div');
    div.textContent = html;
    return div.innerHTML;
  }

  const alert = Object.assign(document.createElement('sl-alert'), {
    variant,
    closable: true,
    duration,
    innerHTML: `
      ${escapeHtml(message)}
    `,
  });

  document.body.append(alert);

  return alert.toast();
};


backstage.drawer = function drawer_(html, width = '50') {
  // console.log('@head/backstage drawer');

  const drawer = document.createElement('sl-drawer');
  drawer.label = 'Backstage';
  drawer.style.setProperty('--size', `${width}vw`);

  if (typeof html === 'string') {
    drawer.innerHTML = `
      ${html}
    `;
  // } else if (content instanceof HTMLElement) {
  //   drawer.innerHTML = ``;
  //   drawer.appendChild(content);
  }

  document.body.appendChild(drawer);

  drawer.show();

  drawer.addEventListener('sl-after-hide', () => {
    drawer.remove();
  });
};


window.postMessage({ source: SOURCE, type: 'READY', from: 'FRONTSTAGE' });
