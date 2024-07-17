<script setup lang="ts">
import { defineModel } from 'vue'
const props = defineProps<{
    color?: string
}>()

const checked = defineModel({ default: false })
</script>

<template>
    <!-- <input type="checkbox" class="tunnel-switch" @click.stop="console.log('xd')"> -->
    <label class="switch">
    <input v-model="checked" type="checkbox" class="no-checkbox">
    <span class="slider"></span>
  </label>
</template>

<style scoped>
    .no-checkbox {
      display: none;
    }

    /* The switch container */
    .switch {
      position: relative;
      display: inline-block;
      width: 60px;
      height: 34px;
      border-top: 500px;
    }

    /* The slider */
    .slider {
      position: absolute;
      cursor: pointer;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      /* background-color: var(--color-outline); */
      background-color: var(--color-surface-container-highest);
      transition: 0.4s;
      border-radius: 34px;
      border: 1px solid;
      border-color: var(--color-outline);
    }

    .slider:before {
      position: absolute;
      content: "";
      height: 26px;
      width: 26px;
      left: 4px;
      bottom: 4px;
      background-color: var(--color-outline);
      transition: 0.4s;
      border-radius: 50%;
    }

    input:checked + .slider {
        background-color: var(--color-primary)
    }

    input:checked + .slider:before {
      transform: translateX(26px);
      background-color: var(--color-on-primary);
    }

</style>