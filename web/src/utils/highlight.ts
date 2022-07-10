// eg.1: Gold => ['Gold']
// eg.2: City:Paris => ['Paris']
// eg.3: City:Paris Gold => ['Paris', 'Gold']
// eg.4: City:par* => ['par']
// eg.5: "Paris Gold" => ['Paris Gold']
import { htmlSpecialChars } from "./html";

export function getKeywords(queryString: string) {
  if (!queryString || queryString.trim().length == 0) {
    return [];
  }

  let arr = [];
  // queryString + " " is for special split regular
  // split by space, but ignore double quotation marks
  const groups = (queryString + " ").split(/ s*(?![^"]*"\ )/);
  for (let i = 0; i < groups.length - 1; i++) {
    const group = groups[i];
    if (!group || group.trim().length == 0) {
      continue;
    }
    // group + ":" is for special split regular
    // split by :, but ignore "
    const fieldWordArr = (group + ":").split(/:s*(?![^"]*"\:)/);
    let keyword = group;
    if (fieldWordArr.length > 2) {
      keyword = fieldWordArr[1];
    }
    // delete start and end of * and "
    keyword = keyword.replace(/(^\**)|(\**$)/g, "").replace(/(^"*)|("*$)/g, "");
    if (keyword.trim().length > 0) {
      // make sure key not empty or not space
      arr.push(keyword);
    }
  }
  return arr;
}

export function highlightAndSpecialChars(value: any, keywords: []) {
  if (!value) {
    return value;
  }

  if (typeof value == "string") {
    value = htmlSpecialChars(value);
    for (const idx in keywords) {
      const keyword = htmlSpecialChars(keywords[idx]);
      const highlightText = "<span class='highlight'>" + keyword + "</span>";
      value = value.replaceAll(keyword, highlightText);
    }
  } else if (Array.isArray(value)) {
    for (let i = 0; i < value.length; i++) {
      value[i] = highlightAndSpecialChars(value[i], keywords);
    }
  } else if (typeof value == "object") {
    for (const key in value) {
      value[key] = highlightAndSpecialChars(value[key], keywords);
    }
  } else {
    // other type direct return value.
  }
  return value;
}
