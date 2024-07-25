<script setup lang="ts">
import { watch } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import { store } from '@/store'

const route = useRoute()

watch(store, (newStore) => {
  localStorage.setItem("autorefresh_interval", newStore.autorefresh_interval.toString())
})
</script>

<template>
  <v-layout>
    <v-main>
      <RouterView v-slot="{ Component }">
        <Transition :name="route.meta.transition" mode="out-in">
          <Component :is="Component" />
        </Transition>
      </RouterView>
    </v-main>
  </v-layout>
</template>

<style>
.header {
  font-size: 3rem;
}

.slide-right-leave-active {
  transition: opacity 0.15s, transform 0.15s;
}

.slide-right-enter-from,
.slide-right-leave-to {
  opacity: 0;
  transform: translateX(60%);
}

.slide-left-leave-active {
  transition: opacity 0.15s, transform 0.15s;
}

.slide-left-enter-from,
.slide-left-leave-to {
  opacity: 0;
  transform: translateX(-60%);
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
