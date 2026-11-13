<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { selectedLightbox } = useTriGameDevStore()
</script>

<template>
  <div
    v-if="selectedLightbox"
    @click.self="selectedLightbox = null"
    class="fixed inset-0 z-50 bg-black/85 backdrop-blur-md flex items-center justify-center p-4"
  >
    <div class="bg-[#232234] border border-white/15 rounded-3xl max-w-2xl w-full overflow-hidden shadow-2xl">
      <img :src="selectedLightbox.image_url" class="w-full max-h-[65vh] object-cover" />
      <div class="p-5 flex items-center justify-between">
        <div>
          <span class="text-xs font-bold text-[#ff4d8d]">{{ selectedLightbox.event_label }}</span>
          <h3 class="text-lg font-bold text-white">{{ selectedLightbox.title }}</h3>
          <p class="text-xs text-[#9e9eb8] mt-0.5">{{ selectedLightbox.caption }}</p>
        </div>
        <button
          @click="selectedLightbox = null"
          class="px-4 py-2 rounded-xl bg-white/10 text-white text-xs font-bold cursor-pointer"
        >
          Close
        </button>
      </div>
    </div>
  </div>
</template>
