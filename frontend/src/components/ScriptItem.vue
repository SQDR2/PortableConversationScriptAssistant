<script setup lang="ts">
import { models } from '../../wailsjs/go/models'
import { useQuasar } from 'quasar'

const props = defineProps<{
  script: models.Script
}>()

const emit = defineEmits(['edit', 'delete'])

const $q = useQuasar()

function copyContent() {
  if (navigator.clipboard) {
    navigator.clipboard
      .writeText(props.script.content)
      .then(() => {
        $q.notify({ type: 'positive', message: '已复制到剪贴板', timeout: 1000 })
      })
      .catch(() => {
        $q.notify({ type: 'negative', message: '复制失败' })
      })
  } else {
    // Fallback or Wails runtime clipboard?
    // runtime.ClipboardSetText(props.script.content) // Need to import runtime if wails
    // For now assume standard browser API works in WebView
    $q.notify({ type: 'warning', message: '当前环境不支持剪贴板操作' })
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
    <q-item-section class="text-left">
      <q-item-label class="script-content" lines="3">{{ script.content }}</q-item-label>
      <q-item-label caption class="row items-center q-mt-xs">
        <span class="text-grey-7">{{ new Date(script.created_at).toLocaleString() }}</span>
        <q-chip v-if="script.tags" size="xs" color="secondary" text-color="white" class="q-ml-sm">{{
          script.tags
        }}</q-chip>
      </q-item-label>
    </q-item-section>

    <q-item-section side top>
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
  </q-item>
</template>

<style lang="scss">
.script-item {
  background: rgba(255, 255, 255, 0.6);
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.9);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

.script-content {
  white-space: pre-wrap;
  font-size: 14px;
  line-height: 1.4;
  color: #333;
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
