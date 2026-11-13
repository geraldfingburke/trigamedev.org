<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview, selectedLightbox } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-8">
    <div>
      <h1 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3 mb-4">
        Community Photo Gallery
      </h1>
      <div
        v-if="!overview.gallery || overview.gallery.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No gallery photos uploaded yet.
      </div>
      <div v-else class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div
          v-for="img in overview.gallery"
          :key="img.id"
          @click="selectedLightbox = img"
          class="group overflow-hidden bg-[#121215] border border-zinc-800 hover:border-white cursor-pointer transition"
        >
          <img :src="img.image_url" class="h-44 w-full object-cover border-b border-zinc-800" />
          <div class="p-3">
            <div class="text-[10px] font-mono uppercase text-[#ef4444]">{{ img.event_label }}</div>
            <div class="text-xs font-bold text-white uppercase truncate mt-0.5">{{ img.title }}</div>
          </div>
        </div>
      </div>
    </div>

    <div>
      <h2 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3 mb-4">
        Talks & Workshop Videos
      </h2>
      <div
        v-if="!overview.videos || overview.videos.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No videos archived yet.
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <a
          v-for="vid in overview.videos"
          :key="vid.id"
          :href="vid.video_url"
          target="_blank"
          class="bg-[#121215] border border-zinc-800 hover:border-white transition flex flex-col"
        >
          <div class="relative h-40 border-b border-zinc-800 bg-black">
            <img :src="vid.thumbnail_url || '/brand/filler.png'" class="w-full h-full object-cover" />
            <span
              v-if="vid.duration"
              class="absolute bottom-2 right-2 px-2 py-0.5 bg-black border border-zinc-700 text-[11px] font-mono text-white"
            >
              {{ vid.duration }}
            </span>
          </div>
          <div class="p-4">
            <h3 class="font-bold text-sm text-white uppercase">{{ vid.title }}</h3>
            <p class="text-xs font-mono text-[#ef4444] mt-1">{{ vid.speaker_name }}</p>
            <p class="text-xs text-zinc-400 mt-1 line-clamp-2">{{ vid.description }}</p>
          </div>
        </a>
      </div>
    </div>
  </div>
</template>
