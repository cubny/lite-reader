$(document).ready(function () {
    stage.init();
    $(document).click(hideAddFeed);
    $('#addfeed > *').click(showAddFeed);
    $('.update').click(function(){
      feeds.update(this.id);
      });
    $('.remove').click(function(){
      feeds.del(this.id);
      });
    feeds.init();
    
});

function addFeed(url){
  var aInput=$("#urlToAdd");
  var af=$('#addfeed .add');
  var currImg='url(public/images/add.png)';
  af.css('background-image','url(public/images/loading.gif)');
  af.text('Adding Feed...');
  $.ajax({
      url:'agg/add',
      type:"POST",
      data:"url="+url,
      dataType:"json",
      success:function(data){
          af.text('Feed');
        if(!data.error){
          $.each(data, function(i,feed){
            feeds.add(feed);
          });
        }else{
            feeds.blink(data.feed);
        }
        af.css('background-image',currImg);
        aInput.val("");
        hideAddFeed();
      }   
  })
}

function showAddFeed(e){
  e.stopPropagation();
  var aButton=$("#addfeed a");
  var aInput=$("#urlToAdd");
  aInput.removeClass('ui-state-error');
  //aButton.parent().animate({width:"230px"},500,null,function(e){
      aInput.show();
      aButton.unbind('click');
      aButton.click(function(e){
        aButton.unbind('click');
        e.stopPropagation();
        var bValid = true;
        aInput.removeClass('ui-state-error');
        bValid = bValid && checkLength(aInput,"password",10,255);
        bValid = bValid && checkRegexp(aInput,/(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?/i,"Please Enter a Valid Url.");
        if(bValid){
          addFeed(aInput.val());
          aInput.val("");
          hideAddFeed();
        }
      });
  //});
}

function hideAddFeed(){
  var aButton=$("#addfeed a");
  var aInput=$("#urlToAdd");
  aInput.hide();
  $('#addfeed > *').click(showAddFeed);
  //aButton.parent().animate({width:"60px"},300,null,function(){
      //$('#addfeed > *').click(showAddFeed);
  //});
}





