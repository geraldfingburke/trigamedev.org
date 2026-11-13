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
            Official Rules (48 Hours)
          </h2>
          <ul class="text-xs text-zinc-300 space-y-2 leading-relaxed">
            <li><strong class="text-white">Solo or Team:</strong> Go solo or build with a team.</li>
            <li><strong class="text-white">Any Engine or Framework:</strong> All engines and frameworks are allowed.</li>
            <li><strong class="text-white">Premade Assets Allowed:</strong> Openly available premade assets are allowed if licenses are respected.</li>
            <li><strong class="text-white">Honor the Theme:</strong> Interpret the theme however you like, but do not submit a pre-made game.</li>
            <li><strong class="text-[#ef4444]">Strictly No NSFW:</strong> Keep submissions welcoming for the whole community.</li>
          </ul>
        </div>

        <div class="bg-[#09090b] border border-zinc-800 p-5 flex flex-col gap-3">
          <h2 class="text-sm font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
            Ratings Criteria
          </h2>
          <ul class="text-xs text-zinc-300 space-y-2 leading-relaxed">
            <li><strong class="text-white">Theme:</strong> How well does the game adhere to or interpret the theme?</li>
            <li><strong class="text-white">Engagement:</strong> How well does the game capture attention and keep players engaged?</li>
            <li><strong class="text-white">Sound:</strong> How fitting and effective are the music and sound effects?</li>
            <li><strong class="text-white">Art:</strong> How cohesive and appealing is the visual presentation?</li>
            <li><strong class="text-white">Impression:</strong> Your personal, gut-feeling score for the entry.</li>
          </ul>
        </div>

        <div class="bg-[#09090b] border border-zinc-800 p-5 flex flex-col justify-between gap-4">
          <div class="flex flex-col gap-3">
            <h2 class="text-sm font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
              What Do I Get for Winning?
            </h2>
            <p class="text-sm font-bold text-[#ef4444] uppercase">
              Glory! Praise! Admiration!
            </p>
            <p class="text-xs text-zinc-400 leading-relaxed">
              This is a humble community jam meant to encourage creative exploration and experimentation. You get the satisfaction of finishing a build and playing what everyone else made.
            </p>
          </div>
          <img
            src="/brand/tres.png"
            alt="Una Substantia, Tres Personae"
            class="h-6 object-contain object-left opacity-90"
          />
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
                THEME: {{ jam.theme }}
              </div>
              <p class="text-xs text-zinc-400 leading-relaxed whitespace-pre-line">{{ jam.description }}</p>
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
