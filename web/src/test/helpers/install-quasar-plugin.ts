import { config } from "@vue/test-utils";
import { cloneDeep } from "lodash-es";
import { Quasar, QuasarPluginOptions } from "quasar";
import { qLayoutInjections } from "./layout-injections";
import { beforeAll, afterAll } from "vitest";

export function installQuasar(options?: Partial<QuasarPluginOptions>) {
  const globalConfigBackup = cloneDeep(config.global);

  beforeAll(() => {
    config.global.plugins.unshift([Quasar, options]);
    config.global.provide = {
      ...config.global.provide,
      ...qLayoutInjections(),
    };
  });

  afterAll(() => {
    config.global = globalConfigBackup;
  });
}
