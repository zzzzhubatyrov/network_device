<template>
  <div v-if="show" class="modal-overlay">
    <div class="modal-content">
      <h3 class="text-lg font-semibold mb-4">Настройка роутера {{ router.name }}</h3>

      <!-- Раздел настройки роутера -->
      <div class="mb-6">
        <h4 class="text-md font-semibold mb-2">Настройки роутера</h4>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Имя</label>
          <input
            v-model="routerConfig.name"
            type="text"
            class="w-full p-2 border rounded"
            placeholder="Введите имя"
          />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Статус</label>
          <select v-model="routerConfig.status" class="w-full p-2 border rounded">
            <option value="active">Активный</option>
            <option value="inactive">Неактивный</option>
          </select>
        </div>
      </div>

      <!-- Раздел настройки портов -->
      <div class="mb-6">
        <h4 class="text-md font-semibold mb-2">Настройки порта</h4>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Статус порта</label>
          <select v-model="config.status" class="w-full p-2 border rounded">
            <option value="up">Активный</option>
            <option value="down">Неактивный</option>
          </select>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Номер порта</label>
          <input
            v-model.number="config.portNumber"
            type="number"
            class="w-full p-2 border rounded"
            min="1"
            max="65535"
          />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Скорость</label>
          <select v-model="config.speed" class="w-full p-2 border rounded">
            <option value="10">10 Мбит/с</option>
            <option value="100">100 Мбит/с</option>
            <option value="1000">1 Гбит/с</option>
          </select>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Режим дуплекса</label>
          <select v-model="config.duplexMode" class="w-full p-2 border rounded">
            <option value="full">Полный дуплекс</option>
            <option value="half">Полудуплекс</option>
          </select>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Описание</label>
          <textarea
            v-model="config.description"
            class="w-full p-2 border rounded"
            rows="3"
            placeholder="Введите описание порта"
          ></textarea>
        </div>
      </div>

      <!-- Список портов роутера -->
      <div v-if="router.ports.length" class="mb-6">
        <h4 class="text-md font-semibold mb-2">Порты роутера</h4>
        <ul>
          <li v-for="port in router.ports" :key="port.id" class="mb-2">
            <strong>Порт {{ port.number }}:</strong> {{ port.status }} 
            <span v-if="port.description">- {{ port.description }}</span>
          </li>
        </ul>
      </div>

      <div class="flex justify-end gap-2">
        <button @click="$emit('close')" class="px-4 py-2 border rounded">
          Отмена
        </button>
        <button @click="saveConfig" class="px-4 py-2 bg-blue-500 text-white rounded">
          Сохранить
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useApi } from '~/composables/useApi'

const props = defineProps({
  show: Boolean,
  router: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'saved'])
const api = useApi()

const config = ref({
  routerId: null,
  portNumber: 1,
  protocol: 'tcp',
  status: 'up',
  speed: '1000',
  duplexMode: 'full',
  description: ''
})

const routerConfig = ref({
  name: '',
  status: 'active',
})

watch(() => props.router, (newRouter) => {
  if (newRouter) {
    config.value.routerId = newRouter.id
    routerConfig.value.ip_address = newRouter.ip_address
    routerConfig.value.status = newRouter.status
  }
}, { immediate: true })

async function saveConfig() {
  try {
    // Save router configuration
    await api.configureRouter(routerConfig.value)
    // Save port configuration
    await api.configureRouterPort(config.value)

    emit('saved')
    emit('close')
  } catch (error) {
    alert('Ошибка при настройке роутера: ' + error.message)
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}
</style>
