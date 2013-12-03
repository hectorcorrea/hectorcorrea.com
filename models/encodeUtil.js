var HtmlEncoder = require('node-html-encoder').Encoder;
var encoder = new HtmlEncoder('entity');

exports.encodeText = function(text) {
  if (text === undefined) {
    return '';
  }
  text = encoder.htmlEncode(text);
  text = text.replace(/&#13;&#10;/g, '<br/>');
  text = text.replace(/&#13;/g,'<br/>');
  text = text.replace(/&#10;/g,'<br/>');
  return text;
};


exports.decodeText = function(text) {
  if (text === undefined) {
    return '';
  }  
  text = encoder.htmlDecode(text);
  text = text.replace(/<br\/>/g, '\r\n')
  return text;
};


exports.urlSafe = function(text) {
  var i, c;
  var url = "";

  text = text.trim().toLowerCase();
  text = text .replace('c#', 'c-sharp')

  for(i=0; i<text.length; i++) {
    c = text[i];
    if(c >= 'a' && c <= 'z')
      url += c;
    else if(c >= '0' && c <= '9')
      url += c;
    else
      url += '-';
  }

  while(url.indexOf('--') > -1) {
    url = url.replace('--', '-');
  }

  if(url.charAt(url.length-1) === '-') {
    url = url.substr(0, url.length-1);
  }

  return url
};
