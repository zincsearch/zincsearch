import http from "./http";

var user = {
  list: () => {
    return http().get("/api/user");
  },
  update: (data: any) => {
    return http().put("/api/user", data);
  },
  delete: (id: string) => {
    return http().delete("/api/user/" + id);
  },
};

export default user;
