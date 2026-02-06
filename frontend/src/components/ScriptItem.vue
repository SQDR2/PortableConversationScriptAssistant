<script setup lang="ts">
import { ref, computed } from 'vue'
import { models } from '../../wailsjs/go/models'
import { useQuasar } from 'quasar'

const props = defineProps<{
  script: models.Script
}>()

const emit = defineEmits(['edit', 'delete'])

const $q = useQuasar()
const showPreview = ref(false)

const scriptImages = computed<string[]>(() => {
  try {
    return props.script.images ? JSON.parse(props.script.images) : []
  } catch (e) {
    return []
  }
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
  <q-item class="script-item q-mb-sm">
    <q-item-section class="text-left cursor-pointer" @click="showPreview = true">
      <div class="row items-start no-wrap">
        <div class="col">
          <q-item-label class="script-content" lines="3">{{ script.content }}</q-item-label>
          <q-item-label caption class="row items-center q-mt-xs">
            <span class="text-grey-7">{{ new Date(script.created_at).toLocaleString() }}</span>
            <q-chip v-if="script.tags" size="xs" color="secondary" text-color="white" class="q-ml-sm">{{
              script.tags
            }}</q-chip>
          </q-item-label>
        </div>
        <div v-if="scriptImages.length > 0" class="q-ml-sm">
          <q-img :src="scriptImages[0]" style="width: 40px; height: 40px" class="rounded-borders thumbnail-preview" />
          <div v-if="scriptImages.length > 1" class="text-right">
            <q-badge color="orange" size="xs" floating>+{{ scriptImages.length - 1 }}</q-badge>
          </div>
        </div>
      </div>
    </q-item-section>

    <q-item-section side>
      <div class="row q-gutter-xs">
        <q-btn flat round dense size="sm" icon="content_copy" color="primary" @click.stop="copyContent">
          <q-tooltip>复制</q-tooltip>
        </q-btn>
        <q-btn flat round dense size="sm" icon="edit" color="grey-7" @click.stop="emit('edit', script)">
          <q-tooltip>编辑</q-tooltip>
        </q-btn>
        <q-btn flat round dense size="sm" icon="delete" color="negative" @click.stop="confirmDelete">
          <q-tooltip>删除</q-tooltip>
        </q-btn>
      </div>
    </q-item-section>

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
            <q-img
              v-for="(img, index) in scriptImages"
              :key="index"
              :src="img"
              class="rounded-borders shadow-1"
              spinner-color="primary" />
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="复制内容" color="primary" @click="copyContent" />
          <q-btn flat label="关闭" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-item>
</template>

<style lang="scss">
.script-item {
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.85);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }
}

.script-content {
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 14px;
  line-height: 1.4;
  color: #333;
}

.preserve-whitespace {
  white-space: pre-wrap;
  word-break: break-all;
}

.thumbnail-preview {
  border: 1px solid rgba(0, 0, 0, 0.05);
}

body.body--dark {
  .script-item {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);

    &:hover {
      background: rgba(255, 255, 255, 0.1);
    }
  }
  .script-content {
    color: #eee;
  }
}
</style>
