<script setup lang="ts">
import { ref, watch } from 'vue'
import { models } from '../../wailsjs/go/models'
import * as CategoryService from '../../wailsjs/go/services/CategoryService'
import { useQuasar } from 'quasar'

const props = defineProps<{
  categories: models.Category[]
  selectedId: number | null // null means 'all' or specific uncategorized handling
}>()

const emit = defineEmits(['select', 'refresh'])

const $q = useQuasar()

// Dialog states
const showDialog = ref(false)
const isEditing = ref(false)
const dialogName = ref('')
const editingId = ref<number | null>(null)

function openCreateDialog() {
  isEditing.value = false
  dialogName.value = ''
  showDialog.value = true
}

function openEditDialog(cat: models.Category) {
  isEditing.value = true
  editingId.value = cat.id
  dialogName.value = cat.name
  showDialog.value = true
}

async function saveCategory() {
  if (!dialogName.value.trim()) return

  try {
    if (isEditing.value && editingId.value) {
      // @ts-ignore
      await CategoryService.UpdateCategory(editingId.value, dialogName.value)
    } else {
      // @ts-ignore
      await CategoryService.CreateCategory(dialogName.value)
    }
    showDialog.value = false
    emit('refresh')
    $q.notify({ type: 'positive', message: isEditing.value ? '目录已更新' : '目录已创建' })
  } catch (e) {
    $q.notify({ type: 'negative', message: '操作失败' })
  }
}

function confirmDelete(cat: models.Category) {
  $q.dialog({
    title: `删除目录 "${cat.name}"`,
    message: '请选择删除方式：',
    options: {
      type: 'radio',
      model: 'keep',
      items: [
        { label: '保留话术（移至未分类）', value: 'keep' },
        { label: '彻底删除（包含话术）', value: 'cascade', color: 'red' },
      ],
    },
    cancel: {
      label: '取消',
      flat: true,
    },
    ok: {
      label: '确定',
      color: 'primary',
    },
    class: 'script-dialog-card',
    persistent: true,
  }).onOk(async data => {
    const cascade = data === 'cascade'
    try {
      // @ts-ignore
      await CategoryService.DeleteCategory(cat.id, cascade)
      emit('refresh')
      // If selected was this, reset to null
      if (props.selectedId === cat.id) {
        emit('select', null)
      }
      $q.notify({ type: 'positive', message: '目录已删除' })
    } catch (e) {
      $q.notify({ type: 'negative', message: '删除失败' })
    }
  })
}
</script>

<template>
  <div class="category-list">
    <div class="row items-center justify-between q-mb-sm q-px-sm">
      <div class="text-subtitle2 text-grey-7">目录</div>
      <q-btn flat round dense icon="add" size="sm" color="primary" @click="openCreateDialog">
        <q-tooltip>新建目录</q-tooltip>
      </q-btn>
    </div>

    <q-list dense separator class="rounded-borders">
      <!-- Uncategorized / All option? Or handled by parent? -->
      <!-- Let's assume parent has an "All" or "Uncategorized" special item, 
                 CategoryList just lists defined categories -->

      <q-item
        v-for="cat in categories"
        :key="cat.id"
        clickable
        v-ripple
        :active="selectedId === cat.id"
        active-class="bg-blue-1 text-primary"
        @click="emit('select', cat.id)">
        <q-item-section avatar style="min-width: 32px">
          <q-icon name="folder" size="xs" />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ cat.name }}</q-item-label>
        </q-item-section>
        <q-item-section side v-if="selectedId === cat.id || $q.platform.is.desktop">
          <div class="row q-gutter-xs">
            <q-btn flat round dense icon="edit" size="xs" color="grey-6" @click.stop="openEditDialog(cat)" />
            <q-btn flat round dense icon="delete" size="xs" color="grey-6" @click.stop="confirmDelete(cat)" />
          </div>
        </q-item-section>
      </q-item>
    </q-list>

    <!-- Category Dialog -->
    <q-dialog v-model="showDialog">
      <q-card style="min-width: 300px" class="script-dialog-card">
        <q-card-section>
          <div class="text-h6">{{ isEditing ? '编辑目录' : '新建目录' }}</div>
        </q-card-section>

        <q-card-section>
          <q-input v-model="dialogName" label="目录名称" autofocus dense @keyup.enter="saveCategory" />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn label="保存" color="primary" @click="saveCategory" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<style lang="scss" scoped>
.category-list {
  /* Optional custom styles */
}
</style>
