<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as ScriptService from '../../wailsjs/go/services/ScriptService'
import * as CategoryService from '../../wailsjs/go/services/CategoryService'
import { models } from '../../wailsjs/go/models'
import { useQuasar } from 'quasar'
import CategoryList from '../components/CategoryList.vue'
import ScriptItem from '../components/ScriptItem.vue'

const $q = useQuasar()
const route = useRoute()
const router = useRouter()

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
const editorBody = ref('')
const editorTags = ref('')
const editorCategoryId = ref<number | null>(null)
type MediaAttachmentType = 'image' | 'video'
interface MediaAttachmentRow {
  type: MediaAttachmentType
  url: string
}
const mediaAttachments = ref<MediaAttachmentRow[]>([])
const mediaTypeOptions: { label: string; value: MediaAttachmentType }[] = [
  { label: '图片', value: 'image' },
  { label: '视频', value: 'video' },
]
const MAX_MEDIA_ATTACHMENTS = 10
const imageInputRef = ref<HTMLInputElement | null>(null)
const imageUploadTargetIndex = ref<number | null>(null)

const VIDEO_EXTENSIONS = ['.mp4', '.webm']
function isVideo(path: string): boolean {
  const lower = path.toLowerCase()
  return VIDEO_EXTENSIONS.some(ext => lower.endsWith(ext))
}

function mediaTypeFromUrl(url: string): MediaAttachmentType {
  return isVideo(url) ? 'video' : 'image'
}

function normalizeMediaAttachments(rows: MediaAttachmentRow[]): string[] {
  return rows
    .map(row => row.url.trim())
    .filter(url => !!url)
}

function addMediaAttachment() {
  if (mediaAttachments.value.length >= MAX_MEDIA_ATTACHMENTS) {
    $q.notify({ type: 'warning', message: `最多添加 ${MAX_MEDIA_ATTACHMENTS} 个媒体附件` })
    return
  }
  mediaAttachments.value.push({ type: 'image', url: '' })
}

function removeMediaAttachment(index: number) {
  mediaAttachments.value.splice(index, 1)
}

function updateMediaType(index: number, type: MediaAttachmentType | null) {
  const nextType = type || 'image'
  const item = mediaAttachments.value[index]
  if (!item) return
  if (item.type !== nextType) {
    item.type = nextType
    item.url = ''
    return
  }
  item.type = nextType
}

function openEditor(script?: models.Script) {
  if (script) {
    editingScript.value = script
    editorBody.value = script.content || ''
    editorTags.value = script.tags
    editorCategoryId.value = script.category_id || null
    mediaAttachments.value = (script.images ? JSON.parse(script.images) : []).map((url: string) => ({
      type: mediaTypeFromUrl(url),
      url,
    }))
  } else {
    editingScript.value = null
    editorBody.value = ''
    editorTags.value = ''
    // Default category: current directory if in directory view
    if (currentView.value === 'directory' && directoryPath.value) {
      editorCategoryId.value = directoryPath.value.id
    } else {
      editorCategoryId.value = null
    }
    mediaAttachments.value = [{ type: 'image', url: '' }]
  }
  showEditor.value = true
}

watch(
  () => route.query.new,
  newToken => {
    if (!newToken) return
    openEditor()
    const query = { ...route.query }
    delete query.new
    router.replace({ path: route.path, query })
  },
  { immediate: true },
)

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = String(reader.result || '')
      const base64 = result.split(',')[1]
      if (!base64) {
        reject(new Error('invalid base64 content'))
        return
      }
      resolve(base64)
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

async function onImageFileSelected(event: Event) {
  const index = imageUploadTargetIndex.value
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file || index === null) {
    imageUploadTargetIndex.value = null
    if (input) input.value = ''
    return
  }

  try {
    const base64 = await fileToBase64(file)
    const ext = file.name.substring(file.name.lastIndexOf('.')) || '.png'
    // @ts-ignore
    const path = await ScriptService.SaveScriptImage(base64, ext)
    mediaAttachments.value[index].url = path
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: '图片上传失败' })
  } finally {
    imageUploadTargetIndex.value = null
    if (input) input.value = ''
  }
}

