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
  isLoggedIn: () => {
    return http().get("/api/login/verify");
  },
  logout() {
    return http().post("/api/logout");
  },
};

export default user;
