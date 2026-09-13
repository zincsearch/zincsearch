import http from "./http";

const search = {
  search: ({ index, query }: { index: string; query: Record<string, unknown> }) => {
    let url = "/es/_search";
    if (index != "") {
      url = "/es/" + index + "/_search";
    }
    return http().post(url, query);
  },
};

export default search;
