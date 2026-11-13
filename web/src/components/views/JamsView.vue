<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-5">
    <h1 class="text-2xl font-extrabold text-white">Monthly Game Jams</h1>
    <div
      v-if="!overview.jams || overview.jams.length === 0"
      class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
    >
      No game jams posted yet.
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="jam in overview.jams"
        :key="jam.id"
        class="bg-[#242336] rounded-3xl overflow-hidden border border-white/[0.08] flex flex-col"
      >
        <img :src="jam.cover_image_url || '/inspo/sunset-studio.jpg'" class="h-48 w-full object-cover" />
        <div class="p-6 flex flex-col gap-3">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold uppercase tracking-wider text-[#ff4d8d]">{{ jam.month_label }}</span>
            <span class="px-2.5 py-0.5 rounded-full bg-[#00f5d4]/15 text-[#00f5d4] text-xs font-bold uppercase">
              {{ jam.status }}
            </span>
          </div>
          <h2 class="text-xl font-extrabold text-white">{{ jam.title }}</h2>
          <div class="text-xs font-mono text-[#ffd166]">THEME: {{ jam.theme }}</div>
          <p class="text-xs text-[#9e9eb8] leading-relaxed whitespace-pre-line">{{ jam.description }}</p>
          <div class="flex items-center justify-between pt-3 border-t border-white/[0.06] text-xs">
            <span class="text-[#8e8ea9]">{{ jam.participants_count }} Jammers • {{ jam.submissions_count }} Entries</span>
            <a
              :href="jam.itch_url"
              target="_blank"
              class="px-4 py-2 rounded-xl bg-[#6c5ce7] text-white font-bold"
            >
              Itch.io Jam Page →
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
