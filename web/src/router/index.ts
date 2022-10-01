import { createRouter, createWebHistory } from "vue-router";
import store from "../store/index.js";

import About from "../views/About.vue";
import Index from "../views/Index.vue";
import Login from "../views/Login.vue";
import Search from "../views/Search.vue";
import Template from "../views/Template.vue";
import User from "../views/User.vue";
import Role from "../views/Role.vue";

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
      },
      {
        path: "template",
        name: "template",
        component: Template,
        meta: {
          keepAlive: true,
        },
      },
      {
        path: "user",
        name: "user",
        component: User,
      },
      {
        path: "role",
        name: "role",
        component: Role,
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

  if (to.path !== "/login" && !isAuthenticated) {
    next({ path: "/login" });
  } else {
    next();
  }
});

export default router;
