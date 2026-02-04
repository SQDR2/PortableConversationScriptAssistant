<script lang="ts" setup>
import MainLayout from './layouts/MainLayout.vue'
import { ref } from 'vue'
import SelectTargetDialog from './components/SelectTargetDialog.vue'

const showTargetDialog = ref(false)
</script>

<template>
  <router-view @open-target-dialog="showTargetDialog = true" />

  <SelectTargetDialog v-model="showTargetDialog" />
</template>

<style>
/* Global styles can go here or in style.css */
</style>
