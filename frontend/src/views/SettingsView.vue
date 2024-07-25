<script setup>
  import { reactive, watch, ref } from 'vue'
  import { store } from '@/store'
  import rules from '@/rules'
  import { mdiArrowLeft } from "@mdi/js"


  const { autorefresh_interval: init_autorefresh_interval } = store

  const settings = reactive({
    autorefresh_interval: init_autorefresh_interval
  })

  const formIsValid = ref(false)
  const unsaved = ref(false)

  const saveSettings = () => {
    store.autorefresh_interval = settings.autorefresh_interval
  }

  watch(settings, () => {
    if (!unsaved.value) {
      unsaved.value = true
    }
  })
</script>

<template>
  <v-container>
    <v-row>
        <v-col class="pa-0">
          <v-btn @click="$router.replace({ name: 'tunnels' })" class="rounded-button pa-0" size="x-large" variant="text"
            :icon="mdiArrowLeft" />
        </v-col>
        <v-col class="d-flex justify-end pa-0">
          <h1 class="header">Settings</h1>
        </v-col>
      </v-row>
    <v-row>
      <p>Client settings are saved locally in browser!</p>
    </v-row>
    <v-row>
      <v-col>
        <v-form v-model="formIsValid">
          <v-card rounded="xl" class="settings-card pa-3 container" variant="flat">
            <v-card-title>Auto-refresh interval</v-card-title>
            <v-card-text>
              <v-row class="d-flex align-center" >
                <v-col cols="5">
                  <v-text-field :rules="[rules.required, rules.integers]" v-model="settings.autorefresh_interval" />
                </v-col>
                <v-col>
                  <p>Seconds</p>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-form>
      </v-col>
    </v-row>
    <v-row>
    </v-row>
    <v-snackbar timeout="-1" v-model="unsaved">
      You have unsaved changes!
      <template v-slot:actions>
        <v-btn @click="saveSettings(); unsaved = false" :disabled="!formIsValid" variant="text">
          Save
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<style>
.settings-card {
  width: 100%;
}
</style>
