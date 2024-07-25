<script setup lang="ts">
import { watch } from 'vue'
import { RouterView } from 'vue-router'
import { store } from '@/store'


watch(store, (newStore) => {
  localStorage.setItem("autorefresh_interval", newStore.autorefresh_interval.toString())
})
</script>

<template>
  <v-layout>
    <v-main>
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
.header {
  font-size: 3rem;
}

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
