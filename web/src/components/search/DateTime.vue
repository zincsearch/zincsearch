<template>
  <div icon="info" class="justify-between">
    <q-btn
      data-cy="date-time-button"
      outline
      no-caps
      :label="displayValue"
      align="between"
      icon-right="schedule"
      class="q-pa-sm date-time-button"
      color="grey-9"
    >
    </q-btn>
    <q-menu class="date-time-dialog">
      <q-tabs v-model="message.tab" dense class="text-primary">
        <q-tab name="relative" label="relative"> </q-tab>
        <q-tab name="absolute" label="absolute"> </q-tab>
      </q-tabs>

      <q-tab-panels v-model="message.tab" animated>
        <q-tab-panel name="relative">
          <table class="date-time-table">
            <tbody>
              <tr>
                <td class="relative-period-name">Minutes</td>
                <td v-for="item in relativePeriodList['Minutes']" :key="item">
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Minutes' &&
                      message.selectedRelativeValue == item
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    :label="item"
                    outline
                    dense
                    @click="selectRelativeValue('Minutes', item)"
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Hours</td>
                <td v-for="item in relativePeriodList['Hours']" :key="item">
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Hours' &&
                      message.selectedRelativeValue == item
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    :label="item"
                    outline
                    dense
                    @click="selectRelativeValue('Hours', item)"
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Days</td>
                <td v-for="item in relativePeriodList['Days']" :key="item">
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Days' &&
                      message.selectedRelativeValue == item
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    :label="item"
                    outline
                    dense
                    @click="selectRelativeValue('Days', item)"
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Weeks</td>
                <td v-for="item in relativePeriodList['Weeks']" :key="item">
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Weeks' &&
                      message.selectedRelativeValue == item
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    :label="item"
                    outline
                    dense
                    @click="selectRelativeValue('Weeks', item)"
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Months</td>
                <td v-for="item in relativePeriodList['Months']" :key="item">
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Months' &&
                      message.selectedRelativeValue == item
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    :label="item"
                    outline
                    dense
                    @click="selectRelativeValue('Months', item)"
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Custom</td>
                <td colspan="6">
                  <div class="row q-gutter-sm">
                    <div class="col">
                      <q-input
                        v-model="message.selectedRelativeValue"
                        type="number"
                        dense
                        filled
                      ></q-input>
                    </div>
                    <div class="col">
                      <q-select
                        v-model="message.selectedRelativePeriod"
                        :options="relativePeriods"
                        dense
                        filled
                      ></q-select>
                    </div>
                  </div>
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">FullTime</td>
                <td colspan="6">
                  <q-toggle v-model="message.selectedFullTime" color="green" />
                </td>
              </tr>
            </tbody>
          </table>
        </q-tab-panel>
        <q-tab-panel name="absolute">
          <table class="date-time-table">
            <tbody>
              <tr>
                <td>Start Date</td>
                <td>End Date</td>
              </tr>
              <tr>
                <td>
                  <q-input
                    v-model="message.startDate"
                    dense
                    filled
                    mask="date"
                    :rules="['date']"
                  >
                    <template #append>
                      <q-icon name="event" class="cursor-pointer">
                        <q-popup-proxy
                          ref="qDateProxy"
                          transition-show="scale"
                          transition-hide="scale"
                        >
                          <q-date v-model="message.startDate">
                            <div class="row items-center justify-end">
                              <q-btn
                                v-close-popup
                                label="Close"
                                color="primary"
                                flat
                              />
                            </div>
                          </q-date>
                        </q-popup-proxy>
                      </q-icon>
                    </template>
                  </q-input>
                </td>
                <td>
                  <q-input
                    v-model="message.endDate"
                    dense
                    filled
                    mask="date"
                    :rules="['date']"
                  >
                    <template #append>
                      <q-icon name="event" class="cursor-pointer">
                        <q-popup-proxy
                          ref="qDateProxy"
                          transition-show="scale"
                          transition-hide="scale"
                        >
                          <q-date v-model="message.endDate">
                            <div class="row items-center justify-end">
                              <q-btn
                                v-close-popup
                                label="Close"
                                color="primary"
                                flat
                              />
                            </div>
                          </q-date>
                        </q-popup-proxy>
                      </q-icon>
                    </template>
                  </q-input>
                </td>
              </tr>
              <tr>
                <td>Start Time</td>
                <td>End Time</td>
              </tr>
              <tr>
                <td>
                  <q-input
                    v-model="message.startTime"
                    dense
                    filled
                    mask="time"
                    :rules="['time']"
                  >
                    <template #append>
                      <q-icon name="access_time" class="cursor-pointer">
                        <q-popup-proxy
                          transition-show="scale"
                          transition-hide="scale"
                        >
                          <q-time v-model="message.startTime">
                            <div class="row items-center justify-end">
                              <q-btn
                                v-close-popup
                                label="Close"
                                color="primary"
                                flat
                              />
                            </div>
                          </q-time>
                        </q-popup-proxy>
                      </q-icon>
                    </template>
                  </q-input>
                </td>
                <td>
                  <q-input
                    v-model="message.endTime"
                    dense
                    filled
                    mask="time"
                    :rules="['time']"
                  >
                    <template #append>
                      <q-icon name="access_time" class="cursor-pointer">
                        <q-popup-proxy
                          transition-show="scale"
                          transition-hide="scale"
                        >
                          <q-time v-model="message.endTime">
                            <div class="row items-center justify-end">
                              <q-btn
                                v-close-popup
                                label="Close"
                                color="primary"
                                flat
                              />
                            </div>
                          </q-time>
                        </q-popup-proxy>
                      </q-icon>
                    </template>
                  </q-input>
                </td>
              </tr>
            </tbody>
          </table>
        </q-tab-panel>
      </q-tab-panels>
    </q-menu>
  </div>
