var feeds={
    container:$("#feeds ul"),
    current:null,          
    init:function(){
      feeds.container.find("li").live('click',function(){
        var $this=$(this);
        var id=this.id;
        $.cookie("feed",id);
        feeds.load(id);
        if(feeds.currentFeed)feeds.currentFeed.removeClass('selected');
        $this.addClass('selected');
        $("#title").html($this.find(".feedtitle").html());
        $(".remove").show();
        $(".remove").attr('id',id);
        $(".update").show();
        $(".update").attr('id',id);
        feeds.currentFeed=$this;
      });
      if(lastFeed=$.cookie("feed")){
        feeds.container.find("li#"+lastFeed).click();
      }
    },
    load:function(id){
       loadr.show();
       $.getJSON('agg/getitems/'+id,function(data){
           items.render(data);
           loadr.hide();
       })
    },
    add:function(options){
       var defaults={
            id:0,
            title:"Untitled Feed",
            url:"http://samplefeedurl.xml",
            unread:20,
       };
       var opts = $.extend(defaults, options);
       var $item=$("<li/>").addClass("new").attr("id",opts.id)
                           .html(
                                $("<img/>").attr("src","http://cdn.netvibes.com/proxy/favIcon.php?url="+urlencode(opts.url))
                               )
                           .append(
                               $("<div/>").addClass("feedtitle").html(opts.title)
                               )
                           .append(
                               $("<div/>").addClass("count").html("(<span>"+opts.unread+"</span>)")
                               );
       feeds.container.append($item);
       $item.animate({backgroundColor:"white"},2000)
    },
    del:function(id){
        $.get('agg/del/'+id,function(){
            feeds.container.find("li[id='"+id+"']").remove()
            $("#items").html("");
            $("#feedbar .remove").hide();
            $("#feedbar .update").hide();
            $("#feedbar #title").html("");
        });
    },
    update:function(id){
        loadr.show();
        $.getJSON('agg/update/'+id,function(data){
            unread=items.render(data);
            var count=feeds.getCurrentCount();
            count.text(unread);
            loadr.hide();
      });
    },
    getCurrentCount:function(){
          return feeds.container.find("li.selected span");
    },
    blink:function(id){
        feeds.container.find("li[id='"+id+"']").effect("pulsate",{times:3},200)
    },
};
