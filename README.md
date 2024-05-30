# Backstage

Tabs can call other tabs, awaiting promised callback, via *magic*.

<img src="./README-magic.jpg" width="360" />

```javascript
// Tab 1: https://tab1.com
const res = await backstage.invoke('https://tab2.com/*', 'GET', '/api/ping');

// Tab 2: https://tab2.com
backstage.route('GET', '/api/ping', (ctx, next) => {
  ctx.body = { code: 0, message: 'pong' };
  next();
});
```

```mermaid
sequenceDiagram
    participant Tab1 / Frontstage
    participant Tab1 / Backstage
    participant Background
    participant Tab2 / Backstage
    participant Tab2 / Frontstage
    Tab1 / Frontstage ->> Tab1 / Backstage: window.postMessage
    Tab1 / Backstage ->> Background: chrome.runtime.sendMessage
    Background ->> Background: chrome.tabs.query : by targe origin
    Background ->> Tab2 / Backstage: chrome.tabs.sendMessage
    Tab2 / Backstage ->> Tab2 / Frontstage: window.postMessage
    Tab2 / Frontstage ->> Tab2 / Backstage: window.postMessage:callback
    Tab2 / Backstage ->> Background: chrome.tabs.sendMessage::respond
    Background ->> Tab1 / Backstage: chrome.runtime.sendMessage::respond
    Tab1 / Backstage ->> Tab1 / Frontstage: window.postMessage::callback
```
