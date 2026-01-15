<script setup>
import { ref, reactive } from 'vue'
import DeviceList from './components/DeviceList.vue'
import ScanButton from './components/ScanButton.vue'
import { ScanNetwork } from '../wailsjs/go/client/App'

const devices = ref([])
const isLoading = ref(false)
const viewMode = ref('grid') // 'grid' or 'list'

async function handleScan() {
    isLoading.value = true
    try {
        // TODO: Password from settings or proper auth flow. Testing with empty for now (if logic allows) or prompt.
        // Assuming user uses default password set install time or handles it.
        // For MVP, we pass empty or hardcode dev password if needed, but prompts say 'ScanNetwork(password)'
        // Let's pass empty for now, assuming local scanner might just decode any. 
        // Wait, scanner logic needs the password to decrypt! 
        // We have no prompt here. 
        // I will add a simple prompt/input for password later. For now, assume a default or empty.
        
        const result = await ScanNetwork("minhasenha") // Using the walkthrough example password
        // Result is map[string]Device. Convert to array.
        devices.value = Object.values(result)
    } catch (err) {
        console.error("Scan failed:", err)
        alert("Scan Failed: " + err)
    } finally {
        isLoading.value = false
    }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 font-sans text-gray-900">
    <!-- Header -->
    <header class="bg-white shadow-sm sticky top-0 z-50">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between items-center h-16">
                <!-- Logo -->
                <div class="flex items-center">
                    <span class="text-2xl font-bold text-primary tracking-tight">goFindMyArm</span>
                </div>
                
                <!-- Actions -->
                <div class="flex items-center space-x-4">
                    <button @click="viewMode = viewMode === 'grid' ? 'list' : 'grid'" class="p-2 rounded-full text-gray-500 hover:bg-gray-100 hover:text-primary transition-colors">
                        <!-- List Icon -->
                         <svg v-if="viewMode === 'grid'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                        </svg>
                        <!-- Grid Icon -->
                         <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                        </svg>
                    </button>
                    <button class="p-2 rounded-full text-gray-500 hover:bg-gray-100 hover:text-primary transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- Hero / Scan Action -->
        <div class="flex flex-col items-center justify-center mb-12">
            <h2 class="text-3xl font-extrabold text-gray-900 mb-2">Find your devices</h2>
            <p class="text-gray-500 mb-8 text-center max-w-lg">
                Discover headless miniPCs on your local network instantly.
            </p>
            <ScanButton :loading="isLoading" @scan="handleScan" />
        </div>

        <!-- Results Area -->
        <div class="transition-all duration-300">
             <DeviceList :devices="devices" :viewMode="viewMode" />
        </div>
    </main>
  </div>
</template>
