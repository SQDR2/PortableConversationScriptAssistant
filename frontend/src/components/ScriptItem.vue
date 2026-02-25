<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { models } from '../../wailsjs/go/models'
import { useQuasar } from 'quasar'
import { Port } from '../../wailsjs/go/services/MediaServer'
import { WindowFullscreen, WindowUnfullscreen } from '../../wailsjs/runtime/runtime'

const props = defineProps<{
  script: models.Script
}>()

const emit = defineEmits(['edit', 'delete'])

const $q = useQuasar()
const showPreview = ref(false)
const showFullscreen = ref(false)
const fullscreenMedia = ref('')
const mediaPort = ref(0)

onMounted(async () => {
  document.addEventListener('fullscreenchange', handleNativeFullscreen)
  document.addEventListener('webkitfullscreenchange', handleNativeFullscreen)

  try {
    mediaPort.value = await Port()
  } catch (e) {
    console.error('Failed to get media server port:', e)
  }
})

const scriptImages = computed<string[]>(() => {
  try {
    return props.script.images ? JSON.parse(props.script.images) : []
  } catch (e) {
    return []
  }
})

const scriptTitle = computed(() => {
  const firstLine = (props.script.content || '').split('\n').find(line => line.trim()) || ''
  return firstLine.length > 48 ? `${firstLine.slice(0, 48)}...` : firstLine || 'Untitled Script'
})

const scriptPreview = computed(() => {
  const normalized = (props.script.content || '').replace(/\n+/g, '\n').trim()
  if (!normalized) return ''
  const withoutTitle = normalized.startsWith(scriptTitle.value)
    ? normalized.slice(scriptTitle.value.length).trim()
    : normalized
  return withoutTitle || normalized
})

const VIDEO_EXTENSIONS = ['.mp4', '.webm']
function isVideo(path: string): boolean {
  const lower = path.toLowerCase()
  return VIDEO_EXTENSIONS.some(ext => lower.endsWith(ext))
}

function videoMimeType(path: string): string {
  const lower = path.toLowerCase()
  if (lower.endsWith('.webm')) return 'video/webm'
  return 'video/mp4'
}

/** Build a localhost URL that hits the dedicated MediaServer for video playback */
function videoSrc(path: string): string {
  if (!mediaPort.value) return path // fallback
  // path is like "images/uuid.mp4" – extract the filename
  const filename = path.split('/').pop() || path
  return `http://127.0.0.1:${mediaPort.value}/${filename}`
}

function openFullscreen(media: string) {
  fullscreenMedia.value = media
  showFullscreen.value = true
  WindowFullscreen()
}

watch(showFullscreen, (val) => {
  if (!val) {
    WindowUnfullscreen()
  }
})

/**
 * Intercept native fullscreen requests on .preview-video elements.
 * Wails WebView does not render video correctly in native fullscreen,
 * so we immediately exit and redirect to our custom fullscreen dialog.
 */
function handleNativeFullscreen() {
  const el = (document as any).fullscreenElement || (document as any).webkitFullscreenElement
  if (!el || !el.classList?.contains('preview-video')) return

  // Exit native fullscreen immediately
  if (document.exitFullscreen) {
    document.exitFullscreen().catch(() => {})
  } else if ((document as any).webkitExitFullscreen) {
    (document as any).webkitExitFullscreen()
  }

  // Find the matching media path from the video source
  const source = el.querySelector('source')
  const src = source?.getAttribute('src') || ''
  for (const media of scriptImages.value) {
    if (videoSrc(media) === src) {
      openFullscreen(media)
      break
    }
  }
}

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleNativeFullscreen)
  document.removeEventListener('webkitfullscreenchange', handleNativeFullscreen)
})

async function copyContent() {
  try {
    await navigator.clipboard.writeText(props.script.content)
    $q.notify({ type: 'positive', message: '已复制文本内容', timeout: 1000 })
  } catch (e) {
    console.error('Copy failed:', e)
    $q.notify({ type: 'negative', message: '复制失败' })
  }
}

function confirmDelete() {
  $q.dialog({
    title: '确认删除',
    message: '确定要删除这条话术吗？此操作无法恢复。',
    cancel: {
      label: '取消',
      flat: true,
    },
    ok: {
      label: '删除',
      color: 'negative',
    },
    persistent: true,
    class: 'script-dialog-card', // Ensure dialog styling
  }).onOk(() => {
    emit('delete', props.script)
  })
}
</script>

