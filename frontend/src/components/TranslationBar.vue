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
  { label: '中文', value: 'zh' },
  { label: 'English', value: 'en' },
  { label: '日本語', value: 'ja' },
  { label: '한국어', value: 'ko' },
  { label: 'Français', value: 'fr' },
  { label: 'Español', value: 'es' },
  { label: 'Deutsch', value: 'de' },
  { label: 'Русский', value: 'ru' },
  { label: 'Português', value: 'pt' },
  { label: 'Italiano', value: 'it' },
  { label: 'Türkçe', value: 'tr' },
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
  <div class="translator-entry-wrap q-pa-sm">
    <q-slide-transition>
      <div v-show="showTranslatorPanel" class="translator-panel q-mb-sm">
        <div class="row items-center q-mb-sm translator-header">
          <div class="text-subtitle1 text-weight-bold translator-title">Translator</div>
          <q-space />
          <q-btn flat round dense icon="close" color="grey-6" @click="showTranslatorPanel = false" />
        </div>

        <div class="row items-center no-wrap q-mb-sm lang-switch-wrap">
          <q-select
            v-model="sourceLang"
            :options="languages"
            emit-value
            map-options
            dense
            borderless
            options-dense
            class="lang-select" />
          <q-btn flat round dense icon="swap_horiz" size="sm" color="grey-6" class="q-mx-xs" @click="swapLanguages" />
          <q-select
            v-model="targetLang"
            :options="languages"
            emit-value
            map-options
            dense
            borderless
            options-dense
            class="lang-select" />
        </div>

        <q-input
          v-model="sourceText"
          outlined
          type="textarea"
          autogrow
          placeholder="Enter Chinese text..."
          class="q-mb-sm"
          bg-color="white"
          @keydown.enter="onKeyEnter"
          debounce="0"
          @update:model-value="onInput">
          <template v-slot:append v-if="sourceText">
            <q-icon name="close" @click="clearInput" class="cursor-pointer" />
          </template>
        </q-input>

        <q-btn
          unelevated
          no-caps
          class="translate-btn q-mb-sm"
          :loading="loading"
          label="Translate"
          @click="doTranslate" />

        <div class="output-area" :class="{ 'has-content': !!targetText }" @click="copyResult">
          <div v-if="!targetText" class="text-grey-6">Translation will appear here...</div>
          <div v-else class="output-text">{{ targetText }}</div>
        </div>
      </div>
    </q-slide-transition>

    <q-btn
      v-if="!showTranslatorPanel"
      no-caps
      unelevated
      class="translator-entry-btn"
      icon="translate"
      label="Open Translator"
      @click="showTranslatorPanel = true" />
  </div>
</template>

<style scoped>
.translator-entry-wrap {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.translator-entry-btn {
  align-self: stretch;
  background: #182136;
  color: #ffffff;
  border-radius: 10px;
  min-height: 44px;
  width: 100%;
  padding: 0 16px;
  box-shadow: 0 4px 10px rgba(20, 27, 39, 0.18);
}

.translator-panel {
  background: #f7f8fc;
  border-radius: 12px;
  border: 1px solid #e0e5ef;
  box-shadow: 0 8px 18px rgba(21, 31, 47, 0.12);
  padding: 12px;
  width: 100%;
}

.translator-header {
  border-bottom: 1px solid #e4e9f2;
  padding-bottom: 6px;
}

.translator-title {
  color: #364156;
}

.lang-switch-wrap {
  background: #f4f6fa;
  border: 1px solid #d9dfeb;
  border-radius: 8px;
  padding: 4px;
}

.lang-select {
  flex: 1;
}

.lang-select :deep(.q-field__control) {
  min-height: 28px;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e2e7f1;
  padding: 0 6px;
}

.lang-select :deep(.q-field__native),
.lang-select :deep(.q-field__input) {
  color: #4b5568;
  font-size: 13px;
  padding: 0;
}

.translate-btn {
  width: 100%;
  background: #8f87ea;
  color: #ffffff;
  border-radius: 8px;
  min-height: 40px;
  font-weight: 600;
}

.output-area {
  min-height: 62px;
  border: 1px dashed #c7ceda;
  background: #fafbfe;
  border-radius: 8px;
  padding: 10px 12px;
  cursor: pointer;
}

.output-area.has-content {
  border-style: solid;
  border-color: #c8d2ea;
  background: #f3f7ff;
}

.output-text {
  white-space: pre-wrap;
  word-break: break-word;
}

</style>
