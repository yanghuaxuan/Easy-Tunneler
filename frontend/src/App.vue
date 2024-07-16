<script setup lang="ts">
import { RouterView } from 'vue-router'
import { mdiWrenchOutline, mdiTunnelOutline, mdiPlus } from '@mdi/js'
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay();
</script>

<template>
    <v-layout>
      <v-main>
        <RouterView />
      </v-main>
    </v-layout>
</template>

<style scoped>
/* body {
  background-color: #FEF9EB
} */
 .container {
  height: dvh;
 }
 .card {
  background-color: var(--color-surface-container);
  color: var(--color-on-surface);
 }
 .toggle {
  border-color: var(--color-outline)
 }
</style>
