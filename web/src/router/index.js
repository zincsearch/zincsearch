import { createRouter, createWebHistory } from "vue-router";
import Search from "../views/Search.vue";
import LogIn from "../views/LogIn.vue";
import Users from "../views/Users.vue";
import About from "../views/About.vue";
import IndexManagement from "../views/IndexManagement.vue";

import store from "../store/index.js";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Search,
  },
  {
    path: "/search",
    name: "Search",
    component: Search,
  },
  {
    path: "/login",
    name: "LogIn",
    component: LogIn,
  },
  {
    path: "/indexmanagement",
    name: "IndexManagement",
    component: IndexManagement,
  },
  {
    path: "/users",
    name: "Users",
    component: Users,
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: About,
    component1: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(window.location.pathname),
  routes,
});

router.beforeEach((to, from, next) => {
  var isAuthenticated = store.state.user.isLoggedIn;

  if (to.path !== "/login" && !isAuthenticated) {
    next({ path: "/login" });
  } else {
    next();
  }
});

export default router;
