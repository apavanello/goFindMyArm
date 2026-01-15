<script setup>
defineProps({
  loading: Boolean
})
defineEmits(['scan'])
</script>

<template>
  <button 
    @click="$emit('scan')"
    :disabled="loading"
    class="relative inline-flex items-center justify-center px-8 py-3 overflow-hidden font-medium text-white transition duration-300 ease-out border-2 border-primary rounded-full shadow-md group disabled:opacity-75 disabled:cursor-not-allowed">
    <span class="absolute inset-0 flex items-center justify-center w-full h-full text-white duration-300 -translate-x-full bg-primary group-hover:translate-x-0 ease">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
    </span>
    <span class="absolute flex items-center justify-center w-full h-full text-primary transition-all duration-300 transform group-hover:translate-x-full ease bg-white group-disabled:bg-gray-100 group-disabled:translate-x-0 group-disabled:text-gray-400">
        {{ loading ? 'Scanning...' : 'Scan Network' }}
    </span>
    <span class="relative invisible">Scan Network</span>
  </button>
</template>
