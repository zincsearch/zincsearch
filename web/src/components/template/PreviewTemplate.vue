<template>
  <q-card class="column full-height">
    <q-card-section>
      <div class="row items-center no-wrap">
        <div class="col">
          <div class="text-h6">{{ templateData.name }}</div>
        </div>

        <div class="col-auto">
          <q-btn v-close-popup flat round color="grey-7" icon="close" />
        </div>
      </div>
    </q-card-section>

    <q-card-section class="col q-pt-none q-w-p50">
      <q-tabs
        v-model="tab"
        dense
        no-caps
        narrow-indicator
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="justify"
      >
        <q-tab name="summary" label="Summary" />
        <q-tab name="settings" label="Index settings" />
        <q-tab name="mappings" label="Mappings" />
        <q-tab name="preview" label="Preview" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="tab">
        <q-tab-panel name="summary">
          <div class="q-pa-md">
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">Name</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ templateData.name }}
              </div>
            </div>
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">Index patterns</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ templateData.patterns }}
              </div>
            </div>
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">Priority</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ templateData.priority }}
              </div>
            </div>
          </div>
        </q-tab-panel>

        <q-tab-panel name="settings" class="q-px-none">
          <div class="q-py-md">
            <json-editor
              v-model="templateData.template.settings"
              name="settings"
              mode="preview"
              :readonly="true"
              :height="501"
            ></json-editor>
          </div>
        </q-tab-panel>

        <q-tab-panel name="mappings" class="q-px-none">
          <div class="q-py-md">
            <json-editor
              v-model="templateData.template.mappings"
              name="mappings"
              mode="preview"
              :readonly="true"
              :height="501"
            ></json-editor>
          </div>
        </q-tab-panel>

        <q-tab-panel name="preview" class="q-px-none">
          <div class="q-py-md">
            <json-editor
              v-model="templateData"
              name="preview"
              mode="preview"
              :readonly="true"
              :height="501"
            ></json-editor>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card-section>
  </q-card>
</template>

<script>
import { defineComponent, ref } from "vue";
import JsonEditor from "../JsonEditor.vue";

export default defineComponent({
  name: "PreviewTemplate",
  components: {
    JsonEditor,
  },
  props: {
    modelValue: {
      type: Object,
      default: () => {},
    },
  },
  setup() {
    return {
      templateData: ref({}),
      tab: ref("summary"),
    };
  },
  created() {
    if (this.modelValue && this.modelValue.name) {
      this.templateData["name"] = this.modelValue.name;
      this.templateData["patterns"] = this.modelValue.patterns;
      this.templateData["priority"] = this.modelValue.priority;
      this.templateData["template"] = {
        settings: this.modelValue.template.settings || {},
        mappings: this.modelValue.template.mappings || {},
      };
    }
  },
});
</script>
