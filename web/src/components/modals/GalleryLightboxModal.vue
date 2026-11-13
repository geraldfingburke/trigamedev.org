<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { selectedLightbox } = useTriGameDevStore()
</script>

<template>
  <div
    v-if="selectedLightbox"
    @click.self="selectedLightbox = null"
    class="fixed inset-0 z-50 bg-black/90 flex items-center justify-center p-4"
  >
    <div class="bg-[#121215] border border-zinc-700 max-w-2xl w-full overflow-hidden">
      <img :src="selectedLightbox.image_url" class="w-full max-h-[65vh] object-cover border-b border-zinc-800" />
      <div class="p-5 flex items-center justify-between">
        <div>
          <span class="text-xs font-mono uppercase text-[#ef4444]">{{ selectedLightbox.event_label }}</span>
          <h3 class="text-lg font-bold text-white uppercase">{{ selectedLightbox.title }}</h3>
          <p class="text-xs text-zinc-400 mt-0.5">{{ selectedLightbox.caption }}</p>
        </div>
        <button
          @click="selectedLightbox = null"
          class="px-4 py-2 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase cursor-pointer"
        >
          Close
        </button>
      </div>
    </div>
  </div>
</template>