<template>
  <q-item class="script-item q-mb-md q-pa-md relative-position">
    <q-item-section class="text-left cursor-pointer" @click="showPreview = true">
      <div class="row items-start no-wrap">
        <div class="col">
          <q-item-label class="script-title" lines="1">{{ scriptTitle }}</q-item-label>
          <q-item-label class="script-content q-mt-xs" lines="3">{{ scriptPreview }}</q-item-label>
          <q-item-label caption class="row items-center q-mt-sm">
            <span class="text-grey-7">{{ new Date(script.created_at).toLocaleString() }}</span>
            <q-chip v-if="script.tags" size="sm" color="grey-2" text-color="grey-8" class="q-ml-sm">{{ script.tags }}</q-chip>
          </q-item-label>
        </div>
        <div v-if="scriptImages.length > 0" class="q-ml-md">
          <div
            v-if="isVideo(scriptImages[0])"
            class="rounded-borders thumbnail-preview flex flex-center bg-grey-3"
            style="width: 60px; height: 60px">
            <q-icon name="play_circle" color="primary" size="md" />
          </div>
          <q-img
            v-else
            :src="scriptImages[0]"
            style="width: 60px; height: 60px"
            class="rounded-borders thumbnail-preview" />
          <div v-if="scriptImages.length > 1" class="text-right">
            <q-badge color="grey-8" text-color="white" size="xs" floating>+{{ scriptImages.length - 1 }}</q-badge>
          </div>
        </div>
      </div>
    </q-item-section>

    <q-btn flat round dense icon="more_horiz" class="script-actions-trigger" @click.stop>
      <q-menu anchor="bottom right" self="top right">
        <q-list style="min-width: 120px">
          <q-item clickable v-close-popup @click="copyContent">
            <q-item-section avatar><q-icon name="content_copy" /></q-item-section>
            <q-item-section>复制</q-item-section>
          </q-item>
          <q-item clickable v-close-popup @click="emit('edit', script)">
            <q-item-section avatar><q-icon name="edit" /></q-item-section>
            <q-item-section>编辑</q-item-section>
          </q-item>
          <q-item clickable v-close-popup @click="confirmDelete">
            <q-item-section avatar><q-icon name="delete" color="negative" /></q-item-section>
            <q-item-section>删除</q-item-section>
          </q-item>
        </q-list>
      </q-menu>
    </q-btn>

    <!-- Preview Dialog -->
    <q-dialog v-model="showPreview">
      <q-card class="script-dialog-card" style="width: 400px; max-width: 90vw">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-subtitle1 text-weight-bold">详情预览</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="scroll" style="max-height: 70vh">
          <div class="text-body1 q-mb-md preserve-whitespace">{{ script.content }}</div>

          <div v-if="scriptImages.length > 0" class="q-gutter-y-md">
            <template v-for="(media, index) in scriptImages" :key="index">
              <!-- Video player -->
              <div v-if="isVideo(media)" class="relative-position">
                <video
                  controls
                  controlslist="nofullscreen"
                  preload="metadata"
                  class="rounded-borders shadow-1 full-width preview-video">
                  <source :src="videoSrc(media)" :type="videoMimeType(media)" />
                </video>
                <q-btn
                  icon="fullscreen"
                  round flat
                  color="white"
                  size="sm"
                  class="video-fullscreen-btn"
                  @click="openFullscreen(media)">
                  <q-tooltip>全屏播放</q-tooltip>
                </q-btn>
              </div>
              <!-- Image -->
              <img
                v-else
                :src="media"
                class="rounded-borders shadow-1 full-width preview-img"
                loading="lazy" />
            </template>
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="复制全文" color="primary" @click="copyContent" />
          <q-btn flat label="关闭" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Fullscreen video dialog (custom replacement for native fullscreen) -->
    <q-dialog v-model="showFullscreen" maximized transition-show="fade" transition-hide="fade">
      <div class="fit bg-black flex flex-center" style="position: relative">
        <video
          controls
          autoplay
          class="fullscreen-video-player"
          @click.stop>
          <source :src="videoSrc(fullscreenMedia)" :type="videoMimeType(fullscreenMedia)" />
        </video>
        <q-btn
          icon="close"
          flat round
          color="white"
          size="md"
          class="absolute-top-right q-ma-md"
          style="z-index: 1"
          @click="showFullscreen = false">
          <q-tooltip>关闭全屏</q-tooltip>
        </q-btn>
      </div>
    </q-dialog>
  </q-item>
</template>

<style lang="scss">
.script-item {
  background: #ffffff;
  border-radius: 10px;
  border: 1px solid #e9ecf2;
  box-shadow: 0 1px 3px rgba(18, 25, 38, 0.04);
  transition: all 0.2s;

  &:hover {
    box-shadow: 0 8px 20px rgba(26, 35, 52, 0.08);

    .script-actions-trigger {
      opacity: 1;
    }
  }
}

.script-title {
  font-size: 15px;
  line-height: 1.3;
  font-weight: 700;
  color: #1f2532;
}

.script-content {
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 14px;
  line-height: 1.4;
  color: #5d6678;
}

.script-actions-trigger {
  position: absolute;
  right: 8px;
  top: 8px;
  color: #7a8395;
  opacity: 0;
  transition: opacity 0.18s ease;
}

.script-item:focus-within .script-actions-trigger {
  opacity: 1;
}

.preserve-whitespace {
  white-space: pre-wrap;
  word-break: break-all;
}

.preview-img {
  display: block;
  max-width: 100%;
  height: auto;
}

.preview-video {
  display: block;
  max-width: 100%;
  height: auto;
}

/* Hide native fullscreen button – Wails WebView does not support
   the HTML5 Fullscreen API correctly for embedded <video>. */
video::-webkit-media-controls-fullscreen-button {
  display: none;
}

.video-fullscreen-btn {
  position: absolute;
  right: 4px;
  bottom: 36px;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1;
}

.fullscreen-video-player {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.thumbnail-preview {
  border-radius: 10px;
  border: 1px solid #eaedf4;
}
</style>
