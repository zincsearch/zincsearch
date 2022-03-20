import http from "./http";

var auth = {
  login: (data) => {
    return http().post("/api/login", data);
  },
};

export default auth;
