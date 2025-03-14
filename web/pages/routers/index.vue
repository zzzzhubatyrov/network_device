<template>
  <div class="container mx-auto p-4">
    <div class="mb-4 flex justify-between items-center">
      <h1 class="text-2xl font-bold">Управление роутерами</h1>
      <button
        @click="showAddForm = true"
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Добавить роутер
      </button>
    </div>

    <!-- Загрузка и ошибки -->
    <div v-if="loading" class="text-center py-4">
      Загрузка роутеров...
    </div>
    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <!-- Граф сети -->
    <div v-else class="mb-8">
      <NetworkGraph
        :routers="routers"
        @router-click="handleRouterClick"
        @connection-created="loadRouters"
      />
    </div>

    <!-- Список роутеров -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="router in routers"
        :key="router.id"
        class="border rounded-lg p-4 bg-white shadow-sm"
      >
        <div class="flex justify-between items-start mb-2">
          <div>
            <h3 class="font-semibold">{{ router.name }}</h3>
            <p class="text-sm text-gray-600">IP: {{ router.ip_address }}</p>
          </div>
          <span
            :class="[
              'px-2 py-1 text-xs rounded-full',
              router.status === 'active'
                ? 'bg-green-100 text-green-800'
                : 'bg-red-100 text-red-800'
            ]"
          >
            {{ router.status === 'active' ? 'Активный' : 'Неактивный' }}
          </span>
        </div>
        <div class="flex gap-2">
          <button
            @click="openConfigModal(router)"
            class="px-3 py-1 text-sm border rounded hover:bg-gray-50"
          >
            Настроить
          </button>
        </div>
      </div>
    </div>

    <!-- Модальное окно добавления роутера -->
    <div v-if="showAddForm" class="modal-overlay">
      <div class="modal-content">
        <h3 class="text-lg font-semibold mb-4">Добавить новый роутер</h3>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Имя роутера</label>
          <input
            v-model="newRouter.name"
            type="text"
            class="w-full p-2 border rounded"
            placeholder="Введите имя роутера"
          />
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">IP адрес</label>
          <input
            v-model="newRouter.ip_address"
            type="text"
            class="w-full p-2 border rounded"
            placeholder="Например: 192.168.1.1"
          />
        </div>
        <div class="flex justify-end gap-2">
          <button
            @click="showAddForm = false"
            class="px-4 py-2 border rounded"
          >
            Отмена
          </button>
          <button
            @click="createRouter"
            class="px-4 py-2 bg-blue-500 text-white rounded"
          >
            Добавить
          </button>
        </div>
      </div>
    </div>

    <!-- Модальное окно настройки роутера -->
    <RouterConfigModal
      :show="showConfigModal"
      :router="selectedRouter"
      @close="closeConfigModal"
      @saved="handleConfigSaved"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

definePageMeta({
  layout: 'default'
})

const api = useApi()
const routers = ref([])
const showAddForm = ref(false)
const showConfigModal = ref(false)
const loading = ref(false)
const error = ref(null)
const selectedRouter = ref(null)

const newRouter = ref({
  name: '',
  ip_address: ''
})

async function loadRouters() {
  loading.value = true
  error.value = null
  try {
    routers.value = await api.getAllRouters()
  } catch (err) {
    error.value = 'Ошибка при загрузке роутеров: ' + err.message
  } finally {
    loading.value = false
  }
}

async function createRouter() {
  if (!newRouter.value.name || !newRouter.value.ip_address) {
    alert('Заполните все поля')
    return
  }

  try {
    await api.createRouter(newRouter.value)
    showAddForm.value = false
    newRouter.value = { name: '', ip_address: '' }
    await loadRouters()
  } catch (err) {
    alert('Ошибка при создании роутера: ' + err.message)
  }
}

function openConfigModal(router) {
  selectedRouter.value = router
  showConfigModal.value = true
}

function closeConfigModal() {
  showConfigModal.value = false
  selectedRouter.value = null
}

async function handleConfigSaved() {
  await loadRouters()
}

function handleRouterClick(routerId) {
  const router = routers.value.find(r => r.id === routerId)
  if (router) {
    console.log('Выбран роутер:', router)
  }
}

onMounted(() => {
  loadRouters()
})
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
  max-width: 500px;
}
</style> 