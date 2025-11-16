// Feed Management - Refactored to work with HTMX
var feeds = {
  container: null,
  current: null,
  currentFeed: null,

  init: function () {
    // Use vanilla JS instead of jQuery for container reference
    feeds.container = document.querySelector("#feeds-list");
    
    // Set up feed click handlers (using event delegation)
    if (feeds.container) {
      feeds.container.addEventListener("click", function(e) {
        // Find the closest li.feed element
        var feedItem = e.target.closest("li.feed");
        if (feedItem && !feedItem.classList.contains("feed-empty")) {
          feeds.handleFeedClick(feedItem);
        }
      });
    }

    // Load initial feed list from server via HTMX
    // The server will return HTML fragments for HTMX requests
    if (typeof htmx !== 'undefined') {
      htmx.ajax('GET', '/feeds', {
        target: '#feeds-list',
        swap: 'innerHTML'
      });
    }

    // Restore last selected feed after a delay
    setTimeout(function () {
      var lastFeed = getCookie("feed");
      if (lastFeed !== undefined && lastFeed !== null) {
        var feedElement = document.getElementById(lastFeed);
        if (feedElement) {
          feedElement.click();
        }
      }
    }, 1000);

    // Initialize unread and starred counts
    feeds.getUnreadItemsCount();
    feeds.getStarredItemsCount();
  },

  handleFeedClick: function(feedItem) {
    var feedId = feedItem.id;
    
    // Remove 'new' class
    feedItem.classList.remove("new");
    
    // Save to cookie
    setCookie("feed", feedId);
    
    // Load feed items
    feeds.load(feedId);
    
    // Update selected state
    if (feeds.currentFeed) {
      feeds.currentFeed.classList.remove("selected");
    }
    feedItem.classList.add("selected");
    feeds.currentFeed = feedItem;
    
    // Update title
    var titleElement = feedItem.querySelector(".feedtitle");
    if (titleElement) {
      var titleBar = document.getElementById("title");
      if (titleBar) {
        titleBar.innerHTML = titleElement.innerHTML;
      }
    }
    
    // Show/hide actions based on feed type
    var actions = document.querySelector("#actions .action");
    if (actions) {
      var allActions = document.querySelectorAll("#actions .action");
      allActions.forEach(function(action) {
        action.setAttribute("id", feedId);
        action.style.display = "block";
      });
      
      var removeBtn = document.querySelector(".remove");
      if (feedId === "unread" || feedId === "starred") {
        if (removeBtn) removeBtn.style.display = "none";
      } else {
        if (removeBtn) removeBtn.style.display = "block";
      }
    }
  },

  load: function (id) {
    const url =
      id === "unread" || id === "starred" ? `items/${id}` : `feeds/${id}/items`;
    
    // Use fetch API instead of jQuery
    fetch(url, {
      headers: {
        'Authorization': 'Bearer ' + getAuthToken()
      }
    })
    .then(response => response.json())
    .then(data => {
      items.render(data);
    })
    .catch(error => {
      console.error('Error loading feed items:', error);
    });
  },

  getCurrentCount: function () {
    var selected = feeds.container.querySelector("li.selected .count span");
    return selected;
  },

  getCount: function (selector) {
    var counter = feeds.container.querySelector("li" + selector + " .count span");
    return counter ? parseInt(counter.textContent, 10) : 0;
  },

  incCount: function (selector) {
    feeds.setCount(selector, feeds.getCount(selector) + 1);
  },

  decCount: function (selector) {
    feeds.setCount(selector, feeds.getCount(selector) - 1);
  },

  getUnreadItemsCount: function () {
    fetch("items/unread/count", {
      headers: {
        'Authorization': 'Bearer ' + getAuthToken()
      }
    })
    .then(response => response.json())
    .then(data => {
      feeds.setCount("#unread", data.count);
    })
    .catch(error => {
      console.error('Error getting unread count:', error);
    });
  },

  getStarredItemsCount: function () {
    fetch("items/starred/count", {
      headers: {
        'Authorization': 'Bearer ' + getAuthToken()
      }
    })
    .then(response => response.json())
    .then(data => {
      feeds.setCount("#starred", data.count);
    })
    .catch(error => {
      console.error('Error getting starred count:', error);
    });
  },

  setCount: function (selector, count) {
    var countElement = feeds.container.querySelector("li" + selector + " .count");
    if (countElement) {
      if (count > 0) {
        countElement.innerHTML = "<span>" + count + "</span>";
      } else {
        countElement.innerHTML = "";
      }
    }
  },

  setCurrentCount: function (count) {
    if (feeds.currentFeed) {
      var countElement = feeds.currentFeed.querySelector(".count");
      if (countElement) {
        if (count > 0) {
          countElement.innerHTML = "<span>" + count + "</span>";
        } else {
          countElement.innerHTML = "";
        }
      }
    }
  },

  // Delete feed using HTMX
  del: function (id) {
    if (confirm("Are you sure you want to remove this feed?")) {
      if (typeof htmx !== 'undefined') {
        htmx.ajax('DELETE', `feeds/${id}`, {
          target: '#feeds-list',
          swap: 'innerHTML'
        }).then(function() {
          // Clear items and actions
          var itemsContainer = document.getElementById("items");
          if (itemsContainer) itemsContainer.innerHTML = "";
          
          var actionsContainer = document.querySelector("#actions .action");
          if (actionsContainer) {
            document.querySelectorAll("#actions .action").forEach(function(el) {
              el.style.display = "none";
            });
          }
          
          var titleBar = document.getElementById("title");
          if (titleBar) titleBar.innerHTML = "";
        });
      }
    }
  },

  // Mark all items in feed as read
  markread: function (id) {
    fetch(`feeds/${id}/read`, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + getAuthToken()
      }
    })
    .then(() => {
      feeds.setCurrentCount(0);
      var itemElements = document.querySelectorAll("#items > li");
      itemElements.forEach(function(item) {
        item.classList.remove("new");
      });
    })
    .catch(error => {
      console.error('Error marking feed as read:', error);
    });
  },

  // Mark all items in feed as unread
  markunread: function (id) {
    fetch(`feeds/${id}/unread`, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + getAuthToken()
      }
    })
    .then(() => {
      var itemElements = document.querySelectorAll("#items li");
      var count = itemElements.length;
      feeds.setCurrentCount(count);
      itemElements.forEach(function(item) {
        item.classList.add("new");
      });
    })
    .catch(error => {
      console.error('Error marking feed as unread:', error);
    });
  },

  // Update/refresh feed items
  update: function (id) {
    fetch(`feeds/${id}/fetch`, {
      method: 'PUT',
      headers: {
        'Authorization': 'Bearer ' + getAuthToken()
      }
    })
    .then(response => response.json())
    .then(data => {
      const unread = items.render(data);
      feeds.setCurrentCount(unread);
    })
    .catch(error => {
      console.error('Error updating feed:', error);
    });
  }
};

// Helper functions for cookies
function getCookie(name) {
  var nameEQ = name + "=";
  var ca = document.cookie.split(';');
  for(var i=0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

function setCookie(name, value, days) {
  var expires = "";
  if (days) {
    var date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

