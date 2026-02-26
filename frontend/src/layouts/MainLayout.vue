<script setup lang="ts">
import { ref, inject, type Ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { utils } from '../../wailsjs/go/models'
import TranslationBar from '../components/TranslationBar.vue'
import ConfigDialog from '../components/ConfigDialog.vue'
import { GetConfigStatus } from '../../wailsjs/go/services/TranslationService'
import { GetVersion } from '../../wailsjs/go/main/App'

const router = useRouter()
const route = useRoute()

const selectedTarget = inject<Ref<utils.WindowInfo | null>>('selectedTarget', ref(null))
const appName = 'Script Assistant'
const headerTitle = computed(() => selectedTarget.value?.title?.trim() || appName)

const emit = defineEmits(['open-target-dialog'])
const configDialogOpen = ref(false)
const isConfigForced = ref(false)
const appVersion = ref('v2.0.4')
const settingsDialogOpen = ref(false)

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

  try {
    appVersion.value = await GetVersion()
  } catch (err) {
    console.error('Failed to get version:', err)
  }
})

function openConfigDialog() {
  settingsDialogOpen.value = false
  isConfigForced.value = false
  configDialogOpen.value = true
}

function openTargetDialog() {
  settingsDialogOpen.value = false
  emit('open-target-dialog')
}

function openSettingsDialog() {
  settingsDialogOpen.value = true
}

function openCreateScript() {
  const query = { ...route.query, new: String(Date.now()) }
  router.replace({ path: '/', query })
}
</script>

<template>
  <q-layout view="lHh lpr lFf" class="app-layout">
    <q-header class="app-header" bordered>
      <q-toolbar class="header-toolbar">
        <div class="row items-center no-wrap q-gutter-sm">
          <q-icon name="chat_bubble_outline" size="20px" class="header-icon" />
          <q-toolbar-title class="header-title">{{ headerTitle }}</q-toolbar-title>
        </div>

        <q-space />

        <q-btn dense flat round icon="add" class="header-action" @click="openCreateScript">
          <q-tooltip>添加话术</q-tooltip>
        </q-btn>
        <q-btn dense flat round icon="settings" class="header-action" @click="openSettingsDialog">
          <q-tooltip>设置</q-tooltip>
        </q-btn>

        <div class="header-version q-ml-sm">{{ appVersion }}</div>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer class="translator-entry-footer">
      <TranslationBar />
    </q-footer>

    <q-dialog v-model="settingsDialogOpen">
      <q-card class="settings-dialog-card" style="width: 360px; max-width: 92vw">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-subtitle1 text-weight-bold">设置</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <q-list bordered separator class="rounded-borders">
            <q-item clickable v-ripple @click="openTargetDialog">
              <q-item-section avatar><q-icon name="link" /></q-item-section>
              <q-item-section>
                <q-item-label>{{ selectedTarget ? '切换关联应用' : '关联应用' }}</q-item-label>
                <q-item-label caption>{{ selectedTarget?.title ? `当前：${selectedTarget.title}` : '选择需要跟随的目标窗口' }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item clickable v-ripple @click="openConfigDialog">
              <q-item-section avatar><q-icon name="vpn_key" /></q-item-section>
              <q-item-section>
                <q-item-label>腾讯翻译配置</q-item-label>
                <q-item-label caption>设置 SecretId 与 SecretKey</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </q-dialog>

    <ConfigDialog v-model="configDialogOpen" :forced="isConfigForced" />
  </q-layout>
</template>

<style lang="scss">
.app-layout {
  background: #f6f7fb;
}

.app-header {
  background: #f6f7fb;
  border-bottom: 1px solid #e8eaf0;
}

.header-toolbar {
  padding: 8px 12px;
}

.header-icon {
  color: #4b5565;
}

.header-title {
  font-size: 16px;
  font-weight: 700;
  color: #202733;
}

.header-action {
  color: #576074;
}

.header-version {
  color: #6f7787;
  font-size: 12px;
}

.settings-dialog-card {
  border-radius: 12px;
}

body.body--dark {
  .app-layout,
  .app-header {
    background: #f6f7fb;
  }
}
</style>

<style scoped>
.translator-entry-footer {
  background: transparent;
  border-top: none;
}
</style>
