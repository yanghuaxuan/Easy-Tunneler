<script setup lang="ts">
import { mdiArrowLeft, mdiWrenchOutline } from "@mdi/js"
import { RouterView } from 'vue-router'
import { useRoute } from 'vue-router'


const route = useRoute()
</script>

<template>
  <v-layout>
    <v-main>
      <v-row v-if="route.name === 'tunnels'">
      <v-col>
        <h1 class="text-h3">Tunnels</h1>
      </v-col>
      <v-col class="d-flex justify-end">
        <v-btn @click="$router.replace({ name: 'settings' })" class="rounded-button" size="x-large" variant="text"
          :icon="mdiWrenchOutline" />
      </v-col>
      </v-row>
      <v-row v-else class="d-flex align-center">
        <v-col>
          <v-btn @click="$router.replace({ name: 'tunnels' })" class="rounded-button" size="x-large" variant="text"
            :icon="mdiArrowLeft" />
        </v-col>
        <v-col class="d-flex justify-end">
          <h1 class="text-h3">Settings</h1>
        </v-col>
      </v-row>
      <!-- <RouterView v-slot="{ Component }">
        <Transition name="slider" mode="out-in">
          <Component :is="Component" :key="$route.path" />
        </Transition>
      </RouterView> -->
      <RouterView />
    </v-main>
  </v-layout>
</template>

<style>
.slider-enter-active,
.slider-leave-active {
  transition: transform 0.25s;
}

.slider-enter-from,
.slider-leave-to {
  opacity: 0;
  transform: translateX(-30%);
}

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
