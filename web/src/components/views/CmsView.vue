<script setup>
import { RefreshCw, Sparkles } from 'lucide-vue-next'
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
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <span class="px-2.5 py-1 rounded-full bg-[#ff9f43]/20 text-[#ff9f43] text-[11px] font-bold">
          OG DISCORD ROLE VERIFIED
        </span>
        <h1 class="text-2xl font-extrabold text-white mt-1">TriGameDev OG Content Management System</h1>
      </div>
      <button
        @click="triggerDiscordSync"
        class="px-4 py-2.5 rounded-xl bg-[#252438] hover:bg-[#312f4a] border border-[#8b5cf6]/40 text-xs font-bold text-white flex items-center gap-2 cursor-pointer"
      >
        <RefreshCw class="w-4 h-4 text-[#00f5d4]" />
        <span>Sync Discord Events & Announcements</span>
      </button>
    </div>

    <div
      v-if="!currentUser?.is_og_admin"
      class="p-6 rounded-3xl bg-[#242336] border border-[#ff9f43]/40 text-sm text-[#ff9f43]"
    >
      You must be signed in via Discord and hold the <strong>OG Role</strong> in the TriGameDev server to access the CMS.
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Publish & Itch.io Importer Form -->
      <div class="bg-[#242336] rounded-3xl p-6 border border-white/[0.07] flex flex-col gap-4">
        <h2 class="text-lg font-bold text-white">Publish New Community Entry</h2>

        <!-- One-Click Itch.io Scraper Box -->
        <div class="p-4 rounded-2xl bg-[#181726] border border-[#8b5cf6]/40 flex flex-col gap-2.5">
          <div class="flex items-center justify-between">
            <label class="text-xs font-bold text-[#a78bfa] flex items-center gap-1.5">
              <Sparkles class="w-3.5 h-3.5 text-[#ff4d8d]" />
              <span>One-Click Itch.io Jam / Game Importer</span>
            </label>
            <span class="text-[11px] text-[#8e8ea9]">Paste a link & publish directly or auto-fill</span>
          </div>
          <div class="flex flex-col sm:flex-row gap-2">
            <input
              v-model="cmsForm.url"
              type="url"
              placeholder="https://itch.io/jam/... or https://creator.itch.io/game"
              class="flex-1 px-3.5 py-2.5 rounded-xl bg-[#12121b] border border-white/10 text-sm text-white placeholder-[#6c6c8a]"
            />
            <button
              type="button"
              @click="autofillFromItch"
              :disabled="scrapingItch"
              class="px-4 py-2.5 rounded-xl bg-[#6c5ce7] hover:bg-[#5a4ad1] disabled:opacity-50 text-white text-xs font-bold whitespace-nowrap cursor-pointer transition"
            >
              {{ scrapingItch ? 'Scraping Itch.io...' : 'Auto-Fill from Itch.io' }}
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="text-xs text-[#8e8ea9]">Target Section</label>
            <select
              v-model="cmsForm.section"
              class="w-full mt-1 px-3 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
            >
              <option value="jams">Monthly Game Jam</option>
              <option value="showcase">Game Showcase</option>
              <option value="news">News / Announcement</option>
              <option value="events">Event / Meetup</option>
            </select>
          </div>
          <div>
            <label class="text-xs text-[#8e8ea9]">Tag / Engine / Month</label>
            <input
              v-model="cmsForm.tag"
              type="text"
              placeholder="Auto-filled for Itch links"
              class="w-full mt-1 px-3 py-2 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
            />
          </div>
        </div>

        <div>
          <label class="text-xs text-[#8e8ea9]">Title (Optional if Itch.io URL is provided above)</label>
          <input
            v-model="cmsForm.title"
            type="text"
            placeholder="Leave blank to auto-scrape on publish, or click Auto-Fill to preview"
            class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
          />
        </div>

        <div>
          <label class="text-xs text-[#8e8ea9]">Subtitle / Location / Theme / Developer</label>
          <input
            v-model="cmsForm.subtitle"
            type="text"
            class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
          />
        </div>

        <div>
          <label class="text-xs text-[#8e8ea9]">Cover Image URL</label>
          <input
            v-model="cmsForm.image_url"
            type="text"
            class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
          />
        </div>

        <div>
          <label class="text-xs text-[#8e8ea9]">Body / Description</label>
          <textarea
            v-model="cmsForm.description"
            rows="3"
            class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
          ></textarea>
        </div>

        <button
          @click="publishCMSItem"
          class="py-3 rounded-xl bg-gradient-to-r from-[#ff9f43] to-[#ff4d8d] text-[#14141d] font-extrabold text-sm cursor-pointer shadow-lg"
        >
          Publish Entry
        </button>
      </div>

      <!-- Manage Live Entries -->
      <div class="bg-[#242336] rounded-3xl p-6 border border-white/[0.07] flex flex-col gap-4">
        <h2 class="text-lg font-bold text-white">Manage Live Entries</h2>
        <div class="flex flex-col gap-2.5 max-h-[540px] overflow-y-auto pr-1">
          <div
            v-for="item in overview.news"
            :key="'n-' + item.id"
            class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="text-[#ff9f43] font-bold mr-2">[NEWS]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('news', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>

          <div
            v-for="item in overview.events"
            :key="'e-' + item.id"
            class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="text-[#00f5d4] font-bold mr-2">[EVENT]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('events', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>

          <div
            v-for="item in overview.jams"
            :key="'j-' + item.id"
            class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="text-[#ff4d8d] font-bold mr-2">[JAM]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('jams', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>

          <div
            v-for="item in overview.showcase"
            :key="'s-' + item.id"
            class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
          >
            <div class="truncate">
              <span class="text-[#a78bfa] font-bold mr-2">[SHOWCASE]</span>
              <span class="text-white font-semibold">{{ item.title }}</span>
            </div>
            <button @click="deleteCMSItem('showcase', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
