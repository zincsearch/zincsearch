<template>
  <div icon="info" class="date-time-selector">
    <q-btn
      outline
      :label="displayValue"
      icon-right="schedule"
      class="date-time-button"
    >
    </q-btn>
    <q-menu class="date-time-dialog">
      <q-tabs dense v-model="message.tab" class="text-teal">
        <q-tab name="relative" label="relative"> </q-tab>
        <q-tab name="absolute" label="absolute"> </q-tab>
      </q-tabs>

      <q-tab-panels v-model="message.tab" animated>
        <q-tab-panel name="relative">
          <table class="relative-period-table">
            <tbody>
              <tr>
                <td class="relative-period-name">Minutes</td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Minutes' &&
                      message.selectedRelativeValue == 5
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="5"
                    @click="selectRelativeValue('Minutes', 5)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Minutes' &&
                      message.selectedRelativeValue == 10
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="10"
                    @click="selectRelativeValue('Minutes', 10)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Minutes' &&
                      message.selectedRelativeValue == 15
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="15"
                    @click="selectRelativeValue('Minutes', 15)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Minutes' &&
                      message.selectedRelativeValue == 30
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="30"
                    @click="selectRelativeValue('Minutes', 30)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Minutes' &&
                      message.selectedRelativeValue == 45
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="45"
                    @click="selectRelativeValue('Minutes', 45)"
                    outline
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Hours</td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Hours' &&
                      message.selectedRelativeValue == 1
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="1"
                    @click="selectRelativeValue('Hours', 1)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Hours' &&
                      message.selectedRelativeValue == 2
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="2"
                    @click="selectRelativeValue('Hours', 2)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Hours' &&
                      message.selectedRelativeValue == 3
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="3"
                    @click="selectRelativeValue('Hours', 3)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Hours' &&
                      message.selectedRelativeValue == 6
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="6"
                    @click="selectRelativeValue('Hours', 6)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Hours' &&
                      message.selectedRelativeValue == 8
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="8"
                    @click="selectRelativeValue('Hours', 8)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Hours' &&
                      message.selectedRelativeValue == 12
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="12"
                    @click="selectRelativeValue('Hours', 12)"
                    outline
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Days</td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Days' &&
                      message.selectedRelativeValue == 1
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="1"
                    @click="selectRelativeValue('Days', 1)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Days' &&
                      message.selectedRelativeValue == 2
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="2"
                    @click="selectRelativeValue('Days', 2)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Days' &&
                      message.selectedRelativeValue == 3
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="3"
                    @click="selectRelativeValue('Days', 3)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Days' &&
                      message.selectedRelativeValue == 4
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="4"
                    @click="selectRelativeValue('Days', 4)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Days' &&
                      message.selectedRelativeValue == 5
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="5"
                    @click="selectRelativeValue('Days', 5)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Days' &&
                      message.selectedRelativeValue == 6
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="6"
                    @click="selectRelativeValue('Days', 6)"
                    outline
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Weeks</td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Weeks' &&
                      message.selectedRelativeValue == 1
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="1"
                    @click="selectRelativeValue('Weeks', 1)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Weeks' &&
                      message.selectedRelativeValue == 2
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="2"
                    @click="selectRelativeValue('Weeks', 2)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Weeks' &&
                      message.selectedRelativeValue == 3
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="3"
                    @click="selectRelativeValue('Weeks', 3)"
                    outline
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Months</td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Months' &&
                      message.selectedRelativeValue == 1
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="1"
                    @click="selectRelativeValue('Months', 1)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Months' &&
                      message.selectedRelativeValue == 2
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="2"
                    @click="selectRelativeValue('Months', 2)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Months' &&
                      message.selectedRelativeValue == 3
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="3"
                    @click="selectRelativeValue('Months', 3)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Months' &&
                      message.selectedRelativeValue == 4
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="4"
                    @click="selectRelativeValue('Months', 4)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Months' &&
                      message.selectedRelativeValue == 5
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="5"
                    @click="selectRelativeValue('Months', 5)"
                    outline
                  />
                </td>
                <td>
                  <q-btn
                    :class="
                      message.selectedRelativePeriod == 'Months' &&
                      message.selectedRelativeValue == 6
                        ? 'rp-selector-selected'
                        : 'rp-selector'
                    "
                    label="6"
                    @click="selectRelativeValue('Months', 6)"
                    outline
                  />
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">Custom</td>
                <td colspan="3" class="custom-relative">
                  <q-input
                    type="number"
                    style="text-align: right"
                    dense
                    filled
                    v-model="message.selectedRelativeValue"
                  ></q-input>
                </td>
                <td colspan="3">
                  <q-select
                    v-model="message.selectedRelativePeriod"
                    :options="relativePeriods"
                    dense
                    filled
                  ></q-select>
                </td>
              </tr>
              <tr>
                <td class="relative-period-name">FullTime</td>
                <td colspan="3" class="custom-relative">
                  <q-toggle v-model="message.selectedFullTime" color="green" />
                </td>
                <td colspan="3"></td>
              </tr>
            </tbody>
          </table>
        </q-tab-panel>
        <q-tab-panel name="absolute">
          <table>
            <tbody>
              <tr>
                <td class="absolute-period-name">Start Date</td>
                <td class="absolute-period-name">End Date</td>
              </tr>
              <tr>
                <td>
                  <q-input
                    dense
                    filled
                    v-model="message.startDate"
                    mask="date"
                    :rules="['date']"
                  >
                    <template v-slot:append>
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
                <td class="absolute-period-name">
                  <q-input
                    dense
                    filled
                    v-model="message.endDate"
                    mask="date"
                    :rules="['date']"
                  >
                    <template v-slot:append>
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
                <td class="absolute-period-name">Start Time</td>
                <td class="absolute-period-name">End Time</td>
              </tr>
              <tr>
                <td class="absolute-period-name">
                  <q-input
                    dense
                    filled
                    v-model="message.startTime"
                    mask="time"
                    :rules="['time']"
                  >
                    <template v-slot:append>
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
                <td class="absolute-period-name">
                  <q-input
                    dense
                    filled
                    v-model="message.endTime"
                    mask="time"
                    :rules="['time']"
                  >
                    <template v-slot:append>
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

    return {
      message,
      relativePeriods,

      selectRelativeValue(period, value) {
        this.message.selectedRelativeValue = value;
        this.message.selectedRelativePeriod = period;
      },
    };
  },
};
</script>

