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
          })
          .catch(error => {
            console.log('There was a problem with the fetch operation:', error.message);
          });
      }
    });
  }
  // removelink
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
          })
          .catch(error => {
            console.log('There was a problem with the fetch operation:', error.message);
          });
      }
    });
  }
});
