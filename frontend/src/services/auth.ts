import http from "./http";

var auth = {
  login: (data: any) => {
    return http().post("/api/login", data);
  },
  updateAccount: (data: { _id: string; password: string; name?: string; new_password?: string }) => {
    // This endpoint verifies the current password in the body, not stored Basic credentials.
    return http({ authenticated: false }).put("/api/account", data);
  },
};

export default auth;
