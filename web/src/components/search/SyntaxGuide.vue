<template>
  <q-btn
    data-cy="syntax-guide-button"
    size="md"
    unelevated
    class="q-pa-sm syntax-menu-button"
    color="grey-4"
    text-color="grey-7"
    icon="info"
  >
    <q-menu>
      <q-card flat>
        <q-card-section class="syntax-guide-title">
          <div class="text-h6">{{t('search.syntaxGuide')}}</div>
        </q-card-section>
        <q-separator />
        <q-card-section class="q-pt-none answers">
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">Search for Gold</div>
            <div class="syntax-guide-text">Gold</div>
          </div>
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">Search for City with Paris</div>
            <div class="syntax-guide-text">City:Paris</div>
          </div>
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">Search for City Not Paris</div>
            <div class="syntax-guide-text">-City:Paris</div>
          </div>
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">
              Search for City with Paris or Gold
            </div>
            <div class="syntax-guide-text">City:Paris Gold</div>
          </div>
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">
              Search for City with Paris and Gold
            </div>
            <div class="syntax-guide-text">+City:Paris +Gold</div>
          </div>
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">
              Search for City with Paris but not Gold
            </div>
            <div class="syntax-guide-text">+City:Paris -Gold</div>
          </div>
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">
              Search for Medal Gold and Year>2000
            </div>
            <div class="syntax-guide-text">+Medal:Gold +Year:>2000</div>
          </div>
          <div class="syntax-section">
            <div class="syntax-guide-sub-title">
              Search text starting with par . e.g. Paris, Part, Paramater
            </div>
            <div class="syntax-guide-text">par*</div>
          </div>
        </q-card-section>
      </q-card>
    </q-menu>
  </q-btn>
</template>


<script>
import { defineComponent } from "vue";
import { useI18n } from "vue-i18n";

export default defineComponent({
  name: "ComponentSearchSyntaxGuide",
  setup() {
    const { t } = useI18n();
    return {
      t,
    };
  },
});
</script>

<style lang="scss">
.syntax-guide {
  /* width: 80%;*/

  margin-right: 10px;
}

.syntax-menu-button {
  border-radius: 0 4px 4px 0;
}

.syntax-guide-title {
  width: 420px;
}

.syntax-guide-sub-title {
  color: $primary;
  font-size: 15px;
  margin-left: 5px;
}

.syntax-guide-text {
  font-size: 12px;
  margin-left: 5px;
}

.syntax-section {
  margin-bottom: 5px;
}

.answers {
  margin-top: 10px;
}
</style>
