import http from "./http";

var template = {
  list: () => {
    return http().get("/es/_index_template");
  },
  update: (data: any) => {
    return http().put("/es/_index_template/" + data.name, data);
  },
  delete: (name: string) => {
    return http().delete("/es/_index_template/" + name);
  },
};

export default template;
