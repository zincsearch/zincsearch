import http from "./http";

var index = {
  list: () => {
    return http().get("/api/index");
  },
  update: (data) => {
    return http().put("/api/index/" + data.name, data);
  },
  delete: (name) => {
    return http().delete("/api/index/" + name);
  },
};

export default index;
