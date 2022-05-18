import http from "./http";

var about = {
  get: () => {
    return http().get("/version");
  },
};

export default about;