<style scoped>
.date-time-selector {
  /* display: flex; */
  /* align-items: center; */
  justify-content: space-between;
  margin-top: 10px;
  margin-right: 10px;
  height: 50px;
}

.date-time-button {
  /* display: block; */
  vertical-align: middle;
  text-align: center;
  height: 40px;
  width: 100%;
}

.relative-period {
  display: flex;
  flex-direction: row;
  /* justify-content: space-between; */
  align-items: center;
  height: 50px;
  width: 450px;
  margin-left: 10px;
  border: brown;
}

.rp-selector {
  height: 40px;
  width: 50px;
  margin-left: 10px;
  margin-bottom: 5px;
  border: brown;
}

.rp-selector-selected {
  height: 40px;
  width: 50px;
  margin-left: 10px;
  margin-bottom: 5px;
  color: rgb(224, 64, 64);
  font-weight: bolder;
  /* font-style: oblique; */
  font-size: 1.5em;

  /* font-family: fantasy; */
}

.relative-period-name {
  width: 50px;
  margin-left: 10px;
}

.relative-period-table {
  margin: 10px;
}

.custom-relative {
  /* display: flex; */
  /* flex-direction: row; */
  justify-content: space-between;
  align-items: center;
  height: 50px;
  /* width: 450px; */
  margin-left: 10px;
  /* border: brown; */
}

.date-time-dialog {
  width: 600px;
}

.absolute-period-name {
  /* width: 50px; */
  margin: 20px;
}
</style>
