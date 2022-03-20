import http from "./http";

var template = {
  list: () => {
    return http().get("/es/_index_template");
  },
  update: (data) => {
    return http().put("/es/_index_template/" + data.name, data);
  },
  delete: (name) => {
    return http().delete("/es/_index_template/" + name);
  },
};

export default template;
