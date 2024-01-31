const items = {
    $elem:$("#items"),
    current:null,
    current_id:null,
    init:function(){
        items.current=null;
        items.$elem.find('> li .title').click(items.click);
    },
    click: function() {
        const $this=$(this).parent();
        const id=$this.attr('id');

        items.read(id);

        $this.find('.desc').toggle();
        $(this).find('.item-link').toggle();
        $this.toggleClass('selected');

        if(items.current && items.current_id && items.current_id !== id){
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
    },
    render:function(data){
        if(data.length===0){
            items.$elem.html("<li class='empty'>No items found</li>");
            return;
        }
        items.$elem.html("");
        let unread=0;
        const r = new RegExp(/[\u0600-\u06FF]/);
        data.forEach(function(item){
            const item_template = $("#item-template").html();

            const $li = $(item_template.format(
                item.id,
                item.title,
                item.desc,
                item.link,
                item.is_new ? "icon-circle":"icon-circle-blank",
                item.starred?"icon-star":"icon-star-empty",
                moment(item.timestamp || new Date(),"YYYY-MM-DD HH:ii:SS").calendar().format('LL')
            ));
            if(r.test(item.title)){
                $li.addClass("rtl");
            }
            $li.find('.desc').hide();
            if(item.is_new){
                $li.addClass("new");
                unread++;
            }

            $li.find('.item-star').click(function(e){
                const $li = $(this).parents('li');
                const id = $li.attr('id');
                if($li.data("starred")) {
                    items.unstar(id);
                }else{
                    items.star(id);
                }
                e.stopPropagation();
            });

            $li.find('.item-read').click(function(e){
                const $li = $(this).parents('li');
                const id = $li.attr('id');
                if($li.hasClass("new")){
                    items.read(id);
                }else{
                    items.unread(id);
                }
                e.stopPropagation();
            });

            $li.data('starred',item.starred);
            items.$elem.append($li);
        });
        items.init();
        return unread;
    },

    star: function(id) {
        const $item = items.$elem.find("li[id='" + id + "']");

        items.update(id, $item.data("is_new"), true, function(){
            $item.data("starred", true);
            $item.find(".item-star > i").removeClass("icon-star-empty").addClass("icon-star");
            feeds.incCount("#starred");
        });
    },

    unstar: function(id) {
        const $item = items.$elem.find("li[id='" + id + "']");

        items.update(id, $item.data("is_new"), false, function(){
            $item.data("starred", false);
            $item.find(".item-star > i").removeClass("icon-star").addClass("icon-star-empty");
            feeds.decCount("#starred");
        });
    },
    unread:function(id){
        const $item=items.$elem.find("li[id='"+id+"']");
        if($item.hasClass("new")) return
        items.update(id, true, $item.data("starred"), function(){
            feeds.incCount(".selected");
            $item.addClass("new");
            $item.find(".item-read i").removeClass("icon-circle-blank").addClass("icon-circle");
        });
    },
    read:function(id){
        const $item=items.$elem.find("li[id='"+id+"']");
        if(!$item.hasClass("new")){
            return;
        }
        items.update(id, false, $item.data("starred"), function(){
            feeds.decCount(".selected");
            $item.removeClass("new");
            $item.find(".item-read i").removeClass("icon-circle").addClass("icon-circle-blank");
        });
    },
    update:function(id, is_new, starred, success, error){
        const $item=items.$elem.find("li[id='"+id+"']");
        $.ajax({
            url: "items/" + id,
            type: "PUT",
            data: JSON.stringify({
                is_new: is_new,
                starred: starred,
            }),
            contentType: "application/json",
            success: success,
            error: error,
        })
    },
};
