var items = {
    
    $elem:$("#items"),
    current:null,
    current_id:null,
    init:function(){
      items.current=null;
      items.$elem.find('> li .title').click(function(){

          var $this=$(this).parent();
          var id=$this.attr('id');

          items.read(id);

          $this.find('.desc').toggle();
          $(this).find('.item-link').toggle();
          $this.toggleClass('selected');

          if(items.current && items.current_id && items.current_id != id){
              items.current.find('.item-link').hide();
              items.current.find('.desc').hide();
              items.current.removeClass('selected');
          }
          items.current=$this;
          items.current_id=id;
          $('.ui-layout-center').scrollTo($this,100);
          $this.find("img.lazy.not-loaded").each(function(i,img){
            $(img).attr('src',$(img).attr('data-original'));
            $(img).removeClass('not-loaded');
          });
      });
    },
    render:function(data){
       items.$elem.html("");
       var unread=0;
       r = new RegExp(/[\u0600-\u06FF]/);
       $.each(data, function(i,item){
         var item_template = $("#item-template").html();

         var $li = $(item_template.format(item.id,item.title,item.desc,item.link,parseInt(item.starred,10)>0?"icon-star":"icon-star-empty"));
         if(r.test(item.title)){
            $li.addClass("rtl");
         }
         $li.find('.desc').hide();
         if(item.is_new=="1"){
           $li.addClass("new");
           unread++;
         }
         $li.find('.item-star').click(function(e){
             var id =$(this).parents('li').attr('id');
             items.starred(id);
            e.stopPropagation();
         })
         $li.data('starred',item.starred);
         items.$elem.append($li);
       });
       items.init();
       return unread;
    },
    starred:function(id){
      var $item=items.$elem.find("li[id='"+id+"']");
      var init_starred = $item.data("starred");
      var action = $item.data("starred")>0?"unstar":"star";
      $item.data("starred",action=="star");
      $item.find(".item-star > i").toggleClass("icon-star").toggleClass("icon-star-empty");
      if(action == "star"){
          feeds.incCount("#starred");
      }else{
          feeds.decCount("#starred");
      }
      $.ajax({
          url: 'agg/'+action+'/'+id,
          dataType: 'json',
          timeout: 2000, //3 second timeout
          error:function(){
              $item.data("starred",init_starred);
              $item.find(".item-star > i").toggleClass("icon-star").toggleClass("icon-star-empty");
              if(action == "star"){
                  feeds.decCount("#starred");
              }else{
                  feeds.incCount("#starred");
              }
          }
      });
    },
    read:function(id){
      var $item=items.$elem.find("li[id='"+id+"']");
      if($item.hasClass("new")){
        $.getJSON('agg/make_read/'+id,function(data){
          if(data){
            var count=feeds.getCurrentCount();
            count.text(parseInt(count.text(),10)-1);
            $item.removeClass("new");
          }
        });
      }
    }
};
