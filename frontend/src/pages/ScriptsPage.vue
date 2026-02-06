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
  let list = scripts.value

  // 1. Search Filter (Global)
  if (searchText.value) {
    const lower = searchText.value.toLowerCase()
    list = list.filter(s => s.content.toLowerCase().includes(lower) || s.tags?.toLowerCase().includes(lower))
  }

  // 2. View Filter
  if (currentView.value === 'directory' && !searchText.value) {
    // If searching, we show all matches regardless of directory?
    // Requirement decision: "Search should be global".
    // So if searchText is present, we ignore directory filter? Or we highlight?
    // Let's make search override view filters for simplicity.

    if (directoryPath.value) {
      // Inside a category
      list = list.filter(s => s.category_id === directoryPath.value?.id)
    } else {
      // Root of Directory view: Don't show any scripts? Or show 'Uncategorized'?
      // Design: "First layer: Show all 'directory cards' + 'Uncategorized scripts'?"
      // Let's show only Uncategorized scripts at root, or nothing?
      // Usually file managers show folders + files.
      // We will handle this in template.
      list = list.filter(s => !s.category_id) // Root shows uncategorized
    }
  }

  // 3. Sort: Always Created Desc for now
  // (Already sorted by backend mostly, but client filtering might mess up if we append?)
  // Re-sort to be safe
  return list.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

// Loading Data
async function loadData() {
  loading.value = true
  try {
    // Load all data parallel
    // @ts-ignore
    const p1 = ScriptService.ListScripts(1, 10000) // Load all for client-side filtering
    // @ts-ignore
    const p2 = CategoryService.ListCategories()

    const [loadedScripts, loadedCats] = await Promise.all([p1, p2])
    scripts.value = loadedScripts
    categories.value = loadedCats
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: '数据加载失败' })
  } finally {
    loading.value = false
  }
}

// Editor State
const showEditor = ref(false)
const editingScript = ref<models.Script | null>(null)
const editorContent = ref('')
const editorTags = ref('')
const editorCategoryId = ref<number | null>(null)

function openEditor(script?: models.Script) {
  if (script) {
    editingScript.value = script
    editorContent.value = script.content
    editorTags.value = script.tags
    editorCategoryId.value = script.category_id || null
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
  }
  showEditor.value = true
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
      )
      $q.notify({ type: 'positive', message: '话术已更新' })
    } else {
      // @ts-ignore
      await ScriptService.CreateScript(editorContent.value, editorTags.value, editorCategoryId.value)
      $q.notify({ type: 'positive', message: '话术已创建' })
    }
    showEditor.value = false
    loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: '保存失败' })
  }
}

async function deleteScript(script: models.Script) {
  try {
    // @ts-ignore
    await ScriptService.DeleteScript(script.id)
    $q.notify({ type: 'positive', message: '话术已删除' })
    loadData()
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
  <q-page class="q-pa-md column no-wrap" :style-fn="tweakPageHeight">
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
    <q-scroll-area class="col">
      <!-- Loading -->
      <div v-if="loading" class="row justify-center q-pa-md">
        <q-spinner color="primary" size="2em" />
      </div>

      <template v-else>
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
            @refresh="loadData" />
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

          <div v-if="displayedScripts.length === 0" class="text-center text-grey q-pa-lg">
            {{ currentView === 'directory' && !directoryPath ? '无未分类话术' : '暂无话术' }}
          </div>
        </div>
      </template>
    </q-scroll-area>

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
</style>
