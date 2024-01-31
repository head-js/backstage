console.log('content.js');

const $app = document.createElement('div');
$app.classList.add('HeadBackstageApp');

const $body = document.querySelector('body');
$body.appendChild($app)

chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
  console.log(request);
  console.log(sender);

  if (!request.payload.length) {
    sendResponse({ code: -1, message: 'empty payload' });
  }

  sendResponse({ code: 0, message: 'ok' });
});
