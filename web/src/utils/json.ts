export const ToString = (o: any) => {
  if (!o) {
    return "";
  }
  if (typeof o == "string") {
    return o;
  } else if (Array.isArray(o)) {
    let tmp = "";
    for (let i = 0; i < o.length; i++) {
      tmp += (tmp === "" ? "" : ", ") + ToString(o[i]);
    }
    return `[${tmp}]`;
  } else if (typeof o == "object") {
    let tmp = "";
    for (const key in o) {
      tmp += (tmp === "" ? "" : ", ") + `"${key}":` + ToString(o[key]);
    }
    return `{${tmp}}`;
  } else {
    return o;
  }
};

export const byString = (o: any, s: string) => {
  if (s == undefined) {
    return "";
  }
  if (s in o) {
    return ToString(o[s]);
  }
  s = s.replace(/\[(\w+)\]/g, ".$1"); // convert indexes to properties
  s = s.replace(/^\./, ""); // strip a leading dot
  let a = s.split(".");
  for (let i = 0, n = a.length; i < n; ++i) {
    const k = a[i];
    if (typeof o == "object" && k in o) {
      o = o[k];
    }
  }
  return ToString(o);
};

export const deepKeys = (o: any) => {
  if (!(o instanceof Object)) {
    return [];
  }
  let results = [];
  let keys = Object.keys(o);
  for (var i in keys) {
    if (o[keys[i]] == undefined || o[keys[i]].length) {
      results.push(keys[i]);
    } else {
      let subKeys = deepKeys(o[keys[i]]);
      if (subKeys.length > 0) {
        subKeys.forEach((key) => {
          results.push(keys[i] + "." + key);
        });
      } else {
        results.push(keys[i]);
      }
    }
  }
  return results;
};
