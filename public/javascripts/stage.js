var stage={
    init: function(){
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
        },
        west:{
          size:250,
        }

      });
      $('body > .ui-layout-center').layout({
        north:{
          resizable:false,
          closable:false,
          size:"auto",
          spacing_open:0
        },
        south:{
          size:250,
        }
        });
  }
}

