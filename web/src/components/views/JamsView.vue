<script setup>
import { Trophy, ExternalLink } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview, selectedWinnersJam, activeSlug, slugify, navigateTo } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-8">
    <!-- Top Section Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
      <h1 class="text-2xl sm:text-3xl font-extrabold text-white uppercase">Jams</h1>
      <a
        href="https://discord.gg/6jdhMV6sqt"
        target="_blank"
        class="px-5 py-2.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider w-fit transition"
      >
        Join Jam on Discord ->
      </a>
    </div>

    <!-- Cards for "Monthly Game Jam" and "Other Jams" -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Monthly Game Jam Card (Spans 2 columns, includes official rules & ratings) -->
      <div class="lg:col-span-2 border border-zinc-800 bg-[#121215] p-6 flex flex-col gap-5">
        <div class="border-b border-zinc-800 pb-3">
          <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">48-HOUR COMMUNITY JAM</div>
          <h2 class="text-xl font-extrabold text-white uppercase mt-1">Monthly Game Jam</h2>
          <p class="text-xs text-zinc-300 mt-1">
            The jam runs for 48 hours. Incorporation of the theme is mandatory, but your interpretation can be whatever you want it to be.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-[#09090b] border border-zinc-800 p-4 flex flex-col gap-2">
            <h3 class="text-xs font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
              The Rules
            </h3>
            <ul class="text-xs text-zinc-300 space-y-1.5 leading-relaxed">
              <li>• NSFW content must be flagged</li>
              <li>• Your submission must be built during the duration of the jam</li>
              <li>• You can use any asset you are legally permitted to use</li>
              <li>• Gen-AI is permitted, lazy Gen-AI is not</li>
              <li>• No spam submissions, your submission must be unique to this jam</li>
              <li>• An interpretation of the theme should be incorporated into your submission</li>
              <li>• An explanation of your interpretation should be included in your submission</li>
              <li>• Respect yourself, respect your fellow jammers, respect the jam</li>
            </ul>
          </div>

          <div class="flex flex-col gap-4">
            <div class="bg-[#09090b] border border-zinc-800 p-4 flex flex-col gap-2">
              <h3 class="text-xs font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
                The Ratings
              </h3>
              <ul class="text-xs text-zinc-300 space-y-1.5 leading-relaxed">
                <li><strong class="text-white">Theme:</strong> How well did the game fit the theme?</li>
                <li><strong class="text-white">Engagement:</strong> Did it draw you in? Make you want to finish it?</li>
                <li><strong class="text-white">Sound:</strong> Does the game have good/interesting sound?</li>
                <li><strong class="text-white">Art:</strong> Does the game look good/interesting/strangeinagoodway?</li>
                <li><strong class="text-white">Impression:</strong> Does the game stick with you after you finish playing?</li>
              </ul>
            </div>

            <div class="bg-[#09090b] border border-zinc-800 p-4 flex flex-col gap-1.5">
              <h3 class="text-xs font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-1.5">
                What Do I Get for Winning?
              </h3>
              <p class="text-xs text-zinc-300">1: Nothing of value</p>
              <p class="text-xs text-zinc-300">2: Bragging rights, a sense of accomplishment, and a cool role on our Discord</p>
              <p class="text-[11px] text-zinc-400">The two are not mutually exclusive. Make of them what you will.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Other Jams Card -->
      <div class="border border-zinc-800 bg-[#121215] p-6 flex flex-col justify-between gap-4">
        <div class="flex flex-col gap-3">
          <div class="border-b border-zinc-800 pb-3">
            <div class="text-[11px] font-mono uppercase tracking-widest text-zinc-400">SPECIAL &amp; EXTERNAL</div>
            <h2 class="text-xl font-extrabold text-white uppercase mt-1">Other Jams</h2>
          </div>
          <p class="text-xs text-zinc-300 leading-relaxed">
            Alongside our 48-hour Monthly Game Jam, TriGameDev members team up for regional and global game jams including East Coast Game Jam, Ludum Dare, Global Game Jam, and seasonal community events.
          </p>
        </div>

        <div class="p-4 bg-[#09090b] border border-zinc-800 flex flex-col gap-2">
          <div class="text-[11px] font-mono uppercase text-[#ef4444]">FIND A TEAM</div>
          <p class="text-xs text-zinc-400">
            Jump into our Discord channels to form a group, share progress, or join an upcoming external jam.
          </p>
        </div>
      </div>
    </div>

    <!-- Past Jams (Compact 3-Column Scrollable Grid with Winners Circle Link) -->
    <div class="flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-2">
        <h2 class="text-lg font-bold text-white uppercase tracking-wider">Past Jams</h2>
        <span class="text-xs font-mono text-zinc-400 uppercase">
          {{ overview.jams?.length || 0 }} ARCHIVED
        </span>
      </div>

      <div
        v-if="!overview.jams || overview.jams.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No jams posted yet.
      </div>

      <div v-else class="max-h-[640px] overflow-y-auto pr-1">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="jam in overview.jams"
            :key="jam.id"
            :id="'slug-' + (jam.slug || slugify(jam.title))"
            class="bg-[#121215] border flex flex-col justify-between transition"
            :class="
              (jam.slug || slugify(jam.title)) === activeSlug
                ? 'border-[#ef4444] ring-1 ring-[#ef4444]'
                : 'border-zinc-800'
            "
          >
            <div>
              <img
                :src="jam.cover_image_url || '/brand/filler.png'"
                class="h-32 w-full object-cover border-b border-zinc-800 bg-black"
              />
              <div class="p-4 flex flex-col gap-2">
                <div class="flex items-center justify-between">
                  <span class="text-[11px] font-mono uppercase tracking-wider text-[#ef4444]">{{ jam.month_label }}</span>
                  <span class="px-1.5 py-0.5 bg-zinc-900 border border-zinc-700 text-zinc-300 text-[10px] font-mono uppercase">
                    {{ jam.jam_type === 'other' ? 'Other Jam' : jam.status }}
                  </span>
                </div>
                <h3
                  @click="navigateTo('jams', jam.slug || slugify(jam.title))"
                  class="text-sm font-extrabold text-white hover:text-[#ef4444] uppercase truncate cursor-pointer transition"
                  :title="jam.title"
                >
                  {{ jam.title }}
                </h3>
                <div class="text-[11px] font-mono text-zinc-300 bg-[#09090b] border border-zinc-800 px-2.5 py-1 truncate">
                  THEME: <span class="text-white font-bold">{{ jam.theme }}</span>
                </div>
                <div class="text-[11px] font-mono text-zinc-400">
                  {{ jam.participants_count }} JAMMERS // {{ jam.submissions_count }} ENTRIES
                </div>
              </div>
            </div>

            <div class="p-3 border-t border-zinc-800 grid grid-cols-2 gap-2 text-[11px]">
              <button
                @click="selectedWinnersJam = jam"
                class="py-1.5 px-2 bg-[#09090b] hover:bg-zinc-900 text-white border border-zinc-700 font-bold uppercase tracking-wider flex items-center justify-center gap-1 cursor-pointer transition"
              >
                <Trophy class="w-3 h-3 text-[#ef4444]" />
                <span>Winners Circle</span>
              </button>
              <a
                :href="jam.itch_url"
                target="_blank"
                class="py-1.5 px-2 bg-white hover:bg-zinc-200 text-black font-bold uppercase tracking-wider flex items-center justify-center gap-1 transition"
              >
                <span>Itch.io</span>
                <ExternalLink class="w-3 h-3" />
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
