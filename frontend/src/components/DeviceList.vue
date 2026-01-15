<script setup>
import { computed } from 'vue'

const props = defineProps({
  devices: {
    type: Array,
    required: true
  },
  viewMode: {
    type: String,
    default: 'grid' // 'grid' or 'list'
  }
})

function copyToClipboard(text) {
    navigator.clipboard.writeText(text);
    // Could add toast notification here
}
</script>

<template>
  <div v-if="devices.length === 0" class="text-center text-gray-500 py-10">
    No devices found yet. Click scan to start.
  </div>

  <!-- Grid View -->
  <div v-else-if="viewMode === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    <div v-for="device in devices" :key="device.ip" class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-shadow relative">
        <div class="flex items-start justify-between">
            <div class="bg-purple-50 p-3 rounded-full">
                <!-- Icon Placeholder (Raspberry Pi / Chip) -->
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                </svg>
            </div>
            <div class="flex items-center space-x-2">
                 <span class="w-3 h-3 rounded-full bg-green-500"></span>
            </div>
        </div>
        
        <div class="mt-4">
            <h3 class="text-lg font-bold text-gray-800">{{ device.hostname || 'Unknown Device' }}</h3>
            <div class="flex items-center space-x-2 mt-1">
                <span class="text-primary font-mono font-medium">{{ device.ip }}</span>
                <button @click="copyToClipboard(device.ip)" class="text-gray-400 hover:text-primary transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 01-2-2V5" />
                    </svg>
                </button>
            </div>
            <p class="text-xs text-gray-500 mt-2 truncate">MAC: {{ device.mac }}</p>
            <p class="text-xs text-gray-400">{{ device.os }}</p>
        </div>
        
        <!-- Actions (Context Menu Placeholder) -->
        <div class="mt-4 pt-4 border-t border-gray-50 flex justify-end">
            <button class="text-sm text-gray-600 hover:text-primary font-medium">Remote > </button>
        </div>
    </div>
  </div>

  <!-- List View -->
  <div v-else class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Hostname</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">IP Address</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">MAC</th>
          <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr v-for="device in devices" :key="device.ip">
            <td class="px-6 py-4 whitespace-nowrap">
                <span class="lex items-center">
                    <span class="w-2.5 h-2.5 rounded-full bg-green-500 inline-block mr-2"></span>
                    <span class="text-sm text-gray-900">Online</span>
                </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ device.hostname }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                {{ device.ip }}
                <button @click="copyToClipboard(device.ip)" class="ml-2 text-gray-300 hover:text-primary">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 01-2-2V5" />
                    </svg>
                </button>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ device.mac }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <a href="#" class="text-primary hover:text-purple-900">Remote</a>
            </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
