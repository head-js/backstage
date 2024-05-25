import * as logger from './utils/logger';


logger.debug('background.js');


function callBackstage(call) {
  return new Promise((resolve, reject) => { // eslint-disable-line no-unused-vars
    chrome.tabs.query({ url: call.svc }, (tabs) => {
      if (tabs.length > 0) {
        const message = { type: 'CALL', from: 'BACKGROUND', to: 'BACKSTAGE', call };
        chrome.tabs.sendMessage(tabs[0].id, message, (resolved, rejected) => {
          resolve({ resolved, rejected });
        });
      }
    });
  });
}


chrome.runtime.onMessage.addListener((message, sender, respond) => {
  const { from, to, type, call } = message;

  if (type === 'CALL' && from === 'BACKSTAGE' && to === 'BACKGROUND') {
    logger.debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage', from, to, type);
    callBackstage(call).then(({ resolved, rejected }) => {
      logger.debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage.respond');
      respond(resolved, rejected);
    });
  } else {
    logger.ignore('\t\t\t\t[BACKGROUND] chrome.runtime.onMessage', message);
  }

  return true;
});
