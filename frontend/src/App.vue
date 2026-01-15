<script setup>
import { ref, onMounted } from 'vue'
import DeviceList from './components/DeviceList.vue'
import ScanButton from './components/ScanButton.vue'
import SettingsModal from './components/SettingsModal.vue'
import { ScanNetwork, RebootDevice } from '../wailsjs/go/client/App'

const devices = ref([])
const isLoading = ref(false)
const viewMode = ref('grid') // 'grid' or 'list'
const showSettings = ref(false)
const appPassword = ref('')

// Load password from local storage on mount
onMounted(() => {
    const saved = localStorage.getItem('gofindmyarm-password')
    if (saved) {
        appPassword.value = saved
    } else {
        // First time? open settings
        showSettings.value = true
    }
})

function saveSettings(newPassword) {
    appPassword.value = newPassword
    localStorage.setItem('gofindmyarm-password', newPassword)
}

async function handleReboot(ip) {
    if (!confirm(`Are you sure you want to reboot ${ip}?`)) {
        return
    }
    
    if (!appPassword.value) {
        alert("Please set a password in settings first.")
        showSettings.value = true
        return
    }

    try {
        await RebootDevice(ip, appPassword.value)
        alert("Reboot command sent successfully!")
    } catch (err) {
        console.error("Reboot failed:", err)
        alert("Reboot Failed: " + err)
    }
}

async function handleScan() {
    if (!appPassword.value) {
        // If no password, prompt settings
        alert("Please set a password first.")
        showSettings.value = true
        return
    }

    isLoading.value = true
    try {
        const result = await ScanNetwork(appPassword.value)
        // Result is map[string]Device. Convert to array.
        devices.value = Object.values(result)
        
        if (devices.value.length === 0) {
           console.log("No devices found (check password?)")
        }

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
                    <button @click="viewMode = viewMode === 'grid' ? 'list' : 'grid'" class="p-2 rounded-full text-gray-500 hover:bg-gray-100 hover:text-primary transition-colors" title="Toggle View">
                        <!-- List Icon -->
                         <svg v-if="viewMode === 'grid'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                        </svg>
                        <!-- Grid Icon -->
                         <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                        </svg>
                    </button>
                    <button @click="showSettings = true" class="p-2 rounded-full text-gray-500 hover:bg-gray-100 hover:text-primary transition-colors" title="Settings">
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
                Discover headless ARMs/MiniPCs/Iots on your local network instantly.
            </p>
            <div class="mb-8">
                <a href="https://github.com/apavanello/goFindMyArm" target="_blank" class="text-sm text-primary hover:underline flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" /></svg>
                    View on GitHub
                </a>
            </div>
            <ScanButton :loading="isLoading" @scan="handleScan" />
        </div>

        <!-- Results Area -->
        <div class="transition-all duration-300">
             <DeviceList :devices="devices" :viewMode="viewMode" @reboot="handleReboot" />
        </div>
    </main>

    <!-- Settings Modal -->
    <SettingsModal 
        :isOpen="showSettings" 
        :currentPassword="appPassword" 
        @close="showSettings = false"
        @save="saveSettings"
    />
  </div>
</template>
