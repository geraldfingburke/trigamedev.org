<script setup>
import { Sparkles } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { toastMessage } = useTriGameDevStore()
</script>

<template>
  <transition name="fade">
    <div
      v-if="toastMessage"
      class="fixed bottom-6 right-6 z-50 flex items-center gap-3 px-5 py-3.5 rounded-2xl bg-[#252338]/95 border border-[#ff4d8d]/50 text-sm font-medium shadow-2xl shadow-[#ff4d8d]/20 backdrop-blur-xl"
    >
      <Sparkles class="w-4 h-4 text-[#ff4d8d] shrink-0" />
      <span>{{ toastMessage }}</span>
    </div>
  </transition>
</template>
