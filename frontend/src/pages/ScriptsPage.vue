<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import * as ScriptService from '../../wailsjs/go/services/ScriptService'
import * as CategoryService from '../../wailsjs/go/services/CategoryService'
import { models } from '../../wailsjs/go/models'
import { useQuasar } from 'quasar'
import CategoryList from '../components/CategoryList.vue'
import ScriptItem from '../components/ScriptItem.vue'

const $q = useQuasar()

// Data State
const scripts = ref<models.Script[]>([])
const categories = ref<models.Category[]>([])
const loading = ref(false)

// View State
const currentView = ref<'timeline' | 'directory'>('timeline')
const selectedCategoryId = ref<number | null>(null) // null for 'all' or 'uncategorized' context?
// Actually in Directory view:
// - If selectedCategoryId is null -> Show category list
// - If selectedCategoryId is set -> Show script list for that category
// But requirements say "Directory View: Show categories as cards/list... Click category to filter script list"
// We can implement a Split View: Left sidebar for categories, Right for scripts (or stacked on mobile)
// Given "Sidekick" is a small window, Tabs might be better.
// "Timeline" tab -> List of all scripts sorted by time
// "Directory" tab -> List of categories. Clicking one enters it (Breadcrumb navigation).

// Let's go with Breadcrumb/Drill-down for small screens.
// Directory View Root: List of Categories + "Uncategorized" folder.
// Clicking a folder -> Shows scripts in that folder.

const directoryPath = ref<models.Category | null>(null) // null = root

// Search
const searchText = ref('')

// Computed Scripts
const displayedScripts = computed(() => {
  try {
    const raw = scripts.value
    console.log('[DISPLAY_DEBUG] Start computing. scripts.value is:', Array.isArray(raw) ? raw.length : 'NOT_ARRAY')

    if (!raw || !Array.isArray(raw)) return []

    let list = [...raw]

    // 1. Search Filter (Global)
    const query = searchText.value?.trim().toLowerCase()
    if (query) {
      list = list.filter(s => {
        const content = (s.content || '').toLowerCase()
        const tags = (s.tags || '').toLowerCase()
        return content.includes(query) || tags.includes(query)
      })
    } else if (currentView.value === 'directory') {
      // 2. View Filter (Only if not searching)
      if (directoryPath.value) {
        const searchDirId = directoryPath.value.id
        list = list.filter(s => s.category_id === searchDirId)
      } else {
        list = list.filter(s => !s.category_id)
      }
    }

    console.log('[DISPLAY_DEBUG] End computing. Result count:', list.length)
    return list
  } catch (err) {
    console.error('[CRITICAL_COMPUTED_ERROR]', err)
    return []
  }
})

// Loading Data
async function loadData(showLoading = true) {
  console.log('loadData start, showLoading:', showLoading)
  if (showLoading) loading.value = true
  try {
    const p1 = ScriptService.ListScripts(1, 1000)
    const p2 = CategoryService.ListCategories()

    const [loadedScripts, loadedCats] = await Promise.all([p1, p2])
    console.log('loadData fetched. scripts:', loadedScripts?.length, 'cats:', loadedCats?.length)

    scripts.value = loadedScripts || []
    categories.value = loadedCats || []
    console.log('loadData state updated')
  } catch (e) {
    console.error('loadData error:', e)
    $q.notify({ type: 'negative', message: '数据加载失败', position: 'top' })
  } finally {
    if (showLoading) loading.value = false
    console.log('loadData finished, loading state:', loading.value)
  }
}

// Editor State
const showEditor = ref(false)
const editingScript = ref<models.Script | null>(null)
const editorContent = ref('')
const editorTags = ref('')
const editorCategoryId = ref<number | null>(null)
const editorImages = ref<string[]>([])
const uploadDummy = ref(null)

function openEditor(script?: models.Script) {
  if (script) {
    editingScript.value = script
    editorContent.value = script.content
    editorTags.value = script.tags
    editorCategoryId.value = script.category_id || null
    editorImages.value = script.images ? JSON.parse(script.images) : []
  } else {
    editingScript.value = null
    editorContent.value = ''
    editorTags.value = ''
    // Default category: current directory if in directory view
    if (currentView.value === 'directory' && directoryPath.value) {
      editorCategoryId.value = directoryPath.value.id
    } else {
      editorCategoryId.value = null
    }
    editorImages.value = []
  }
  showEditor.value = true
}

