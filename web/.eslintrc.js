module.exports = {
    root: true,
    env: { node: true },
    // https://github.com/vuejs/vue-eslint-parser#parseroptionsparser
    parser: "vue-eslint-parser",
    parserOptions: {
      parser: "@typescript-eslint/parser",
    },
    plugins: ["@typescript-eslint", "prettier"],
    extends: [
      "plugin:vue/vue3-recommended",
      "prettier",
    ],
    rules: {
      "prettier/prettier": "error",
    }
  }