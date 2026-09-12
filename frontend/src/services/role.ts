import http from "./http";

var role = {
  list: () => {
    return http().get("/api/role");
  },
  update: (data: any) => {
    return http().put("/api/role", data);
  },
  delete: (id: string) => {
    return http().delete("/api/role/" + id);
  },
};

export default role;
