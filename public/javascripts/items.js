var items = {
    
    $elem:$("#items"),
    current:null,
    current_id:null,
    init:function(){
      items.current=null;
      items.$elem.find('li .title').click(function(){

          var $this=$(this).parent();
          var id=$this.attr('id');

          items.read(id);

          $this.find('.desc').toggle();
          $this.toggleClass('selected');
          if(items.current && items.current_id && items.current_id != id){
            items.current.find('.desc').hide();
            items.current.removeClass('selected')

          }
          items.current=$this;
          items.current_id=id;

      });

    },
    render:function(data){
       items.$elem.html("");
       var unread=0;
       r = new RegExp(/[\u0600-\u06FF]/);
       $.each(data, function(i,item){
         var item_template = $("#item-template").html();

         var $li = $(item_template.format(item.id,item.title,item.desc))
         if(r.test(item.title)){
            $li.css({
              "direction":"rtl"
            })
         }
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
        $.getJSON('agg/make_read/'+id,function(data){
          if(data){
            var count=feeds.getCurrentCount();
            count.text(parseInt(count.text())-1);
            $item.removeClass("new");
          }
        })
      }
    },
}
