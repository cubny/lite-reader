$(document).ready(function () {
    stage.init();
    loadr.init();
    
    // Feed form toggle handlers
    var showFeedFormBtn = document.getElementById('show-feed-form');
    var addFeedForm = document.getElementById('add-feed-form');
    var urlInput = document.getElementById('urlToAdd');
    
    if (showFeedFormBtn && addFeedForm && urlInput) {
      showFeedFormBtn.addEventListener('click', function(e) {
        e.preventDefault();
        e.stopPropagation();
        
        if (addFeedForm.style.display === 'none') {
          // Show the form
          addFeedForm.style.display = 'inline-block';
          urlInput.style.display = 'inline-block';
          urlInput.focus();
          showFeedFormBtn.classList.remove('btn-purple');
          showFeedFormBtn.classList.add('btn-green');
          showFeedFormBtn.querySelector('span').textContent = '';
        } else {
          // Hide the form
          addFeedForm.style.display = 'none';
          urlInput.value = '';
          showFeedFormBtn.classList.remove('btn-green');
          showFeedFormBtn.classList.add('btn-purple');
          showFeedFormBtn.querySelector('span').textContent = 'Feed';
        }
      });
      
      // Handle Enter key in input
      urlInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
          e.preventDefault();
          addFeedForm.querySelector('button[type="submit"]').click();
        }
      });
      
      // Reset form after successful submission
      document.body.addEventListener('htmx:afterSwap', function(evt) {
        if (evt.detail.target.id === 'feeds-list') {
          addFeedForm.style.display = 'none';
          urlInput.value = '';
          showFeedFormBtn.classList.remove('btn-green');
          showFeedFormBtn.classList.add('btn-purple');
          showFeedFormBtn.querySelector('span').textContent = 'Feed';
          
          // Re-initialize counts after feed list update
          setTimeout(function() {
            feeds.getUnreadItemsCount();
            feeds.getStarredItemsCount();
          }, 500);
        }
      });
      
      // Hide form when clicking outside
      document.addEventListener('click', function(e) {
        if (!e.target.closest('#addfeed') && addFeedForm.style.display !== 'none') {
          addFeedForm.style.display = 'none';
          urlInput.value = '';
          showFeedFormBtn.classList.remove('btn-green');
          showFeedFormBtn.classList.add('btn-purple');
          showFeedFormBtn.querySelector('span').textContent = 'Feed';
        }
      });
    }
    
    // Action button handlers (keep using jQuery for compatibility with existing code)
    $('.update').click(function(){
      feeds.update(this.id);
    });
    $('#mark-read-all').click(function(){
      feeds.markread(this.id);
    });
    $('#mark-unread-all').click(function(){
      feeds.markunread(this.id);
    });
    $('.remove').click(function(){
      feeds.del(this.id);
    });
    $('.logout').click(function(){
      logout()
    });
    
    feeds.init();

    if(feeds.container && feeds.container.querySelectorAll('li').length < 2){
      var feedActions = document.getElementById('feeds-actions');
      if (feedActions) {
        feedActions.style.display = 'none';
      }
    }
    
    $(document).bind('keydown',function(e){
      var code = (e.keyCode ? e.keyCode : e.which);
      if(code == 32) {
        if(items.current !== null){
          var h = items.current.height();
          var top = items.current.offset().top;
          console.log(h+top);
          if(top+h<$('.ui-layout-center').height()){
            items.current.next().find('.title').click();
          }
        }
      }
    });
});

function logout() {
    $.ajax({
        url: 'logout',
        type: "POST",
        complete: function() {
          clearAuthToken();
          window.location.href = '/login.html';
        }
    })
}





