console.log('background.js');

chrome.action.onClicked.addListener((tab) => {
  const queryOptions = { active: true, currentWindow: true };
  chrome.tabs.query(queryOptions).then(function (tabs) {
    var activeTab = tabs[0];
    chrome.tabs.sendMessage(activeTab.id, { action: 'ActionOnClick', payload: [ tab.id, activeTab.id ] }, function (resp) {
      console.log(resp);
    });
  });
});
