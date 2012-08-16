function checkLength(o,n,min,max) {
  if ( o.val().length > max || o.val().length < min ) {
    o.addClass('ui-state-error');
    alert("Length of " + n + " must be between "+min+" and "+max+".");
    return false;
  } else {
    return true;
  }
}
function checkRegexp(o,regexp,n) {
  if ( !( regexp.test( o.val() ) ) ) {
    o.addClass('ui-state-error');
    alert(n);
    return false;
  } else {
    return true;
  }
}
function urlencode(s) {
  s = encodeURIComponent(s);
  return s.replace(/~/g,'%7E').replace(/%20/g,'+');
}
var loadr = {
  elem:$("#msg"),
  show:function(){
         loadr.elem.animate({top:"15px"},500);
   },
  hide:function(){
         loadr.elem.animate({top:"-20px"},500);
   },
}
