var items = {
    
    $elem:$("#items"),
    current:null,
    init:function(){
      items.current=null;
      items.$elem.find('li').click(function(){
          var id=this.id;
          var $this=$(this);
          items.load(id);
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
           var $li=$('<li/>').attr("id",item.id).html(item.title);
           if(item.is_new=="1"){
             $li.addClass("new"); 
             unread++;
           }
           $li.appendTo(items.$elem);
       })
       items.init();
       return unread;
    },
    load:function(id){
       $.getJSON('agg/getdesc/'+id,function(data){
           if(data.desc==null){
             $("#south").html($('<a href="'+data.link+'"/>').html(data.title));
           }else{
             $("#south").html(data.desc);
           }
           items.read(id);
       });
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
