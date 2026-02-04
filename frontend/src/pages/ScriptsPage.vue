<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as ScriptService from '../../wailsjs/go/services/ScriptService'
import { models } from '../../wailsjs/go/models'
import { useQuasar } from 'quasar'

const $q = useQuasar()

const scripts = ref<models.Script[]>([])
const loading = ref(false)
const searchText = ref('')
const page = ref(1)
const pageSize = 50
const hasMore = ref(true)

// Editor state
const showEditor = ref(false)
const editingScript = ref<models.Script | null>(null)
const editorContent = ref('')
const editorTags = ref('')

async function loadScripts(reset = false) {
  if (reset) {
    page.value = 1
    scripts.value = []
    hasMore.value = true
  }

  if (!hasMore.value && !reset) return

  loading.value = true
  try {
    let newScripts: models.Script[] = []
    if (searchText.value) {
      // Search mode - currently fetch all matches (limit 50 in backend)
      // @ts-ignore
      newScripts = await ScriptService.SearchScripts(searchText.value)
      hasMore.value = false // Search API doesn't support pagination yet in this V1
    } else {
      // List mode
      // @ts-ignore
      newScripts = await ScriptService.ListScripts(page.value, pageSize)
      if (newScripts.length < pageSize) {
        hasMore.value = false
      }
      page.value++
    }

    if (reset) {
      scripts.value = newScripts
    } else {
      scripts.value.push(...newScripts)
    }
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: '加载话术失败' })
  } finally {
    loading.value = false
  }
}

function onScroll(index: number, done: () => void) {
  if (!searchText.value) {
    // Only paginate in list mode
    loadScripts().then(() => done())
  } else {
    done()
  }
}

function openEditor(script?: models.Script) {
  if (script) {
    editingScript.value = script
    editorContent.value = script.content
    editorTags.value = script.tags
  } else {
    editingScript.value = null
    editorContent.value = ''
    editorTags.value = ''
  }
  showEditor.value = true
}

async function saveScript() {
  try {
    if (editingScript.value) {
      // @ts-ignore
      await ScriptService.UpdateScript(editingScript.value.id, editorContent.value, editorTags.value)
      $q.notify({ type: 'positive', message: '话术已更新' })
    } else {
      // @ts-ignore
      await ScriptService.CreateScript(editorContent.value, editorTags.value)
      $q.notify({ type: 'positive', message: '话术已创建' })
    }
    showEditor.value = false
    loadScripts(true)
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: '保存话术失败' })
  }
}

function confirmDelete(script: models.Script) {
  $q.dialog({
    title: '确认',
    message: '确定要删除这条话术吗？',
    cancel: '取消',
    ok: '确定',
    persistent: true,
  }).onOk(async () => {
    try {
      // @ts-ignore
      await ScriptService.DeleteScript(script.id)
      $q.notify({ type: 'positive', message: '话术已删除' })
      loadScripts(true)
    } catch (e) {
      $q.notify({ type: 'negative', message: '删除失败' })
    }
  })
}

// Import state
const showImportDialog = ref(false)
const importDelimiter = ref('\\n\\n')
const importCount = ref(0)
const importing = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
let scriptsToImport: string[] = []

function onFileSelected(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = e => {
    const content = e.target?.result as string
    if (!content) return

    // Parse regex delimiter
    let delimiter = importDelimiter.value
    // If it looks like a regex (starts/ends with /), parse it? No, just string or simplified regex.
    // For now, support string delimiter or simple newline handling.
    // Actually JS split supports regex string if passed to new RegExp

    let separator: string | RegExp = delimiter
    try {
      // Handle escaped newlines
      const unescaped = delimiter.replace(/\\n/g, '\n').replace(/\\r/g, '\r')
      separator = unescaped
    } catch (e) {}

    scriptsToImport = content
      .split(separator)
      .map(s => s.trim())
      .filter(s => s.length > 0)
    importCount.value = scriptsToImport.length

    if (importCount.value > 0) {
      showImportDialog.value = true
    } else {
      $q.notify({ type: 'warning', message: '使用当前分隔符未找到话术' })
    }

    // Reset input
    if (fileInput.value) fileInput.value.value = ''
  }
  reader.readAsText(file)
}

async function confirmImport() {
  importing.value = true
  try {
    // @ts-ignore
    await ScriptService.ImportScripts(scriptsToImport)
    $q.notify({ type: 'positive', message: `成功导入 ${scriptsToImport.length} 条话术` })
    showImportDialog.value = false
    loadScripts(true)
  } catch (e) {
    $q.notify({ type: 'negative', message: '导入失败' })
  } finally {
    importing.value = false
  }
}

function triggerImport() {
  fileInput.value?.click()
}

// Search debounce could be added, for now enter to search or lazy update
// We'll use @update:model-value with debounce in template

onMounted(() => {
  loadScripts(true)
})
</script>

<template>
  <q-page class="q-pa-md">
    <div class="row q-mb-md">
      <q-input
        dense
        outlined
        v-model="searchText"
        placeholder="搜索话术..."
        class="col"
        debounce="300"
        @update:model-value="loadScripts(true)">
        <template v-slot:append>
          <q-icon name="search" />
        </template>
      </q-input>
      <q-btn color="primary" icon="add" label="新建" class="q-ml-md" @click="openEditor()" />
      <q-btn color="secondary" icon="upload" label="导入" class="q-ml-md" @click="triggerImport" />
    </div>

    <!-- Virtual Scroll for performance with large lists -->
    <q-virtual-scroll style="height: calc(100vh - 150px)" :items="scripts" separator v-slot="{ item, index }">
      <q-item :key="item.id" clickable v-ripple @click="openEditor(item)">
        <q-item-section>
          <q-item-label class="text-body1">{{ item.content }}</q-item-label>
          <q-item-label caption lines="1">
            <q-chip v-if="item.tags" size="xs" color="secondary" text-color="white">{{ item.tags }}</q-chip>
            {{ new Date(item.created_at).toLocaleString() }}
          </q-item-label>
        </q-item-section>

        <q-item-section side>
          <q-btn flat round color="negative" icon="delete" @click.stop="confirmDelete(item)" />
        </q-item-section>
      </q-item>
    </q-virtual-scroll>

    <!-- 编辑对话框 -->
    <q-dialog v-model="showEditor" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editingScript ? '编辑话术' : '新建话术' }}</div>
        </q-card-section>

        <q-card-section>
          <q-input v-model="editorContent" type="textarea" label="内容" filled autogrow autofocus />
          <q-input v-model="editorTags" label="标签 (逗号分隔)" class="q-mt-md" outlined />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn label="保存" color="primary" @click="saveScript" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 导入对话框 -->
    <q-dialog v-model="showImportDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">导入话术</div>
        </q-card-section>

        <q-card-section>
          <q-input v-model="importDelimiter" label="分隔符 (支持正则)" hint="例如：\\n\\n 表示双换行符" />

          <div class="q-mt-md text-subtitle2">预览：找到 {{ importCount }} 条话术。</div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn label="导入" color="primary" @click="confirmImport" :loading="importing" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Hidden File Input -->
    <input type="file" ref="fileInput" style="display: none" @change="onFileSelected" accept=".txt" />
  </q-page>
</template>
