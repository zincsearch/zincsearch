import { createRouter, createWebHistory } from "vue-router";
import store from "../store/index.js";

import About from "../views/About.vue";
import Index from "../views/Index.vue";
import Login from "../views/Login.vue";
import Search from "../views/Search.vue";
import Template from "../views/Template.vue";
import User from "../views/User.vue";

const routes = [
  {
    path: "/",
    component: () => import("../layouts/MainLayout.vue"),
    children: [
      {
        path: "",
        name: "home",
        redirect: (to: any) => {
          return { path: "/search", query: { q: to.params.searchText } };
        },
      },
      {
        path: "search",
        name: "search",
        component: Search,
        meta: {
          keepAlive: true,
        },
      },
      {
        path: "index",
        name: "index",
        component: Index,
        meta: {
          keepAlive: true,
        },
        role: true,
      },
      {
        path: "template",
        name: "template",
        component: Template,
        meta: {
          keepAlive: true,
        },
        role: true,
      },
      {
        path: "user",
        name: "user",
        component: User,
        role: true,
      },
      {
        path: "about",
        name: "about",
        component: About,
      },
    ],
  },
  {
    path: "/login",
    name: "login",
    component: Login,
  },
  // Always leave this as last one,
  // but you can also remove it
  {
    path: "/:catchAll(.*)*",
    component: () => import("../views/Error404.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(window.location.pathname),
  routes,
});

router.beforeEach((to: any, from: any, next: any) => {
  var isAuthenticated = store.state.user.isLoggedIn;
  var role = store.state.user.role === "admin";

  if (to.path !== "/login" && !isAuthenticated) {
    next({ path: "/login" });
  } else if (to.role && role) {
    next();
  } else {
    next();
  }
});

export default router;
