import http from "./http";

var permission = {
  list: () => {
    return http().get("/api/permission");
  },
};

export default permission;
