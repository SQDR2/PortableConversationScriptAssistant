<script setup lang="ts">
import { ref } from 'vue'

const leftDrawerOpen = ref(false)

const emit = defineEmits(['open-target-dialog'])

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}
</script>

<template>
  <q-layout view="lHh Lpr lFf" class="glass-layout">
    <q-header elevated class="glass-header">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="菜单" @click="toggleLeftDrawer" />

        <q-toolbar-title> Sidekick 话术助手 </q-toolbar-title>

        <div>v0.0.1</div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="glass-drawer">
      <q-list>
        <q-item-label header> 工具 </q-item-label>

        <q-item clickable to="/">
          <q-item-section avatar>
            <q-icon name="assignment" />
          </q-item-section>
          <q-item-section>
            <q-item-label>话术管理</q-item-label>
            <q-item-label caption>管理您的常用话术</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable @click="emit('open-target-dialog')">
          <q-item-section avatar>
            <q-icon name="gps_fixed" />
          </q-item-section>
          <q-item-section>
            <q-item-label>选择目标</q-item-label>
            <q-item-label caption>关联目标应用程序</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<style lang="scss">
.glass-layout {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.glass-header {
  background: rgba(25, 118, 210, 0.8) !important;
  backdrop-filter: blur(5px);
}

.glass-drawer {
  background: rgba(255, 255, 255, 0.5) !important;
  backdrop-filter: blur(10px);
}

body.body--dark {
  .glass-layout {
    background: rgba(0, 0, 0, 0.2);
  }
  .glass-drawer {
    background: rgba(30, 30, 30, 0.7) !important;
  }
}
</style>
