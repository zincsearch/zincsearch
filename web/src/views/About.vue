<template>
  <q-page class="q-pa-lg">
    <h3 class="q-ma-none">Zinc Search</h3>
    <p class="q-mt-md">{{ t("about.introduction") }}</p>

    <div class="q-pa-md">
      <div class="row items-center q-mb-md">
        <div class="col-sm-3 col-12">Version</div>
        <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
          {{ version.version }}
        </div>
      </div>
      <div class="row items-center q-mb-md">
        <div class="col-sm-3 col-12">Build</div>
        <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
          {{ version.build }}
        </div>
      </div>
      <div class="row items-center q-mb-md">
        <div class="col-sm-3 col-12">CommitHash</div>
        <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
          {{ version.commit_hash }}
        </div>
      </div>
      <div class="row items-center q-mb-md">
        <div class="col-sm-3 col-12">Branch</div>
        <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
          {{ version.branch }}
        </div>
      </div>
      <div class="row items-center q-mb-md">
        <div class="col-sm-3 col-12">BuildDate</div>
        <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
          {{ version.build_date }}
        </div>
      </div>
    </div>
  </q-page>
</template>

<script>
import { defineComponent, ref } from "vue";
import { useStore } from "vuex";
import aboutService from "../services/about";
import { useI18n } from "vue-i18n";

export default defineComponent({
  name: "PageAbout",

  setup() {
    const { t } = useI18n();
    const store = useStore();

    const version = ref({});
    const getVersion = () => {
      aboutService.get().then((res) => {
        version.value = res.data;
      });
    };
    getVersion();

    return {
      t,
      version,
    };
  },
});
</script>
