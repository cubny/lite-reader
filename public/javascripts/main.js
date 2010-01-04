	$(document).ready(function () {
     $('body').layout({
       north:{
        resizable:false,
        closable:false,
        size:"auto",
        spacing_open:0
       },
       center:{
         size:250,
         spacing_closed:0,
       }
   });
     $('body > .ui-layout-center').layout({
        south:{
          size:250,
        }
     });
     $(document).click(hideAddFeed);
     $('#addfeed > *').click(showAddFeed);

     var currentFeed=null;
     $('#feeds li').click(function(){
       getItems($(this).attr('id'));
       if(currentFeed)
         currentFeed.removeClass('selected');
       $(this).addClass('selected');
       currentFeed=$(this);
      });

     var currentItem=null;
     

     
   });
   function initItemsBehavior(){
     currentItem=null;
      $('#items li').click(function(){
          getDesc($(this).attr('id'));
          if(currentItem)
            currentItem.removeClass('selected');
          $(this).addClass('selected');
          currentItem=$(this);
          });
   }
   function getDesc(id){
    $.getJSON('?index/getdesc/'+id,function(data){
      if(data.desc==null){
        $("#south").html($('<a href="'+data.link+'"/>').html(data.title));
      }else{
        $("#south").html(data.desc);
      }
    });
   }

   function getItems(id){
    $.getJSON('?index/getitems/'+id,function(data){
       $("#items").html("");
       $.each(data, function(i,item){
          $('<li/>').attr("id",item.id).html(item.title).appendTo("#items");
      })
       initItemsBehavior();
    })
   }
   function urlencode(s) {
     s = encodeURIComponent(s);
     return s.replace(/~/g,'%7E').replace(/%20/g,'+');
   }

   function addFeed(url){
    $.getJSON('?index/add/'+urlencode(url),function(data){
       var container=$(".ul-layout-west ul");
       container.html("");
       $.each(data, function(i,feed){
          $('<li/>').html($('<a href="javascript:getItems('+feed.id+')" />').html(feed.title)).appendTo(container);
      })
    })
   }
 
   function showAddFeed(e){
     e.stopPropagation();
     var aButton=$("#addfeed a");
     var aInput=$("#urlToAdd");
     aInput.removeClass('ui-state-error');
     aButton.parent().animate({width:"230px"},500,null,function(e){
            aInput.show();
            aButton.unbind('click');
            aButton.click(function(e){
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
     });
   }
   function hideAddFeed(){
     var aButton=$("#addfeed a");
     var aInput=$("#urlToAdd");
     aInput.hide();
     aButton.parent().animate({width:"60px"},300,null,function(){
       aButton.unbind('click');
       $('#addfeed > *').click(showAddFeed);
     });
   }
   function checkLength(o,n,min,max) {
     if ( o.val().length > max || o.val().length < min ) {
       o.addClass('ui-state-error');
       alert("Length of " + n + " must be between "+min+" and "+max+".");
       return false;
     } else {
       return true;
     }
   }

   function checkRegexp(o,regexp,n) {
     if ( !( regexp.test( o.val() ) ) ) {
       o.addClass('ui-state-error');
       alert(n);
       return false;
     } else {
       return true;
     }
   }



