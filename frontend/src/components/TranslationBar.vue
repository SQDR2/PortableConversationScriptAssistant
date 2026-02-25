<script setup lang="ts">
import { ref } from 'vue'
import * as TranslationService from '../../wailsjs/go/services/TranslationService'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State
const sourceText = ref('')
const targetText = ref('')
const loading = ref(false)

// Languages
const languages = [
  { label: '自动识别', value: 'auto' },
  { label: '简体中文', value: 'zh' },
  { label: '英语', value: 'en' },
  { label: '日语', value: 'ja' },
  { label: '韩语', value: 'ko' },
  { label: '法语', value: 'fr' },
  { label: '西班牙语', value: 'es' },
  { label: '德语', value: 'de' },
  { label: '俄语', value: 'ru' },
]

const sourceLang = ref('zh')
const targetLang = ref('en')
const showTranslatorPanel = ref(false)

// Debounce timer
let debounceTimer: number | null = null

function onInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  // Disabled auto-translate to save quota
  // debounceTimer = window.setTimeout(doTranslate, 800)
}

function onKeyEnter(e: KeyboardEvent) {
  if (!e.shiftKey) {
    e.preventDefault()
    doTranslate()
  }
}

async function doTranslate() {
  if (debounceTimer) clearTimeout(debounceTimer)
  if (!sourceText.value.trim()) {
    targetText.value = ''
    return
  }

  loading.value = true
  try {
    // @ts-ignore
    const result = await TranslationService.Translate(sourceText.value, sourceLang.value, targetLang.value)
    targetText.value = result
  } catch (e: any) {
    console.error(e)
    $q.notify({
      type: 'negative',
      message: '翻译失败: ' + (e.message || e),
      position: 'top',
      timeout: 2000,
    })
  } finally {
    loading.value = false
  }
}

function copyResult() {
  if (!targetText.value) return
  navigator.clipboard.writeText(targetText.value)
  $q.notify({
    type: 'positive',
    message: '已复制',
    position: 'center',
    timeout: 500,
  })
}

function swapLanguages() {
  const temp = sourceLang.value
  sourceLang.value = targetLang.value
  targetLang.value = temp

  // Also swap text if meaningful
  if (targetText.value) {
    sourceText.value = targetText.value
    targetText.value = ''
    doTranslate()
  }
}

function clearInput() {
  sourceText.value = ''
  targetText.value = ''
}
</script>

<template>
  <div class="translator-entry-wrap q-pa-md">
    <q-slide-transition>
      <div v-show="showTranslatorPanel" class="translator-panel q-mb-sm">
        <div class="row items-center justify-between q-mb-xs">
          <div class="row items-center no-wrap col-grow">
            <q-select
              v-model="sourceLang"
              :options="languages"
              dense
              borderless
              options-dense
              emit-value
              map-options
              class="col text-caption lang-select"
              popup-content-class="text-caption lang-menu-popup" />
            <q-btn flat round dense icon="swap_horiz" size="sm" color="grey-7" @click="swapLanguages" />
            <q-select
              v-model="targetLang"
              :options="languages"
              dense
              borderless
              options-dense
              emit-value
              map-options
              class="col text-caption lang-select target-lang-select"
              popup-content-class="text-caption lang-menu-popup" />
          </div>
          <q-btn
            flat
            round
            dense
            :icon="loading ? 'hourglass_empty' : 'translate'"
            size="sm"
            color="primary"
            :loading="loading"
            @click="doTranslate" />
        </div>

        <q-input
          v-model="sourceText"
          dense
          outlined
          placeholder="输入..."
          autogrow
          class="q-mb-xs"
          bg-color="white"
          @keydown.enter="onKeyEnter"
          debounce="0"
          @update:model-value="onInput">
          <template v-slot:append v-if="sourceText">
            <q-icon name="close" @click="clearInput" class="cursor-pointer" />
          </template>
        </q-input>

        <div
          class="output-area q-pa-sm relative-position rounded-borders cursor-pointer"
          :class="{ 'has-content': !!targetText }"
          @click="copyResult"
          v-ripple>
          <div v-if="!targetText" class="text-grey-5 text-italic text-caption">翻译结果...</div>
          <div v-else class="text-body2 output-content">{{ targetText }}</div>

          <q-icon
            v-if="targetText"
            name="content_copy"
            size="xs"
            class="absolute-bottom-right q-ma-xs text-grey-6" />
        </div>
      </div>
    </q-slide-transition>

    <q-btn
      no-caps
      unelevated
      class="translator-entry-btn"
      icon="translate"
      label="Open Translator"
      @click="showTranslatorPanel = !showTranslatorPanel" />
  </div>
</template>

<style scoped>
.translator-entry-wrap {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.translator-entry-btn {
  align-self: center;
  background: #1f2735;
  color: #ffffff;
  border-radius: 24px;
  min-height: 44px;
  min-width: 220px;
  padding: 0 20px;
  box-shadow: 0 8px 20px rgba(20, 27, 39, 0.24);
}

.translator-panel {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e6eaf2;
  box-shadow: 0 10px 24px rgba(19, 27, 39, 0.12);
  padding: 12px;
  max-width: 720px;
  width: 100%;
  align-self: center;
}

.output-content {
  white-space: pre-wrap;
  word-break: break-word;
  min-height: 1.5em; /* Ensure consistent line height */
}

.lang-select :deep(.q-field__control) {
  min-height: 24px;
}
.lang-select :deep(.q-field__native) {
  padding-top: 0;
  padding-bottom: 0;
  font-size: 12px;
}
.target-lang-select :deep(.q-field__native) {
  justify-content: flex-end;
  text-align: right;
}

.output-area {
  min-height: 38px;
  width: 100%; /* Prevent shrinking in flex container */
  background: rgba(0, 0, 0, 0.03);
  border: 1px dashed rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
}

.output-area.has-content {
  background: rgba(25, 118, 210, 0.05);
  border: 1px solid rgba(25, 118, 210, 0.1);
}

.output-area:hover.has-content {
  background: rgba(25, 118, 210, 0.1);
}

/* Fix dropdown text color visibility */
.lang-select :deep(.q-field__native) {
  padding-top: 0;
  padding-bottom: 0;
  font-size: 12px;
  color: #1d1d1d; /* Force dark text for visibility against white bg */
}

</style>

<!-- Global styles for popup -->
<style>
.lang-menu-popup .q-item {
  color: #1d1d1d !important; /* Force visible text in dropdown */
}
.lang-menu-popup .q-item--active {
  color: #1976d2 !important; /* Primary color for active */
  font-weight: bold;
}
</style>
