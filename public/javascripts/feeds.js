var feeds={
    container:$("#feeds ul"),
    current:null,
    init:function(){
      feeds.container.find("li.feed").live('click',function(){
        $('.ui-layout-center').scrollTop(0);
        var $this=$(this);
        var id=this.id;
        $(this).removeClass('new');
        $.cookie("feed",id);
        feeds.load(id);
        if(feeds.currentFeed) { 
            feeds.currentFeed.removeClass('selected');
        }
        $this.addClass('selected');
        $("#title").html($this.find(".feedtitle").html());
        $("#actions .action").attr('id',id).show();
        feeds.currentFeed=$this;
      });
      var lastFeed = $.cookie("feed");
      if(lastFeed !== undefined) {
        feeds.container.find("li#"+lastFeed).click();
      }
    },
    load:function(id){
       loadr.show();
       $.getJSON('agg/getitems/'+id,function(data){
           items.render(data);
           loadr.hide();
       });
    },
    add:function(options){
       var defaults={
            id:0,
            title:"Untitled Feed",
            url:"http://samplefeedurl.xml",
            unread:20,
       };
       var opts = $.extend(defaults, options);
       var $item=$("<li/>").addClass("new").addClass("feed").attr("id",opts.id)
         .html($("<img/>").attr("src","http://geticon.org/of/"+get_hostname(opts.url)))
         .append($("<div/>").addClass("feedtitle").html(opts.title))
         .append($("<div/>").addClass("count").html("<span>"+opts.unread+"</span>"));
       feeds.container.append($item);
       feeds.blink(opts);
       // $item.animate({backgroundColor:"white"},2000);
    },
    del:function(id){
        $.get('agg/del/'+id,function(){
            feeds.container.find("li[id='"+id+"']").remove();
            $("#items").html("");
            $("#actions .action").hide();
            $("#feedbar #title").html("");
        });
    },
    markread: function(id){
        loadr.show();
        $.getJSON('agg/mark_read_all/'+id,function(data){
          feeds.setCurrentCount(0);
          items.$elem.find("> li").removeClass("new");
          loadr.hide();
        });
    },
    markunread: function(id){
        loadr.show();
        $.getJSON('agg/mark_unread_all/'+id,function(data){
          feeds.setCurrentCount(items.$elem.find('> li').length);
          items.$elem.find("li").addClass("new");
          loadr.hide();
        });
    },
    update:function(id){
        loadr.show();
        $.getJSON('agg/update/'+id,function(data){
            unread=items.render(data);
            var count=feeds.setCurrentCount(unread);
            loadr.hide();
      });
    },
    update_all:function(){
        loadr.show();
        $.getJSON('agg/update_all',function(data){
            $.each(data,function(id,feed_items){
                var unreads = 0;
                $.each(feed_items,function(item_id,item){
                    if(item.is_new == 1){
                        unreads++;
                    }
                });
                if(unreads>0){
                    feeds.container.find("li[id='"+id+"' .count").html("<span>"+unreads+"</span>")
                }else{
                    feeds.container.find("li[id='"+id+"' .count").html("")
                }
                if(id == feeds.currentFeed.attr('id')){
                    items.render(feed_items);
                }
            });
            loadr.hide();
        });    
    },
    getCurrentCount:function(){
          return feeds.container.find("li.selected span");
    },
    setCurrentCount: function(count){
      if(count > 0){
        feeds.container.find("li.selected .count").html("<span>"+count+"</span>");
      }else{
        feeds.container.find("li.selected .count").html("");
      }
    },
    blink:function(feed){
        feeds.container.find("li[id='"+feed.id+"']").effect("pulsate",{times:3},200);
        feeds.container.find("li[id='"+feed.id+"']").click();
    },
};
