import http from "./http";

var index = {
  list: (
    page_num: number,
    page_size: number,
    sort_by: string,
    desc: boolean,
    name: string
  ) => {
    return http().get(
      `/api/index?page_num=${page_num}&page_size=${page_size}&sort_by=${sort_by}&desc=${desc}&name=${name}`
    );
  },
  update: (data: any) => {
    return http().put("/api/index/" + data.name, data);
  },
  delete: (names: string) => {
    return http().delete("/api/index/" + names);
  },
  nameList: (name: string) => {
    return http().get("/api/index_name?name=" + name);
  },
};

export default index;
