<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-8">
    <!-- Header Banner with Official Brandkit Jam Copy -->
    <div class="border border-zinc-800 bg-[#121215] p-6 sm:p-8 flex flex-col gap-6">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-800 pb-5">
        <div>
          <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">48-HOUR COMMUNITY JAM</div>
          <h1 class="text-2xl sm:text-3xl font-extrabold text-white uppercase mt-1">TriGameDev Monthly Game Jams</h1>
        </div>
        <a
          href="https://discord.gg/6jdhMV6sqt"
          target="_blank"
          class="px-5 py-2.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider w-fit transition"
        >
          Join Jam on Discord ->
        </a>
      </div>

      <!-- Rules, Ratings & Prizes Grid from tgd-mgj-template.html -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="bg-[#09090b] border border-zinc-800 p-5 flex flex-col gap-3">
          <h2 class="text-sm font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
            The Rules (48 Hours)
          </h2>
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

        <div class="bg-[#09090b] border border-zinc-800 p-5 flex flex-col gap-3">
          <h2 class="text-sm font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
            The Ratings
          </h2>
          <ul class="text-xs text-zinc-300 space-y-2 leading-relaxed">
            <li><strong class="text-white">Theme:</strong> How well did the game fit the theme?</li>
            <li><strong class="text-white">Engagement:</strong> Did it draw you in? Make you want to finish it?</li>
            <li><strong class="text-white">Sound:</strong> Does the game have good/interesting sound?</li>
            <li><strong class="text-white">Art:</strong> Does the game look good/interesting/strangeinagoodway?</li>
            <li><strong class="text-white">Impression:</strong> Does the game stick with you after you finish playing?</li>
          </ul>
        </div>

        <div class="bg-[#09090b] border border-zinc-800 p-5 flex flex-col justify-between gap-4">
          <div class="flex flex-col gap-3">
            <h2 class="text-sm font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
              What Do I Get for Winning?
            </h2>
            <p class="text-xs text-zinc-300">There are two ways to look at this:</p>
            <div class="text-xs text-zinc-300 space-y-2">
              <p><strong class="text-white">1:</strong> Nothing of value</p>
              <p><strong class="text-white">2:</strong> Bragging rights, a sense of accomplishment, and a cool role on our Discord</p>
            </div>
            <p class="text-xs text-zinc-400 leading-relaxed">
              The two are not mutually exclusive. Make of them what you will.
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Live & Archived Jam Entries -->
    <div class="flex flex-col gap-4">
      <h2 class="text-lg font-bold text-white uppercase tracking-wider">Published Jams</h2>

      <div
        v-if="!overview.jams || overview.jams.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No game jams posted yet. OG Admins can import an Itch.io jam URL in the CMS.
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <div
          v-for="jam in overview.jams"
          :key="jam.id"
          class="bg-[#121215] border border-zinc-800 flex flex-col justify-between"
        >
          <div>
            <img
              :src="jam.cover_image_url || '/brand/filler.png'"
              class="h-48 w-full object-cover border-b border-zinc-800 bg-black"
            />
            <div class="p-6 flex flex-col gap-3">
              <div class="flex items-center justify-between">
                <span class="text-xs font-mono uppercase tracking-wider text-[#ef4444]">{{ jam.month_label }}</span>
                <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-white text-[11px] font-mono uppercase">
                  {{ jam.status }}
                </span>
              </div>
              <h3 class="text-xl font-extrabold text-white uppercase">{{ jam.title }}</h3>
              <div class="text-xs font-mono text-zinc-300 bg-[#09090b] border border-zinc-800 px-3 py-1.5 w-fit">
                THEME: <span class="text-white font-bold">{{ jam.theme }}</span>
              </div>
            </div>
          </div>

          <div class="px-6 py-4 border-t border-zinc-800 flex items-center justify-between text-xs">
            <span class="font-mono text-zinc-400">
              {{ jam.participants_count }} JAMMERS // {{ jam.submissions_count }} ENTRIES
            </span>
            <a
              :href="jam.itch_url"
              target="_blank"
              class="px-4 py-2 bg-white hover:bg-zinc-200 text-black font-bold uppercase tracking-wider transition"
            >
              Itch.io Jam Page ->
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
