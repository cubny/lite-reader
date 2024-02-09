var feeds = {
  container: $("#feeds ul"),
  current: null,
  init: function () {
    $.getJSON("/feeds", function (data) {
      // add all the feeds to the page, bind click events
      $.each(data, function (id, feed) {
        feeds.add(feed);
      });
    });
    feeds.container.find("li.feed").live("click", function () {
      $(".ui-layout-center").scrollTop(0);
      var $this = $(this);
      var id = this.id;
      $(this).removeClass("new");
      $.cookie("feed", id);
      feeds.load(id);
      if (feeds.currentFeed) {
        feeds.currentFeed.removeClass("selected");
      }
      $this.addClass("selected");
      $("#title").html($this.find(".feedtitle").html());
      $("#actions .action").attr("id", id).show();
      if(id === "unread" || id === "starred"){
        $(".remove").hide()
      } else {
        $(".remove").show()
      }
      feeds.currentFeed = $this;
    });
    setTimeout(function () {
      var lastFeed = $.cookie("feed");
      if (lastFeed !== undefined) {
        feeds.container.find("li#" + lastFeed).click();
      }
    }, 1000);
    this.getUnreadItemsCount();
    this.getStarredItemsCount();
  },
  load: function (id) {
    const url =
      id === "unread" || id === "starred" ? `items/${id}` : `feeds/${id}/items`;
    $.getJSON(url, function (data) {
      items.render(data);
    });
  },
  add: function (options) {
    var defaults = {
      id: 0,
      title: "Untitled Feed",
      url: "http://samplefeedurl.xml",
      unread: 20,
    };
    var opts = $.extend(defaults, options);
    var $item = $("<li/>")
      .addClass("new")
      .addClass("feed")
      .attr("id", opts.id)
      .html(
        $("<img/>").attr("src", "http://api.byi.pw/favicon?url=" + opts.link),
      )
      .append($("<div/>").addClass("feedtitle").html(opts.title))
      .append(
        $("<div/>")
          .addClass("count")
          .html("<span>" + opts.unread_count + "</span>"),
      );
    feeds.container.append($item);
    feeds.blink(opts);
    $("#feeds-actions").show();
  },
  del: function (id) {
    $.ajax({
        url: `feeds/${id}`,
        type: "DELETE",
        success: function (data) {
          feeds.container.find("li[id='" + id + "']").remove();
          $("#items").html("");
          $("#actions .action").hide();
          $("#feedbar #title").html("");
        },
    });
  },
  markread: function (id) {
    $.ajax({
        url: `feeds/${id}/read`,
        type: "POST",
        success: function (data) {
            feeds.setCurrentCount(0);
            items.$elem.find("> li").removeClass("new");
        },
    });
  },
  markunread: function (id) {
    $.ajax({
        url: `feeds/${id}/unread`,
        type: "POST",
        success: function (data) {
            feeds.setCurrentCount(items.$elem.find("> li").length);
            items.$elem.find("li").addClass("new");
        },
    });
  },
  update: function (id) {
    $.ajax({
      url: `feeds/${id}/fetch`,
      type: "PUT",
      content: "application/json",
      success: function (encryptedData) {
        const data = JSON.parse(encryptedData);
        const unread = items.render(data);
        feeds.setCurrentCount(unread);
      },
    });
  },
  getCurrentCount: function () {
    return feeds.container.find("li.selected span");
  },
  getCount: function (selector) {
    var counter = feeds.container.find("li" + selector + " .count span");
    return counter.length > 0 ? parseInt(counter.text(), 10) : 0;
  },
  incCount: function (selector) {
    feeds.setCount(selector, feeds.getCount(selector) + 1);
  },
  decCount: function (selector) {
    feeds.setCount(selector, feeds.getCount(selector) - 1);
  },
  getUnreadItemsCount: function () {
    $.getJSON("items/unread/count", function (data) {
        feeds.setCount("#unread", data.count);
    });
  },
  getStarredItemsCount: function () {
    $.getJSON("items/starred/count", function (data) {
        feeds.setCount("#starred", data.count);
    });
  },
  setCount: function (selector, count) {
    if (count > 0) {
      feeds.container
        .find("li" + selector + " .count")
        .html("<span>" + count + "</span>");
    } else {
      feeds.container.find("li" + selector + " .count").html("");
    }
  },
  setCurrentCount: function (count) {
    feeds.setCount(count, ".selected");
  },
  blink: function (feed) {
    feeds.container
      .find("li[id='" + feed.id + "']")
      .effect("pulsate", { times: 3 }, 200);
    feeds.container.find("li[id='" + feed.id + "']").click();
  },
};
