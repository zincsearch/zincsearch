import http from "./http";

var auth = {
  login: (data: any) => {
    return http().post("/api/login", data);
  },
};

export default auth;
