<script setup lang="ts">
import { ref, inject, computed, type Ref, onMounted } from 'vue'
import { utils } from '../../wailsjs/go/models'
import TranslationBar from '../components/TranslationBar.vue'
import ConfigDialog from '../components/ConfigDialog.vue'
import { GetConfigStatus } from '../../wailsjs/go/services/TranslationService'

const leftDrawerOpen = ref(false)

const selectedTarget = inject<Ref<utils.WindowInfo | null>>('selectedTarget', ref(null))

const targetTitle = computed(() => selectedTarget.value?.title || '未关联应用')
const targetProcess = computed(() => selectedTarget.value?.process || '')
const targetHandle = computed(() => selectedTarget.value?.handle || '')

const emit = defineEmits(['open-target-dialog'])
const showTranslation = ref(true)
const configDialogOpen = ref(false)
const isConfigForced = ref(false)

onMounted(async () => {
  try {
    const isConfigured = await GetConfigStatus()
    if (!isConfigured) {
      isConfigForced.value = true
      configDialogOpen.value = true
    }
  } catch (err) {
    console.error('Failed to check config status:', err)
  }
})

function openConfigDialog() {
  isConfigForced.value = false
  configDialogOpen.value = true
}

function toggleTranslation() {
  showTranslation.value = !showTranslation.value
}

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

        <q-btn
          dense
          flat
          round
          icon="translate"
          @click="toggleTranslation"
          :color="showTranslation ? 'primary' : 'grey'">
          <q-tooltip>切换翻译栏</q-tooltip>
        </q-btn>
        <div class="q-ml-sm">v0.0.1</div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="glass-drawer">
      <q-list class="drawer-list">
        <q-item-label header class="drawer-header">Sidekick</q-item-label>

        <q-item class="drawer-card">
          <q-item-section avatar>
            <q-avatar color="primary" text-color="white" icon="link" />
          </q-item-section>
          <q-item-section class="drawer-card__content">
            <div class="drawer-card__title">
              <span>当前关联</span>
              <q-btn dense flat color="primary" label="更换" @click="emit('open-target-dialog')" />
            </div>
            <q-item-label class="drawer-card__name" lines="1">{{ targetTitle }}</q-item-label>
            <q-item-label caption lines="1" v-if="targetProcess">{{ targetProcess }}</q-item-label>
            <q-item-label caption lines="1" v-if="targetHandle">句柄: {{ targetHandle }}</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable to="/" class="drawer-item">
          <q-item-section avatar>
            <q-icon name="assignment" />
          </q-item-section>
          <q-item-section>
            <q-item-label>话术管理</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable @click="openConfigDialog" class="drawer-item">
          <q-item-section avatar>
            <q-icon name="settings" />
          </q-item-section>
          <q-item-section>
            <q-item-label>翻译配置</q-item-label>
            <q-item-label caption>设置腾讯云 API 密钥</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable @click="emit('open-target-dialog')" class="drawer-item">
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

    <q-footer v-if="showTranslation" elevated class="glass-footer">
      <TranslationBar />
    </q-footer>

    <ConfigDialog v-model="configDialogOpen" :forced="isConfigForced" />
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
  background: rgba(255, 255, 255, 0.9) !important;
  backdrop-filter: blur(16px);
  border-right: 1px solid rgba(0, 0, 0, 0.06);
  color: black;
}

.drawer-list {
  padding: 12px 10px 16px;
}

.drawer-header {
  color: rgba(0, 0, 0, 0.55);
  font-weight: 600;
  letter-spacing: 0.5px;
  margin: 4px 6px 8px;
}

.drawer-card {
  margin: 6px 6px 18px;
  padding: 10px 12px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(25, 118, 210, 0.08), rgba(25, 118, 210, 0.02));
  border: 1px solid rgba(25, 118, 210, 0.15);
}

.drawer-card__content {
  gap: 2px;
}

.drawer-card__title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.55);
}

.drawer-card__name {
  font-weight: 600;
  margin-top: 2px;
}

.drawer-item {
  margin: 6px 6px;
  border-radius: 10px;
}

.drawer-item:hover {
  background: rgba(25, 118, 210, 0.08);
}

body.body--dark {
  .glass-layout {
    background: rgba(0, 0, 0, 0.2);
  }
  .glass-drawer {
    background: rgba(20, 20, 20, 0.9) !important;
    border-right: 1px solid rgba(255, 255, 255, 0.08);
    color: white;
  }
  .drawer-header {
    color: rgba(255, 255, 255, 0.65);
  }
  .drawer-card__title {
    color: rgba(255, 255, 255, 0.6);
  }
  .drawer-card {
    background: linear-gradient(135deg, rgba(66, 165, 245, 0.18), rgba(66, 165, 245, 0.05));
    border: 1px solid rgba(66, 165, 245, 0.2);
  }
  .drawer-item:hover {
    background: rgba(66, 165, 245, 0.14);
  }
}
</style>

<style scoped>
.glass-footer {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  color: black;
}

body.body--dark .glass-footer {
  background: rgba(30, 30, 30, 0.95);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  color: white;
}
</style>
