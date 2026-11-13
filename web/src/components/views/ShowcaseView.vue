<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { filteredShowcase } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-5">
    <h1 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3">
      Member Game Showcase
    </h1>
    <div
      v-if="filteredShowcase.length === 0"
      class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
    >
      No showcase builds published yet.
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <article
        v-for="game in filteredShowcase"
        :key="game.id"
        class="bg-[#121215] border border-zinc-800 flex flex-col justify-between"
      >
        <div>
          <div class="relative h-48 overflow-hidden bg-black border-b border-zinc-800">
            <img :src="game.thumbnail_url || '/brand/filler.png'" class="w-full h-full object-cover" />
          </div>
          <div class="p-4">
            <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[10px] font-mono uppercase text-[#ef4444]">
              {{ game.engine }}
            </span>
            <h3 class="font-bold text-base text-white uppercase mt-2">{{ game.title }}</h3>
            <p class="text-xs font-mono text-zinc-400">BY {{ game.developer_name }}</p>
            <p class="text-xs text-zinc-400 mt-2 leading-relaxed">{{ game.description }}</p>
          </div>
        </div>
        <div class="p-4 pt-0 grid grid-cols-2 gap-2">
          <a
            :href="game.play_url"
            target="_blank"
            class="py-2 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider text-center transition"
          >
            Play Game
          </a>
          <a
            :href="game.source_url"
            target="_blank"
            class="py-2 border border-zinc-700 hover:border-white text-zinc-300 hover:text-white text-xs font-bold uppercase tracking-wider text-center transition"
          >
            Source / Store
          </a>
        </div>
      </article>
    </div>
  </div>
</template>
