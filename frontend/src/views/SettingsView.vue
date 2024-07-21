<script setup>
  import { reactive, watch, ref } from 'vue'
  import { store } from '@/store'

  const { autorefresh_interval: init_autorefresh_interval } = store

  const settings = reactive({
    autorefresh_interval: init_autorefresh_interval
  })

  const unsaved = ref(false)

  const saveSettings = () => {
    store.autorefresh_interval = settings.autorefresh_interval
  }

  watch(settings, () => {
    if (!unsaved.value) {
      unsaved.value = true
      console.log("You have unsaved changes!")
    }
  })
</script>

<template>
  <v-container>
    <v-row>
      <p>Client settings are saved locally in browser!</p>
    </v-row>
    <v-row>
      <v-col>
        <v-card rounded="xl" class="settings-card pa-3 container" variant="flat">
          <v-card-title>Auto-refresh interval</v-card-title>
          <v-card-text>
            <v-row class="d-flex align-center">
              <v-col cols="5">
                <v-text-field v-model="settings.autorefresh_interval" />
              </v-col>
              <v-col>
                <p>Seconds</p>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
    </v-row>
    <v-snackbar timeout="-1" v-model="unsaved">
      You have unsaved changes!
      <template v-slot:actions>
        <v-btn @click="saveSettings(); unsaved = false" variant="text">
          Save
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<style>
.settings-card {
  width: 100%;
  background-color: var(--color-surface-container);
  color: var(--color-on-surface);
}
</style>
