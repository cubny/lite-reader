var items = {
    
    $elem:$("#items"),
    current:null,
    init:function(){
      items.current=null;
      items.$elem.find('li').click(function(){
          var id=this.id;
          var $this=$(this);
          $this.find('.desc').toggle();
          items.read(id);
          if(items.current)
            items.current.removeClass('selected');
          $this.addClass('selected');
          items.current=$this;
      });

    },
    render:function(data){
       items.$elem.html("");
       var unread=0;
       $.each(data, function(i,item){
         var item_template = $("#item-template").html();
         var $li = $(item_template.format(item.id,item.title,item.desc))
         $li.find('.desc').hide();
         if(item.is_new=="1"){
           $li.addClass("new");
           unread++;
         }
         items.$elem.append($li)
       })
       items.init();
       return unread;
    },
    read:function(id){
      var $item=items.$elem.find("li[id='"+id+"']");
      if($item.hasClass("new")){
        var count=feeds.getCurrentCount();
        count.text(parseInt(count.text())-1);
        $item.removeClass("new");
      }
    },
}
