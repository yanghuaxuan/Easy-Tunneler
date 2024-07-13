<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import { mdiWrenchOutline, mdiTunnelOutline } from '@mdi/js'
</script>

<template>
  <v-container class="d-flex justify-center w-screen h-screen container">
    <v-card class="pa-3 card w-75 h-55" variant="flat" rounded="lg">
      <v-row>
        <v-col cols="8">
          <h1 class="text-h3">Easy-Tunneler</h1>
        </v-col>
        <v-col class="d-flex justify-end">
          <v-btn-toggle color="var(--color-primary)"  class="toggle" rounded="lg" variant="outlined" divided mandatory>
            <v-btn @click="$router.replace({ name: 'tunnels' })"><v-icon size="25" :icon="mdiTunnelOutline" /></v-btn>
            <v-btn @click="$router.replace({ name: 'settings' })"><v-icon size="25" :icon="mdiWrenchOutline" /></v-btn>
          </v-btn-toggle>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <RouterView />
        </v-col>
      </v-row>
    </v-card>
  </v-container>
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
