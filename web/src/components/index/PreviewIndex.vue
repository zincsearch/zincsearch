<template>
  <q-card class="column full-height">
    <q-card-section>
      <div class="row items-center no-wrap">
        <div class="col">
          <div class="text-h6">{{ indexData.name }}</div>
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
        <q-tab name="settings" label="Settings" />
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
                {{ indexData.name }}
              </div>
            </div>
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">Docs Count</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ indexData.doc_num }}
              </div>
            </div>
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">Shards Num</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ indexData.shard_num }}
              </div>
            </div>
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">Storage Size</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ indexData.storage_size }}
              </div>
            </div>
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">Storage Type</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ indexData.storage_type }}
              </div>
            </div>
            <div class="row items-center q-mb-md">
              <div class="col-sm-3 col-12">WAL Entries</div>
              <div class="col-sm-9 col-12 q-mb-none q-pl-md q-pt-sm q-pb-sm">
                {{ indexData.wal_size }}
              </div>
            </div>
          </div>
        </q-tab-panel>

        <q-tab-panel name="settings" class="q-px-none">
          <div class="q-py-md">
            <json-editor
              v-model="indexData.settings"
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
              v-model="indexData.mappings"
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
              v-model="indexData"
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
  name: "PreviewIndex",
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
      indexData: ref({}),
      tab: ref("summary"),
    };
  },
  created() {
    if (this.modelValue && this.modelValue.name) {
      this.indexData["name"] = this.modelValue.name;
      this.indexData["doc_num"] = this.modelValue.doc_num;
      this.indexData["shard_num"] = this.modelValue.shard_num;
      this.indexData["storage_type"] = this.modelValue.storage_type;
      this.indexData["storage_size"] = this.modelValue.storage_size;
      this.indexData["wal_size"] = this.modelValue.wal_size || 0;
      this.indexData["settings"] = this.modelValue.settings || {};
      this.indexData["mappings"] = this.modelValue.mappings || {};
    }
  },
});
</script>
