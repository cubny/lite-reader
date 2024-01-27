var stage={
    init: function(){
      $('body').layout({
        north:{
          resizable:false,
          closable:false,
          size:"100%",
          spacing_open:0
        },
        west:{
          //resizable:false,
          size:250,
          closable:false,
        }

      });
      $('body > .ui-layout-center').layout({
        north:{
          resizable:false,
          closable:false,
          size:"100%",
          spacing_open:0
        },
        });
  }
}

