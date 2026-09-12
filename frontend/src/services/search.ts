import http from "./http";

var search = {
  search: ({ index, query }: { index: string; query: string }) => {
    let url = "/es/_search";
    if (index != "") {
      url = "/es/" + index + "/_search";
    }
    return http().post(url, query);
  },
};

export default search;
