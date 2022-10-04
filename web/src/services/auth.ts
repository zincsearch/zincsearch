import http from "./http";

var auth = {
  login: (data: any) => {
    return http().post("/api/login", data);
  },
  refresh() {
    return http().get("/api/login/refresh");
  },
};

export default auth;
