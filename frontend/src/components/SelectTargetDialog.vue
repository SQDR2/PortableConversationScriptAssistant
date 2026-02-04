<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { GetStartApps, SetTarget } from '../../wailsjs/go/services/WindowService'
import { utils } from '../../wailsjs/go/models'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits(['update:modelValue', 'target-selected'])

const isOpen = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
})

const windows = ref<utils.WindowInfo[]>([])
const loading = ref(false)
const filter = ref('')

const filteredWindows = computed(() => {
  if (!filter.value) return windows.value
  const lowerFilter = filter.value.toLowerCase()
  return windows.value.filter(
    w => (w.title || '').toLowerCase().includes(lowerFilter) || (w.handle || '').includes(lowerFilter),
  )
})

async function fetchWindows() {
  loading.value = true
  try {
    windows.value = await GetStartApps()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function selectTarget(win: utils.WindowInfo) {
  try {
    await SetTarget(win.handle)
    emit('target-selected', win)
    isOpen.value = false
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  if (isOpen.value) {
    fetchWindows()
  }
})
</script>

<template>
  <q-dialog v-model="isOpen" @show="fetchWindows">
    <q-card style="min-width: 350px">
      <q-card-section>
        <div class="text-h6">选择目标应用程序</div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-input dense v-model="filter" placeholder="搜索..." autofocus />
      </q-card-section>

      <q-card-section class="q-pt-none scroll" style="max-height: 50vh">
        <div v-if="loading" class="row justify-center q-my-md">
          <q-spinner color="primary" size="2em" />
        </div>

        <q-list v-else separator bordered>
          <q-item v-for="win in filteredWindows" :key="win.handle" clickable v-ripple @click="selectTarget(win)">
            <q-item-section avatar>
              <q-icon name="window" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ win.title || '无标题窗口' }}</q-item-label>
              <q-item-label caption v-if="win.process">{{ win.process }}</q-item-label>
              <q-item-label caption>句柄: {{ win.handle }}</q-item-label>
            </q-item-section>
          </q-item>

          <q-item v-if="filteredWindows.length === 0">
            <q-item-section>
              <q-item-label class="text-grey italic">未找到窗口</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="取消" color="primary" v-close-popup />
        <q-btn flat label="刷新" color="secondary" @click="fetchWindows" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
