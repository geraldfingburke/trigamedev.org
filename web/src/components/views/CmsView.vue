<script setup>
import { RefreshCw, Terminal } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useCms } from '../../composables/useCms'

const { currentUser, overview } = useTriGameDevStore()
const {
  scrapingItch,
  cmsForm,
  autofillFromItch,
  publishCMSItem,
  deleteCMSItem,
  triggerDiscordSync
} = useCms()
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
      <div>
        <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[#ef4444] text-[10px] font-mono uppercase tracking-widest">
          OG DISCORD ROLE VERIFIED
        </span>
        <h1 class="text-2xl font-extrabold text-white uppercase mt-1">TriGameDev OG Content Management System</h1>
      </div>
      <button
        @click="triggerDiscordSync"
        class="px-4 py-2.5 bg-[#121215] hover:bg-zinc-800 border border-zinc-700 text-xs font-bold uppercase tracking-wider text-white flex items-center gap-2 cursor-pointer transition"
      >
        <RefreshCw class="w-4 h-4 text-[#ef4444]" />
        <span>Sync Discord Events & Announcements</span>
      </button>
    </div>

    <div
      v-if="!currentUser?.is_og_admin"
      class="p-6 bg-[#121215] border border-[#ef4444] text-sm text-zinc-200"
    >
      You must be signed in via Discord and hold the <strong class="text-white">OG Role</strong> in the TriGameDev server to access the CMS.
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Publish & Itch.io Importer Form -->
      <div class="bg-[#121215] p-6 border border-zinc-800 flex flex-col gap-4">
        <h2 class="text-base font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
          Publish New Community Entry
        </h2>

        <!-- One-Click Itch.io Scraper Box -->
        <div class="p-4 bg-[#09090b] border border-zinc-700 flex flex-col gap-2.5">
          <div class="flex items-center justify-between">
            <label class="text-xs font-mono uppercase text-white flex items-center gap-1.5">
              <Terminal class="w-3.5 h-3.5 text-[#ef4444]" />
              <span>One-Click Itch.io Jam / Game Importer</span>
            </label>
            <span class="text-[11px] font-mono text-zinc-500">PASTE LINK & AUTO-FILL</span>
          </div>
          <div class="flex flex-col sm:flex-row gap-2">
            <input
              v-model="cmsForm.url"
              type="url"
              placeholder="https://itch.io/jam/... or https://creator.itch.io/game"
              class="flex-1 px-3.5 py-2.5 bg-[#121215] border border-zinc-700 text-sm text-white placeholder-zinc-600 focus:outline-none focus:border-white"
            />
            <button
              type="button"
              @click="autofillFromItch"
              :disabled="scrapingItch"
              class="px-4 py-2.5 bg-white hover:bg-zinc-200 disabled:opacity-50 text-black text-xs font-bold uppercase tracking-wider whitespace-nowrap cursor-pointer transition"
            >
              {{ scrapingItch ? 'Scraping Itch.io...' : 'Auto-Fill from Itch.io' }}
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Target Section</label>
            <select
              v-model="cmsForm.section"
              class="w-full mt-1 px-3 py-2.5 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            >
              <option value="jams">Monthly Game Jam</option>
              <option value="showcase">Game Showcase</option>
              <option value="news">News / Announcement</option>
              <option value="events">Event / Meetup</option>
            </select>
          </div>
          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Tag / Engine / Month</label>
            <input
              v-model="cmsForm.tag"
              type="text"
              placeholder="Auto-filled for Itch links"
              class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
          </div>
        </div>

        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Title (Optional if Itch.io URL is provided above)</label>
          <input
            v-model="cmsForm.title"
            type="text"
            placeholder="Leave blank to auto-scrape on publish, or click Auto-Fill to preview"
            class="w-full mt-1 px-3.5 py-2.5 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
        </div>

        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Subtitle / Location / Theme / Developer</label>
          <input
            v-model="cmsForm.subtitle"
            type="text"
            class="w-full mt-1 px-3.5 py-2.5 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
        </div>

        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Cover Image URL</label>
          <input
            v-model="cmsForm.image_url"
            type="text"
            class="w-full mt-1 px-3.5 py-2.5 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
        </div>

        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Body / Description</label>
          <textarea
            v-model="cmsForm.description"
            rows="3"
            class="w-full mt-1 px-3.5 py-2.5 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          ></textarea>
        </div>

        <button
          @click="publishCMSItem"
          class="py-3 bg-[#ef4444] hover:bg-red-500 text-white font-extrabold text-xs uppercase tracking-wider cursor-pointer transition"
        >
          Publish Entry
        </button>
      </div>

      <!-- Manage Live Entries -->
      <div class="bg-[#121215] p-6 border border-zinc-800 flex flex-col gap-4">
        <h2 class="text-base font-bold text-white uppercase tracking-wider border-b border-zinc-800 pb-2">
          Manage Live Entries
        </h2>
        <div class="flex flex-col gap-2.5 max-h-[540px] overflow-y-auto pr-1">
          <div
            v-for="item in overview.news"
            :key="'n-' + item.id"
            class="p-3 bg-[#09090b] border border-zinc-800 flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="font-mono text-zinc-400 mr-2">[NEWS]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('news', item.id)" class="text-[#ef4444] font-mono uppercase hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>

          <div
            v-for="item in overview.events"
            :key="'e-' + item.id"
            class="p-3 bg-[#09090b] border border-zinc-800 flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="font-mono text-zinc-400 mr-2">[EVENT]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('events', item.id)" class="text-[#ef4444] font-mono uppercase hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>

          <div
            v-for="item in overview.jams"
            :key="'j-' + item.id"
            class="p-3 bg-[#09090b] border border-zinc-800 flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="font-mono text-zinc-400 mr-2">[JAM]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('jams', item.id)" class="text-[#ef4444] font-mono uppercase hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>

          <div
            v-for="item in overview.showcase"
            :key="'s-' + item.id"
            class="p-3 bg-[#09090b] border border-zinc-800 flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="font-mono text-zinc-400 mr-2">[SHOWCASE]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('showcase', item.id)" class="text-[#ef4444] font-mono uppercase hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