async function handleMediaUpload(index: number) {
  const row = mediaAttachments.value[index]
  if (!row) return

  if (row.type === 'video') {
    try {
      // @ts-ignore
      const path = await ScriptService.SelectAndSaveMedia()
      if (path) {
        mediaAttachments.value[index].url = path
      }
    } catch (e: any) {
      const msg = String(e)
      if (!msg.includes('no file selected')) {
        console.error(e)
        $q.notify({ type: 'negative', message: msg || '视频上传失败' })
      }
    }
    return
  }

  imageUploadTargetIndex.value = index
  imageInputRef.value?.click()
}

function removeEmptyMediaRows() {
  mediaAttachments.value = mediaAttachments.value.filter(item => item.url.trim())
}

async function saveScript() {
  try {
    removeEmptyMediaRows()
    const images = normalizeMediaAttachments(mediaAttachments.value)
    const content = editorBody.value.trim()

    if (!content) {
      $q.notify({ type: 'warning', message: '文本内容不能为空' })
      return
    }

    if (editingScript.value) {
      // @ts-ignore
      await ScriptService.UpdateScript(
        editingScript.value.id,
        content,
        editorTags.value,
        editorCategoryId.value,
        JSON.stringify(images),
      )
      $q.notify({ type: 'positive', message: '话术已更新' })
    } else {
      // @ts-ignore
      await ScriptService.CreateScript(
        content,
        editorTags.value,
        editorCategoryId.value,
        JSON.stringify(images),
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
  <q-page :style-fn="tweakPageHeight" class="q-pa-md column no-wrap" style="overflow: hidden">
    <!-- Header: Search -->
    <div class="row q-mb-md items-center">
      <q-input
        dense
        outlined
        v-model="searchText"
        placeholder="Search scripts..."
        class="col search-box"
        debounce="300"
        bg-color="white">
        <template v-slot:append>
          <q-icon name="search" />
        </template>
      </q-input>
    </div>

    <!-- View Switcher -->
    <div class="row q-mb-md">
      <q-btn-group unelevated class="full-width view-switcher">
        <q-btn
          :class="['col switch-btn', { 'switch-btn--active': currentView === 'timeline' }]"
          flat
          no-caps
          label="Timeline"
          class="col"
          size="sm"
          @click="currentView = 'timeline'" />
        <q-btn
          :class="['col switch-btn', { 'switch-btn--active': currentView === 'directory' }]"
          flat
          no-caps
          label="By Category"
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
    <div class="col relative-position scroll" style="flex: 1; overflow-y: auto; overflow-x: hidden; padding-right: 8px">
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
      <q-card class="script-dialog-card script-editor-dialog">
        <q-card-section class="row items-center editor-header">
          <div class="text-h6">{{ editingScript ? '编辑话术' : '新建话术' }}</div>
          <q-space />
          <q-btn flat round dense icon="close" color="grey-6" v-close-popup />
        </q-card-section>

        <q-card-section class="editor-body">
          <div class="editor-field-label">所属目录</div>
          <q-select
            v-model="editorCategoryId"
            :options="categoryOptions"
            dense
            outlined
            emit-value
            map-options
            popup-content-class="script-dialog-card"
            class="q-mb-md" />

          <div class="editor-field-label">文本内容</div>
          <q-input
            v-model="editorBody"
            type="textarea"
            outlined
            autogrow
            autofocus
            placeholder="请输入话术内容..."
            class="q-mb-md"
            style="max-height: 240px; overflow-y: auto" />

          <div class="row items-center q-mb-sm">
            <div class="editor-field-label q-mb-none">媒体附件</div>
            <q-space />
            <q-btn
              flat
              no-caps
              color="primary"
              icon="add"
              label="添加媒体"
              :disable="mediaAttachments.length >= MAX_MEDIA_ATTACHMENTS"
              @click="addMediaAttachment" />
          </div>

          <div class="q-gutter-y-sm">
            <div v-for="(media, index) in mediaAttachments" :key="`media-${index}`" class="media-row">
              <q-select
                :model-value="media.type"
                :options="mediaTypeOptions"
                emit-value
                map-options
                outlined
                dense
                class="media-type-select"
                @update:model-value="value => updateMediaType(index, value as MediaAttachmentType)" />
              <q-btn
                unelevated
                color="primary"
                text-color="white"
                no-caps
                :label="media.url ? '重新上传' : '上传'"
                class="media-upload-btn"
                @click="handleMediaUpload(index)" />
              <div class="media-status" :class="{ 'media-status--uploaded': !!media.url }">
                {{ media.url ? '已上传' : '未上传' }}
              </div>
              <q-btn
                flat
                round
                dense
                color="negative"
                icon="delete"
                class="media-delete-btn"
                @click="removeMediaAttachment(index)" />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="editor-footer">
          <q-btn flat no-caps label="取消" color="grey-8" v-close-popup />
          <q-btn unelevated no-caps color="primary" icon="save" label="保存话术" @click="saveScript" />
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
    <input type="file" ref="imageInputRef" style="display: none" accept="image/*" @change="onImageFileSelected" />
  </q-page>
</template>

<style lang="scss">
.search-box {
  :deep(.q-field__control) {
    border-radius: 12px;
    background: #eff2f7;
    border-color: #eff2f7;
  }
}

.view-switcher {
  border-radius: 14px;
  padding: 4px;
  background: #eceff5;
}

.switch-btn {
  border-radius: 12px;
  color: #596275;
  font-weight: 600;
}

.switch-btn--active {
  background: #ffffff;
  color: #1f2735;
  box-shadow: 0 2px 6px rgba(27, 39, 61, 0.12);
}

.script-dialog-card {
  color: rgba(0, 0, 0, 0.87);
  .q-field__label {
    color: rgba(0, 0, 0, 0.6);
  }
}

.script-editor-dialog {
  width: 560px;
  max-width: calc(100vw - 20px);
  border-radius: 12px;
}

.editor-header {
  border-bottom: 1px solid #eceff4;
}

.editor-body {
  max-height: 70vh;
  overflow-y: auto;
  padding: 14px 16px;
}

.editor-field-label {
  font-size: 14px;
  font-weight: 600;
  color: #3b4558;
  margin-bottom: 8px;
}

.media-row {
  display: grid;
  grid-template-columns: 88px 96px 48px 30px;
  align-items: center;
  column-gap: 8px;
  padding: 8px;
  border: 1px solid #e6ebf4;
  border-radius: 8px;
  background: #fbfcff;
}

.media-type-select {
  min-width: 0;
}

.media-upload-btn {
  min-width: 96px;
  height: 36px;
  border-radius: 8px;
  font-weight: 600;
}

.media-status {
  text-align: center;
  font-size: 12px;
  line-height: 1;
  font-weight: 600;
  color: #8a94a8;
  background: #eef1f6;
  border-radius: 10px;
  padding: 4px 0;
}

.media-status--uploaded {
  color: #1e8e3e;
  background: #eaf7ee;
}

.media-delete-btn {
  justify-self: center;
  background: #fff1f1;
}

.editor-footer {
  border-top: 1px solid #eceff4;
  padding: 12px 16px;
}

@media (max-width: 480px) {
  .media-row {
    grid-template-columns: 84px 92px 44px 28px;
    column-gap: 6px;
  }

  .media-upload-btn {
    min-width: 92px;
    height: 34px;
  }

  .media-status {
    font-size: 11px;
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
