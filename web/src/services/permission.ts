import http from "./http";

var permission = {
  list: () => {
    return http().get("/api/permissions");
  },
};

export default permission;
