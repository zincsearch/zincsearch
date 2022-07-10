export function htmlSpecialChars(str: string) {
  if (!str || str.length == 0) {
    return "";
  }

  let s = "";
  for (let i = 0; i < str.length; i++) {
    switch (str.substring(i, i + 1)) {
      case "<":
        s += "&lt;";
        break;
      case ">":
        s += "&gt;";
        break;
      case "&":
        s += "&amp;";
        break;
      case " ":
        if (str.substring(i + 1, i + 1 + 1) == " ") {
          s += " &nbsp;";
          i++;
        } else {
          s += " ";
        }
        break;
      case '"':
        s += "&quot;";
        break;
      case "'":
        s += "&#39;";
        break;
      case "\n":
        s += "<br>";
        break;
      default:
        s += str.substring(i, i + 1);
        break;
    }
  }
  return s;
}
