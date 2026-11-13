<script setup>
import { Terminal } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { toastMessage } = useTriGameDevStore()
</script>

<template>
  <transition name="fade">
    <div
      v-if="toastMessage"
      class="fixed bottom-6 right-6 z-50 flex items-center gap-3 px-5 py-3.5 bg-[#09090b] border border-white text-xs font-mono uppercase tracking-wider text-white"
    >
      <Terminal class="w-4 h-4 text-[#ef4444] shrink-0" />
      <span>{{ toastMessage }}</span>
    </div>
  </transition>
</template>
