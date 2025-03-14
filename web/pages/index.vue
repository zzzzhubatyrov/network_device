<template>
  <div class="flex flex-col items-center justify-center min-h-[80vh]">
    <h1 class="text-4xl font-bold mb-4">Welcome to Network Management</h1>
    <p class="text-gray-600 mb-8">Manage your network infrastructure with ease</p>
    <NuxtLink
      to="/routers"
      class="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium"
    >
      Go to Routers
    </NuxtLink>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'default'
})
</script>