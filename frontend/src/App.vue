<script lang="ts" setup>
import MainLayout from './layouts/MainLayout.vue'
import { ref, provide, onMounted } from 'vue'
import SelectTargetDialog from './components/SelectTargetDialog.vue'
import { utils } from '../wailsjs/go/models'

const showTargetDialog = ref(false)
const selectedTarget = ref<utils.WindowInfo | null>(null)

provide('selectedTarget', selectedTarget)

function onTargetSelected(win: utils.WindowInfo) {
  selectedTarget.value = win
  localStorage.setItem('sidekick.target', JSON.stringify(win))
}

onMounted(() => {
  const cached = localStorage.getItem('sidekick.target')
  if (cached) {
    try {
      selectedTarget.value = JSON.parse(cached) as utils.WindowInfo
    } catch (e) {
      localStorage.removeItem('sidekick.target')
    }
  }
})
</script>

<template>
  <router-view @open-target-dialog="showTargetDialog = true" />

  <SelectTargetDialog v-model="showTargetDialog" @target-selected="onTargetSelected" />
</template>

<style>
/* Global styles can go here or in style.css */
</style>
