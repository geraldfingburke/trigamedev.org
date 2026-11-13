<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { filteredShowcase } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-5">
    <h1 class="text-2xl font-extrabold text-white">Member Game Showcase</h1>
    <div
      v-if="filteredShowcase.length === 0"
      class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
    >
      No showcase builds published yet.
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <article
        v-for="game in filteredShowcase"
        :key="game.id"
        class="bg-[#242336] rounded-3xl p-4 border border-white/[0.07] flex flex-col justify-between gap-4"
      >
        <div class="relative h-48 rounded-2xl overflow-hidden">
          <img :src="game.thumbnail_url || '/inspo/sunset-studio.jpg'" class="w-full h-full object-cover" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-[#00f5d4]">{{ game.engine }}</span>
          <h3 class="font-bold text-base text-white mt-0.5">{{ game.title }}</h3>
          <p class="text-xs text-[#a78bfa]">by {{ game.developer_name }}</p>
          <p class="text-xs text-[#9e9eb8] mt-2">{{ game.description }}</p>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <a :href="game.play_url" target="_blank" class="py-2 rounded-xl bg-[#6c5ce7] text-white text-xs font-bold text-center">
            Play Game
          </a>
          <a :href="game.source_url" target="_blank" class="py-2 rounded-xl border border-white/20 text-white text-xs font-bold text-center">
            Source / Store
          </a>
        </div>
      </article>
    </div>
  </div>
</template>
