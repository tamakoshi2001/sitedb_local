chrome.commands.onCommand.addListener((command) => {
  if (command === "getLink") {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      let currentTab = tabs[0];
      if (currentTab) {
        let url = currentTab.url;
        console.log(url);
        // send url to server
        let urlToSend = "http://localhost:8080/site";
        fetch(urlToSend, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ url: url })
        })
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.text();
        })
        .then(data => {
          console.log(data);
          // コマンド成功の通知を表示
          chrome.notifications.create({
            type: 'basic',
            iconUrl: 'icon48.png',
            title: 'URL Sent',
            message: `The URL was successfully sent: ${url}`
          });
        })
        .catch(error => {
          console.log('There was a problem with the fetch operation:', error.message);
          chrome.notifications.create({
            type: 'basic',
            iconUrl: 'icon48.png',
            title: 'Error',
            message: `Failed to send the URL: ${error.message}`
          });
        });
      }
    });
  }

  if (command === "removeLink") {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      let currentTab = tabs[0];
      if (currentTab) {
        let url = currentTab.url;
        console.log(url);
        // send url to server
        let urlToSend = "http://localhost:8080/site";
        fetch(urlToSend, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ url: url })
        })
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.text();
        })
        .then(data => {
          console.log(data);
          // コマンド成功の通知を表示
          chrome.notifications.create({
            type: 'basic',
            iconUrl: 'icon48.png',
            title: 'URL Removed',
            message: `The URL was successfully removed: ${url}`
          });
        })
        .catch(error => {
          console.log('There was a problem with the fetch operation:', error.message);
          chrome.notifications.create({
            type: 'basic',
            iconUrl: 'icon48.png',
            title: 'Error',
            message: `Failed to remove the URL: ${error.message}`
          });
        });
      }
    });
  }
});