async function handleFileUpload(file: File) {
  if (!file) return
  const reader = new FileReader()
  reader.onload = async () => {
    const base64 = (reader.result as string).split(',')[1]
    const ext = file.name.substring(file.name.lastIndexOf('.'))
    try {
      // @ts-ignore
      const path = await ScriptService.SaveScriptImage(base64, ext)
      editorImages.value.push(path)
    } catch (e) {
      console.error(e)
      $q.notify({ type: 'negative', message: '图片上传失败' })
    }
  }
  reader.readAsDataURL(file)
}

function removeImage(index: number) {
  editorImages.value.splice(index, 1)
}

async function saveScript() {
  try {
    if (editingScript.value) {
      // @ts-ignore
      await ScriptService.UpdateScript(
        editingScript.value.id,
        editorContent.value,
        editorTags.value,
        editorCategoryId.value,
        JSON.stringify(editorImages.value),
      )
      $q.notify({ type: 'positive', message: '话术已更新' })
    } else {
      // @ts-ignore
      await ScriptService.CreateScript(
        editorContent.value,
        editorTags.value,
        editorCategoryId.value,
        JSON.stringify(editorImages.value),
      )
      $q.notify({ type: 'positive', message: '话术已创建' })
    }
    showEditor.value = false
    loadData(false)
  } catch (e) {
    $q.notify({ type: 'negative', message: '保存失败' })
  }
}

async function deleteScript(script: models.Script) {
  try {
    // @ts-ignore
    await ScriptService.DeleteScript(script.id)
    $q.notify({ type: 'positive', message: '话术已删除' })
    loadData(false)
  } catch (e) {
    $q.notify({ type: 'negative', message: '删除失败' })
  }
}

// Import (simplified from previous)
const showImportDialog = ref(false)
const importDelimiter = ref('\\n\\n')
const importFile = ref<File | null>(null)
const scriptsToImport = ref<string[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

function onFileSelected(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    // Read file... (Logic same as before, simplified for brevity in this rewrite)
    const reader = new FileReader()
    reader.onload = e => {
      const content = e.target?.result as string
      if (content) {
        scriptsToImport.value = content
          .split('\n\n')
          .map(s => s.trim())
          .filter(s => !!s)
        if (scriptsToImport.value.length > 0) showImportDialog.value = true
      }
    }
    reader.readAsText(file)
  }
  if (fileInput.value) fileInput.value.value = ''
}

// Category Navigation
function enterCategory(cat: models.Category) {
  directoryPath.value = cat
}
function goUp() {
  directoryPath.value = null
}

const categoryOptions = computed(() => {
  return [{ label: '未分类', value: null }, ...categories.value.map(c => ({ label: c.name, value: c.id }))]
})

onMounted(loadData)

function tweakPageHeight(offset: number) {
  return { height: `calc(100vh - ${offset}px)` }
}
</script>

