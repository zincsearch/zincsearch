<template>
  <div class="about">
    <h1>{{ t("about.name") }}</h1>
    <div>{{ t("about.introduction") }}</div>
  </div>
</template>
<script>
import { useI18n } from "vue-i18n";
export default {
  setup() {
    const { t } = useI18n();
    return {
      t,
    };
  },
};
</script>
