<script setup lang="ts">
import { defineModel } from 'vue'

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
    }

    /* The slider */
    .slider {
      position: absolute;
      cursor: pointer;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: #ccc;
      transition: 0.4s;
      border-radius: 34px;
    }

    .slider:before {
      position: absolute;
      content: "";
      height: 26px;
      width: 26px;
      left: 4px;
      bottom: 4px;
      background-color: white;
      transition: 0.4s;
      border-radius: 50%;
    }

    input:checked + .slider {
      background-color: #2196F3;
    }

    input:checked + .slider:before {
      transform: translateX(26px);
    }

</style>