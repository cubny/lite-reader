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
      setTimeout(function(){
        var lastFeed = $.cookie("feed");
        if(lastFeed !== undefined) {
          feeds.container.find("li#"+lastFeed).click();
        }
      },1000);
    },
    load:function(id){
       $.getJSON('agg/getitems/'+id,function(data){
           items.render(data);
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
         .html($("<img/>").attr("src","http://api.byi.pw/favicon?url="+opts.url))
         .append($("<div/>").addClass("feedtitle").html(opts.title))
         .append($("<div/>").addClass("count").html("<span>"+opts.unread+"</span>"));
       feeds.container.append($item);
       feeds.blink(opts);
       $("#feeds-actions").show();
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
        $.getJSON('agg/mark_read_all/'+id,function(data){
          feeds.setCurrentCount(0);
          items.$elem.find("> li").removeClass("new");
        });
    },
    markunread: function(id){
        $.getJSON('agg/mark_unread_all/'+id,function(data){
          feeds.setCurrentCount(items.$elem.find('> li').length);
          items.$elem.find("li").addClass("new");
        });
    },
    update:function(id){
        $.getJSON('agg/update/'+id,function(data){
            unread=items.render(data);
            var count=feeds.setCurrentCount(unread);
      });
    },
    update_all:function(){
        $.getJSON('agg/update_all',function(data){
            $.each(data,function(id,feed_items){
                var unreads = 0;
                $.each(feed_items,function(item_id,item){
                    if(item.is_new == 1){
                        unreads++;
                    }
                });
                feeds.setCount("#"+id,unreads);
                if(id == feeds.currentFeed.attr('id')){
                    items.render(feed_items);
                }
            });
        });
    },
    getCurrentCount:function(){
          return feeds.container.find("li.selected span");
    },
    getCount:function(selector){
        var counter = feeds.container.find("li"+selector+" .count span");
        return counter.length>0?parseInt(counter.text(),10):0;
    },
    incCount:function(selector){
        feeds.setCount(selector,feeds.getCount(selector)+1);
    },
    decCount:function(selector){
        feeds.setCount(selector,feeds.getCount(selector)-1);
    },
    setCount: function(selector,count){
      if(count > 0){
        feeds.container.find("li"+selector+" .count").html("<span>"+count+"</span>");
      }else{
        feeds.container.find("li"+selector+" .count").html("");
      }
    },
    setCurrentCount: function(count){
        feeds.setCount(count,".selected");
    },
    blink:function(feed){
        feeds.container.find("li[id='"+feed.id+"']").effect("pulsate",{times:3},200);
        feeds.container.find("li[id='"+feed.id+"']").click();
    },
};
