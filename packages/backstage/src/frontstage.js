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


backstage.onClick = function onClick(handler) {
  const button = document.createElement('button');
  button.type = 'button';
  button.title = 'Backstage';
  button.innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
      <rect x="3" y="3" width="7" height="7" rx="1"></rect>
      <rect x="14" y="3" width="7" height="7" rx="1"></rect>
      <rect x="3" y="14" width="7" height="7" rx="1"></rect>
      <rect x="14" y="14" width="7" height="7" rx="1"></rect>
    </svg>
  `;
  Object.assign(button.style, {
    position: 'fixed',
    right: '24px',
    bottom: '24px',
    width: '48px',
    height: '48px',
    borderRadius: '50%',
    border: 'none',
    padding: '0',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
    color: '#f5f5f5',
    backgroundColor: '#1f1f1f',
    boxShadow: '0 2px 8px rgba(0, 0, 0, 0.3)',
    zIndex: '2147483646',
    transition: 'background-color 0.15s ease, transform 0.15s ease',
  });
  button.addEventListener('mouseenter', () => {
    button.style.backgroundColor = '#2b2b2b';
    button.style.transform = 'scale(1.05)';
  });
  button.addEventListener('mouseleave', () => {
    button.style.backgroundColor = '#1f1f1f';
    button.style.transform = 'scale(1)';
  });
  button.addEventListener('click', () => {
    try {
      const result = handler();
      if (result && typeof result.then === 'function') {
        result.catch((err) => console.error('[BACKSTAGE] onClick handler error', err));
      }
    } catch (err) {
      console.error('[BACKSTAGE] onClick handler error', err);
    }
  });
  document.body.appendChild(button);
  return button;
};


window.postMessage({ source: SOURCE, type: 'READY', from: 'FRONTSTAGE' });
