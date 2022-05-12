import { createStore } from "vuex";

export default createStore({
  state: {
    API_ENDPOINT: "http://localhost:4080",
    user: {
      isLoggedIn: false,
      _id: "",
      password: "",
      base64encoded: "",
      name: "",
      email: "",
      role: "",
    },
  },
  mutations: {
    login(state, payload) {
      if (payload && payload._id && payload.base64encoded) {
        state.user.isLoggedIn = true;
        state.user._id = payload._id;
        state.user.name = payload.name || payload._id;
        state.user.role = payload.role;
        state.user.base64encoded = payload.base64encoded;
      }
    },
    logout(state) {
      state.user.isLoggedIn = false;
      state.user._id = "";
      state.user.name = "";
      state.user.role = "";
      state.user.base64encoded = "";
    },
    endpoint(state, payload) {
      state.API_ENDPOINT = payload;
    },
  },
  actions: {
    login(context, payload) {
      context.commit("login", payload);
    },
    logout(context) {
      context.commit("logout");
    },
    endpoint(context, payload) {
      context.commit("endpoint", payload);
    },
  },
  modules: {},
});