</template>

<script>
import { ref, computed } from "vue";

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({
        selectedRelativePeriod: "Days",
        selectedRelativeValue: 1,
        selectedFullTime: false,
        startDate: "",
        endDate: "",
        startTime: "",
        endTime: "",
      }),
    },
  },
  emits: ["update:modelValue"],
  setup(props, { emit }) {
    const message = computed({
      get: () => props.modelValue,
      set: (value) => emit("update:modelValue", value),
    });

    const relativePeriods = ref([
      "Minutes",
      "Hours",
      "Days",
      "Weeks",
      "Months",
    ]);
    const relativePeriodList = {
      Minutes: [1, 5, 10, 15, 30, 45],
      Hours: [1, 2, 3, 6, 8, 12],
      Days: [1, 2, 3, 4, 5, 6],
      Weeks: [1, 2, 3],
      Months: [1, 2, 3, 4, 5, 6],
    };

    return {
      message,
      relativePeriods,
      relativePeriodList,
      selectRelativeValue(period, value) {
        this.message.selectedRelativeValue = value;
        this.message.selectedRelativePeriod = period;
      },
    };
  },
  computed: {
    displayValue() {
      if (this.message.selectedFullTime) {
        return "FullTime";
      }
      if (this.message.tab === "relative") {
        return `${this.message.selectedRelativeValue} ${this.message.selectedRelativePeriod}`;
      } else {
        return `${this.message.startDate} ${this.message.startTime} - ${this.message.endDate} ${this.message.endTime}`;
      }
    },
  },
};
</script>

<style lang="scss">
.date-time-button {
  min-width: 138px;
}

.date-time-dialog {
  width: 370px;
}
.date-time-table {
  width: 100%;
}

.date-time-table td {
  padding: 0 2px 6px 2px;
}

.relative-period-name {
  width: 35px;
}

.rp-selector,
.rp-selector-selected {
  height: 32px;
  width: 35px;
  border: $secondary;
}

.rp-selector-selected {
  color: $secondary;
  font-weight: bolder;
}
</style>
