<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview, selectedLightbox } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-7">
    <div>
      <h1 class="text-2xl font-extrabold text-white mb-4">Community Photo Gallery</h1>
      <div
        v-if="!overview.gallery || overview.gallery.length === 0"
        class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
      >
        No gallery photos uploaded yet.
      </div>
      <div v-else class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div
          v-for="img in overview.gallery"
          :key="img.id"
          @click="selectedLightbox = img"
          class="group rounded-3xl overflow-hidden bg-[#242336] border border-white/[0.07] cursor-pointer"
        >
          <img :src="img.image_url" class="h-44 w-full object-cover group-hover:scale-105 transition" />
          <div class="p-3">
            <div class="text-[10px] font-bold text-[#ff4d8d]">{{ img.event_label }}</div>
            <div class="text-xs font-bold text-white truncate">{{ img.title }}</div>
          </div>
        </div>
      </div>
    </div>

    <div>
      <h2 class="text-2xl font-extrabold text-white mb-4">Talks & Workshop Videos</h2>
      <div
        v-if="!overview.videos || overview.videos.length === 0"
        class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
      >
        No videos archived yet.
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <a
          v-for="vid in overview.videos"
          :key="vid.id"
          :href="vid.video_url"
          target="_blank"
          class="bg-[#242336] rounded-3xl overflow-hidden border border-white/[0.07] hover:border-[#8b5cf6] transition flex flex-col"
        >
          <div class="relative h-40">
            <img :src="vid.thumbnail_url || '/inspo/sunset-studio.jpg'" class="w-full h-full object-cover" />
            <span v-if="vid.duration" class="absolute bottom-2 right-2 px-2 py-0.5 rounded bg-black/80 text-[11px] font-mono text-white">
              {{ vid.duration }}
            </span>
          </div>
          <div class="p-4">
            <h3 class="font-bold text-sm text-white">{{ vid.title }}</h3>
            <p class="text-xs text-[#a78bfa] mt-1">{{ vid.speaker_name }}</p>
            <p class="text-xs text-[#9e9eb8] mt-1 line-clamp-2">{{ vid.description }}</p>
          </div>
        </a>
      </div>
    </div>
  </div>
</template>
