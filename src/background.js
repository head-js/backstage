function callBackstage(call) {
  return new Promise((resolve, reject) => {
    chrome.tabs.query({ url: call.svc }, (tabs) => {
      if (tabs.length > 0) {
        const message = { type: 'CALL', from: 'BACKGROUND', to: 'BACKSTAGE', call };
        chrome.tabs.sendMessage(tabs[0].id, message, ({ resolved, rejected }) => {
          if (rejected) {
            reject(rejected);
          } else {
            resolve(resolved);
          }
        });
      }
    });
  });
}


chrome.runtime.onMessage.addListener((message, sender, respond) => {
  const { from, to, type, call } = message;

  if (type === 'CALL' && from === 'BACKSTAGE' && to === 'BACKGROUND') {
    console.debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage', from, to, type);
    callBackstage(call)
      .then((resolved) => {
        console.debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage.respond');
        respond({ resolved, rejected: null });
      })
      .catch((rejected) => {
        console.debug('\t\t\t\t[BACKGOURD] chrome.runtime.onMessage.respond');
        respond({ resolved: null, rejected });
      });
  } else {
    console.verbose('\t\t\t\t[BACKGROUND] chrome.runtime.onMessage', message);
  }

  return true;
});
