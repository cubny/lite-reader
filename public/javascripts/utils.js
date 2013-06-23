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

function get_hostname(url) {
        var m = ((url||'')+'').match(/^http:\/\/([^/]+)/);
        return m ? m[1] : null;
}

var loadr = {
  elem:$("#msg"),
  show:function(){
         loadr.elem.animate({top:"60px"},500);
   },
  hide:function(){
         loadr.elem.animate({top:"0px"},500);
   },
}

String.prototype.format = function() {
    var args = arguments;
    return this.replace(/{(\d+)}/g, function(match, number) {
      return typeof args[number] != 'undefined'? args[number]: match;
    });
};
