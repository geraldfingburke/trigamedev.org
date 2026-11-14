<script setup>
import { X, Trophy, ExternalLink } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { selectedWinnersJam } = useTriGameDevStore()

function medalBadgeClass(medal) {
  switch ((medal || '').toLowerCase()) {
    case 'gold':
      return 'bg-[#eab308] text-black border-[#eab308]'
    case 'silver':
      return 'bg-zinc-300 text-black border-zinc-300'
    default:
      return 'bg-[#b45309] text-white border-[#b45309]'
  }
}
</script>

<template>
  <div
    v-if="selectedWinnersJam"
    @click.self="selectedWinnersJam = null"
    class="fixed inset-0 z-50 bg-black/85 flex items-center justify-center p-4"
  >
    <div class="bg-[#121215] border border-zinc-700 max-w-xl w-full p-6 flex flex-col gap-5">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <div class="flex items-center gap-2">
          <Trophy class="w-5 h-5 text-[#ef4444]" />
          <div>
            <div class="text-[10px] font-mono uppercase tracking-widest text-[#ef4444]">WINNERS CIRCLE</div>
            <h3 class="text-base font-extrabold text-white uppercase">{{ selectedWinnersJam.title }}</h3>
          </div>
        </div>
        <button @click="selectedWinnersJam = null" class="text-zinc-400 hover:text-white cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>

      <div
        v-if="!selectedWinnersJam.winners || selectedWinnersJam.winners.length === 0"
        class="p-8 bg-[#09090b] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-400"
      >
        No ranked results recorded for this jam yet.
      </div>

      <div v-else class="flex flex-col gap-3">
        <div
          v-for="winner in selectedWinnersJam.winners"
          :key="winner.rank"
          class="p-4 bg-[#09090b] border border-zinc-800 flex items-center justify-between gap-4"
        >
          <div class="flex items-center gap-4 min-w-0">
            <img
              :src="winner.cover_url || '/brand/filler.png'"
              :alt="winner.title"
              class="w-16 h-16 object-cover border border-zinc-700 bg-black shrink-0"
            />
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <span
                  class="px-2 py-0.5 text-[10px] font-mono font-bold uppercase border"
                  :class="medalBadgeClass(winner.medal)"
                >
                  {{ winner.medal }} // #{{ winner.rank }}
                </span>
              </div>
              <h4 class="font-bold text-base text-white uppercase truncate mt-1">{{ winner.title }}</h4>
              <p class="text-xs font-mono text-zinc-400 truncate">BY {{ winner.authors }}</p>
            </div>
          </div>

          <a
            :href="winner.game_url || selectedWinnersJam.itch_url"
            target="_blank"
            class="px-3.5 py-2 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider flex items-center gap-1 shrink-0 transition"
          >
            <span>Play</span>
            <ExternalLink class="w-3.5 h-3.5" />
          </a>
        </div>
      </div>

      <div class="flex items-center justify-between pt-2 border-t border-zinc-800">
        <a
          :href="selectedWinnersJam.itch_url + '/results'"
          target="_blank"
          class="text-xs font-mono uppercase text-zinc-400 hover:text-white"
        >
          Full Itch.io Results Page ->
        </a>
        <button
          @click="selectedWinnersJam = null"
          class="px-4 py-2 bg-zinc-900 border border-zinc-700 hover:border-white text-white text-xs font-bold uppercase cursor-pointer"
        >
          Close
        </button>
      </div>
    </div>
  </div>
</template>
