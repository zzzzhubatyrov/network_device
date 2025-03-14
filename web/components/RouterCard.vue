<template>
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex justify-between items-start mb-4">
      <div>
        <h3 class="text-lg font-medium">{{ router.name }}</h3>
        <p class="text-sm text-gray-500">IP: {{ router.ip_address }}</p>
      </div>
      <span
        :class="[
          'px-2 py-1 text-sm rounded-full',
          router.connected
            ? 'bg-green-100 text-green-800'
            : 'bg-red-100 text-red-800'
        ]"
      >
        {{ router.connected ? 'Connected' : 'Disconnected' }}
      </span>
    </div>
    <div class="flex space-x-2">
      <button
        @click="$emit('connect', router)"
        class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
      >
        {{ router.connected ? 'Disconnect' : 'Connect' }}
      </button>
      <button
        @click="$emit('configure', router)"
        class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded"
      >
        Configure
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  router: {
    type: Object,
    required: true
  }
})

defineEmits(['connect', 'configure'])
</script> 