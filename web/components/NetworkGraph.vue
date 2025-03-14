<template>
  <div class="network-container">
    <div ref="container" class="network-graph"></div>
    
    <!-- Модальное окно для создания соединения -->
    <div v-if="showConnectionModal" class="modal-overlay">
      <div class="modal-content">
        <h3 class="text-lg font-semibold mb-4">Создать соединение между роутерами</h3>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Исходный роутер</label>
          <select v-model="selectedSourceRouter" class="w-full p-2 border rounded">
            <option v-for="router in routers" :key="router.id" :value="router">
              {{ router.name }} ({{ router.ip_address }})
            </option>
          </select>
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Целевой роутер</label>
          <select v-model="selectedTargetRouter" class="w-full p-2 border rounded">
            <option v-for="router in routers" :key="router.id" :value="router">
              {{ router.name }} ({{ router.ip_address }})
            </option>
          </select>
        </div>
        <div class="flex justify-end gap-2">
          <button @click="closeConnectionModal" class="px-4 py-2 border rounded hover:bg-gray-50">
            Отмена
          </button>
          <button @click="createConnection" class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
            Создать соединение
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { Network } from 'vis-network'
import { DataSet } from 'vis-data'
import { useApi } from '~/composables/useApi'

const props = defineProps({
  routers: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['router-click', 'connection-created'])
const container = ref(null)
const showConnectionModal = ref(false)
const selectedSourceRouter = ref(null)
const selectedTargetRouter = ref(null)
let network = null
const connections = ref([])

const api = useApi()

// Загрузка соединений
async function loadConnections() {
  try {
    connections.value = await api.getRouterConnections()
  } catch (error) {
    console.error('Failed to load connections:', error)
  }
}

// Преобразование данных роутеров в формат для vis.js
function prepareData(routers) {
  const nodes = new DataSet(
    routers.map(router => ({
      id: router.id,
      label: `${router.name}\n${router.ip_address}`,
      title: `IP: ${router.ip_address}\nStatus: ${router.status}`,
      color: {
        background: router.status === 'active' ? '#dcfce7' : '#fee2e2',
        border: router.status === 'active' ? '#166534' : '#991b1b'
      },
      shape: 'circle',
      size: 30
    }))
  )

  // Создаем ребра на основе данных о соединениях
  const edges = new DataSet(
    connections.value.map(conn => ({
      from: routers.find(r => r.ip_address === conn.source_ip)?.id,
      to: routers.find(r => r.ip_address === conn.target_ip)?.id,
      arrows: 'to',
      title: `${conn.source_ip} → ${conn.target_ip}`
    })).filter(edge => edge.from && edge.to) // Фильтруем невалидные соединения
  )

  return { nodes, edges }
}

// Опции для визуализации
const options = {
  nodes: {
    shape: 'circle',
    margin: 10,
    borderWidth: 2,
    shadow: true
  },
  edges: {
    width: 2,
    color: {
      color: '#999999',
      highlight: '#666666',
      hover: '#666666'
    },
    smooth: {
      type: 'continuous'
    },
    arrows: {
      to: {
        enabled: true,
        scaleFactor: 0.5
      }
    }
  },
  physics: {
    enabled: true,
    barnesHut: {
      gravitationalConstant: -2000,
      centralGravity: 0.3,
      springLength: 150
    }
  },
  interaction: {
    hover: true,
    tooltipDelay: 200
  }
}

// Инициализация сети
async function initNetwork() {
  if (!container.value) return
  
  await loadConnections() // Загружаем соединения перед инициализацией графа
  const data = prepareData(props.routers)
  network = new Network(container.value, data, options)

  network.on('click', function(params) {
    if (params.nodes.length > 0) {
      emit('router-click', params.nodes[0])
    }
  })

  network.on('doubleClick', function(params) {
    if (params.nodes.length > 0) {
      showConnectionModal.value = true
      const selectedRouter = props.routers.find(r => r.id === params.nodes[0])
      selectedSourceRouter.value = selectedRouter
    }
  })
}

async function createConnection() {
  if (!selectedSourceRouter.value || !selectedTargetRouter.value) {
    alert('Выберите оба роутера')
    return
  }

  try {
    await api.connectRouter({
      router_from_ip: selectedSourceRouter.value.ip_address,
      router_to_ip: selectedTargetRouter.value.ip_address
    })

    // Обновляем соединения после создания нового
    await loadConnections()
    const data = prepareData(props.routers)
    network.setData(data)

    emit('connection-created')
    closeConnectionModal()
  } catch (error) {
    alert('Ошибка при создании соединения: ' + error.message)
  }
}

function closeConnectionModal() {
  showConnectionModal.value = false
  selectedSourceRouter.value = null
  selectedTargetRouter.value = null
}

watch(() => props.routers, async (newRouters) => {
  if (network) {
    await loadConnections() // Перезагружаем соединения при обновлении роутеров
    const data = prepareData(newRouters)
    network.setData(data)
  }
}, { deep: true })

onMounted(() => {
  initNetwork()
})
</script>

<style scoped>
.network-container {
  position: relative;
  width: 100%;
  height: 600px;
}

.network-graph {
  width: 100%;
  height: 100%;
  background: white;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
}

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