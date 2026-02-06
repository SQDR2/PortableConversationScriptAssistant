<template>
  <q-dialog v-model="isOpen" :persistent="forced">
    <q-card style="min-width: 350px">
      <q-card-section>
        <div class="text-h6">翻译配置</div>
        <div class="text-caption text-grey">
          <span v-if="forced">检测到翻译服务凭证缺失或无效。请输入腾讯云 API 密钥以继续使用翻译功能。</span>
          <span v-else>在此修改腾讯云 API 密钥。</span>
        </div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-input v-model="secretId" label="Secret ID" dense autofocus :rules="[(val: string) => !!val || '此项必填']" />
        <q-input
          v-model="secretKey"
          label="Secret Key"
          dense
          type="password"
          :rules="[(val: string) => !!val || '此项必填']" />
      </q-card-section>

      <q-card-actions align="right" class="text-primary">
        <q-btn flat label="取消" v-if="!forced" @click="isOpen = false" />
        <q-btn flat label="保存" @click="save" :disable="!isValid" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { UpdateCredentials } from '../../wailsjs/go/services/TranslationService'
import { useQuasar } from 'quasar'

const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    forced?: boolean
  }>(),
  {
    modelValue: false,
    forced: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'saved'): void
}>()

const $q = useQuasar()
const secretId = ref('')
const secretKey = ref('')

const isOpen = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
})

const isValid = computed(() => {
  return secretId.value.trim() !== '' && secretKey.value.trim() !== ''
})

async function save() {
  if (!isValid.value) return

  try {
    await UpdateCredentials(secretId.value, secretKey.value)
    $q.notify({
      type: 'positive',
      message: '配置保存成功',
      position: 'top',
    })
    emit('saved')
    isOpen.value = false
  } catch (error) {
    console.error('Failed to save credentials:', error)
    $q.notify({
      type: 'negative',
      message: '配置保存失败: ' + error,
      position: 'top',
    })
  }
}
</script>