<template>
  <q-page class="q-pa-md column no-wrap" style="height: 100%; min-height: unset; display: flex">
    <!-- Header: Search & Actions -->
    <div class="row q-mb-sm items-center q-gutter-sm">
      <q-input dense outlined v-model="searchText" placeholder="搜索..." class="col" debounce="300">
        <template v-slot:append>
          <q-icon name="search" />
        </template>
      </q-input>
      <q-btn color="primary" icon="add" round dense size="sm" @click="openEditor()">
        <q-tooltip>新建话术</q-tooltip>
      </q-btn>
    </div>

    <!-- View Switcher -->
    <div class="row q-mb-sm">
      <q-btn-group unelevated class="full-width">
        <q-btn
          :color="currentView === 'timeline' ? 'primary' : 'grey-3'"
          :text-color="currentView === 'timeline' ? 'white' : 'black'"
          label="时间轴"
          class="col"
          size="sm"
          @click="currentView = 'timeline'" />
        <q-btn
          :color="currentView === 'directory' ? 'primary' : 'grey-3'"
          :text-color="currentView === 'directory' ? 'white' : 'black'"
          label="目录视图"
          class="col"
          size="sm"
          @click="
            () => {
              currentView = 'directory'
              directoryPath = null
            }
          " />
      </q-btn-group>
    </div>

    <!-- Directory Navigation Breadcrumb -->
    <div
      v-if="currentView === 'directory' && directoryPath && !searchText"
      class="row items-center q-mb-sm bg-grey-2 text-grey-9 q-pa-xs rounded-borders">
      <q-btn flat dense icon="arrow_back" size="sm" @click="goUp" />
      <span class="text-caption q-ml-sm text-weight-bold">{{ directoryPath.name }}</span>
    </div>

    <!-- Main Content Area -->
    <div class="col relative-position scroll" style="flex: 1; overflow-y: auto; overflow-x: hidden">
      <div>
        <!-- Directory View: Category List (Only at Root) -->
        <div v-if="currentView === 'directory' && !directoryPath && !searchText">
          <CategoryList
            :categories="categories"
            :selectedId="null"
            @select="
              id => {
                const c = categories.find(x => x.id === id)
                if (c) enterCategory(c)
              }
            "
            @refresh="loadData(false)" />
          <q-separator class="q-my-sm" />
          <div class="text-caption text-grey-6 q-mb-xs q-px-sm">未分类话术</div>
        </div>

        <!-- Script List -->
        <div class="q-gutter-y-sm">
          <ScriptItem
            v-for="script in displayedScripts"
            :key="script.id"
            :script="script"
            @edit="openEditor"
            @delete="deleteScript" />

          <div v-if="displayedScripts.length === 0 && !loading" class="text-center text-grey q-pa-lg">
            {{ currentView === 'directory' && !directoryPath ? '无未分类话术' : '暂无话术' }}
          </div>
        </div>
      </div>

      <q-inner-loading :showing="loading" style="z-index: 9999">
        <q-spinner-tail color="primary" size="3em" />
        <div class="q-mt-md">加载中...</div>
      </q-inner-loading>
    </div>

    <!-- Editor Dialog -->
    <q-dialog v-model="showEditor" persistent>
      <q-card class="script-dialog-card" style="min-width: 300px">
        <q-card-section>
          <div class="text-h6">{{ editingScript ? '编辑话术' : '新建话术' }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-select
            v-model="editorCategoryId"
            :options="categoryOptions"
            label="所属目录"
            dense
            outlined
            emit-value
            map-options
            popup-content-class="script-dialog-card"
            class="q-mb-md" />
          <q-input
            v-model="editorContent"
            type="textarea"
            label="内容"
            filled
            autogrow
            autofocus
            style="max-height: 200px; overflow-y: auto" />
          <q-input v-model="editorTags" label="标签 (逗号分隔)" class="q-mt-md" outlined dense />

          <!-- Image Selection -->
          <div class="q-mt-lg">
            <div class="row items-center q-mb-sm">
              <span class="text-subtitle2 text-grey-8">图片附件</span>
              <q-space />
              <span class="text-caption text-grey-6">{{ editorImages.length }} / 10</span>
            </div>
            <div class="row q-gutter-md">
              <div v-for="(img, index) in editorImages" :key="index" class="relative-position">
                <q-img :src="img" style="width: 70px; height: 70px" class="rounded-borders shadow-1 border-grey" />
                <q-btn
                  round
                  dense
                  color="negative"
                  icon="close"
                  size="xs"
                  class="absolute-top-right"
                  style="top: -8px; right: -8px; z-index: 10"
                  @click="removeImage(index)" />
              </div>
              <q-file
                v-if="editorImages.length < 10"
                v-model="uploadDummy"
                borderless
                dense
                accept="image/*"
                display-value=""
                @update:model-value="handleFileUpload"
                style="width: 70px; height: 70px"
                class="bg-grey-2 rounded-borders overflow-hidden upload-box">
                <template v-slot:default>
                  <div class="full-width full-height flex flex-center">
                    <q-icon name="add_a_photo" color="grey-6" size="sm" />
                  </div>
                </template>
              </q-file>
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn label="保存" color="primary" @click="saveScript" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Import Dialog (Minimal placeholder to restore functional link) -->
    <q-dialog v-model="showImportDialog">
      <q-card class="script-dialog-card" style="min-width: 300px">
        <q-card-section><div class="text-h6">导入话术</div></q-card-section>
        <q-card-section>
          <div>找到 {{ scriptsToImport.length }} 条话术</div>
        </q-card-section>
        <q-card-actions align="right">
          <!-- @ts-ignore -->
          <q-btn
            label="确认导入"
            color="primary"
            @click="
              async () => {
                await ScriptService.ImportScripts(scriptsToImport)
                showImportDialog = false
                loadData()
              }
            " />
        </q-card-actions>
      </q-card>
    </q-dialog>
    <input type="file" ref="fileInput" style="display: none" @change="onFileSelected" accept=".txt" />
  </q-page>
</template>

<style lang="scss">
.script-dialog-card {
  color: rgba(0, 0, 0, 0.87);
  .q-field__label {
    color: rgba(0, 0, 0, 0.6);
  }
}

body.body--dark .script-dialog-card {
  color: white;
  .q-field__label {
    color: rgba(255, 255, 255, 0.7);
  }
}

.upload-box {
  border: 2px dashed #ddd;
  cursor: pointer;
  &:hover {
    border-color: var(--q-primary);
    background: #f5f5f5;
  }
}

body.body--dark .upload-box {
  border-color: #444;
  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }
}

.border-grey {
  border: 1px solid rgba(0, 0, 0, 0.1);
}
</style>
