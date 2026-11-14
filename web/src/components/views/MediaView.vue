<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview, selectedLightbox, activeSlug, slugify } = useTriGameDevStore()
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Column 1: Photos -->
    <div class="flex flex-col gap-4">
      <h1 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3">
        Photos
      </h1>
      <div
        v-if="!overview.gallery || overview.gallery.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No photos uploaded yet.
      </div>
      <div v-else class="grid grid-cols-2 gap-4">
        <div
          v-for="img in overview.gallery"
          :key="img.id"
          :id="'slug-' + (img.slug || slugify(img.title))"
          @click="selectedLightbox = img"
          class="group overflow-hidden bg-[#121215] border hover:border-white cursor-pointer transition"
          :class="
            (img.slug || slugify(img.title)) === activeSlug
              ? 'border-[#ef4444] ring-1 ring-[#ef4444]'
              : 'border-zinc-800'
          "
        >
          <img :src="img.image_url" class="h-44 w-full object-cover border-b border-zinc-800 bg-black" />
          <div class="p-3">
            <div class="text-[10px] font-mono uppercase text-[#ef4444]">{{ img.event_label }}</div>
            <div class="text-xs font-bold text-white uppercase truncate mt-0.5">{{ img.title }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Column 2: Videos -->
    <div class="flex flex-col gap-4">
      <h2 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3">
        Videos
      </h2>
      <div
        v-if="!overview.videos || overview.videos.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No videos posted yet.
      </div>
      <div v-else class="flex flex-col gap-4">
        <a
          v-for="vid in overview.videos"
          :key="vid.id"
          :id="'slug-' + (vid.slug || slugify(vid.title))"
          :href="vid.video_url"
          target="_blank"
          class="bg-[#121215] border hover:border-white transition flex flex-col sm:flex-row"
          :class="
            (vid.slug || slugify(vid.title)) === activeSlug
              ? 'border-[#ef4444] ring-1 ring-[#ef4444]'
              : 'border-zinc-800'
          "
        >
          <div class="relative w-full sm:w-48 h-36 border-b sm:border-b-0 sm:border-r border-zinc-800 bg-black shrink-0">
            <img :src="vid.thumbnail_url || '/brand/filler.png'" class="w-full h-full object-cover" />
            <span
              v-if="vid.duration"
              class="absolute bottom-2 right-2 px-2 py-0.5 bg-black border border-zinc-700 text-[10px] font-mono text-white"
            >
              {{ vid.duration }}
            </span>
          </div>
          <div class="p-4 flex flex-col justify-center">
            <h3 class="font-bold text-sm text-white uppercase">{{ vid.title }}</h3>
            <p class="text-xs font-mono text-[#ef4444] mt-1">{{ vid.speaker_name }}</p>
            <p class="text-xs text-zinc-400 mt-1 line-clamp-2">{{ vid.description }}</p>
          </div>
        </a>
      </div>
    </div>
  </div>
</template>
