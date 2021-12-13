<template>
  <q-btn outline class="syntax-menu-button" color="secondary" icon="info">
    <q-menu>
      <q-card>
        <q-card-section class="syntax-guide-title">
          <div>Syntax guide</div>
        </q-card-section>
        <!-- <q-separator /> -->
        <div class="answers">
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
              Search text starting with par . e.g. Paris, Parthenon, Part, Paramater
            </div>
            <div class="syntax-guide-text">par*</div>
          </div>
        </div>
      </q-card>
    </q-menu>
  </q-btn>
</template>

<style scoped>
.syntax-guide {
  /* width: 80%; */
  margin-right: 10px;
}

.syntax-menu-button {
  width: 50px;
  height: 40px;
  margin-top: 10px;
  margin-right: 10px;
}

.syntax-guide-title {
  /* font-family: "Roboto Mono", Consolas, Menlo, Courier, monospace; */
  width: 600px;
  font-size: 20px;
  background-color: rgb(207, 203, 203);
}

.syntax-guide-sub-title {
  /* font-family: "Roboto Mono", Consolas, Menlo, Courier, monospace; */
  color: green;
  /* font-weight: bold; */
  font-size: 15px;
  margin-left: 5px;
}

.syntax-guide-text {
  font-family: "Roboto Mono", Consolas, Menlo, Courier, monospace;
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
