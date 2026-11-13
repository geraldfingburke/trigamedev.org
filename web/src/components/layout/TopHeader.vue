<script setup>
import { computed } from 'vue'
import { Search } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { activeSection, searchQuery, currentUser } = useTriGameDevStore()

const tabs = computed(() => [
  { id: 'explore', label: 'EXPLORE' },
  { id: 'jams', label: 'MONTHLY JAM' },
  { id: 'events', label: 'EVENTS' },
  { id: 'news', label: 'ANNOUNCEMENTS' },
  { id: 'showcase', label: 'SHOWCASE' },
  { id: 'media', label: 'GALLERY & VIDEO' },
  { id: 'keys', label: 'KEY VAULT' },
  { id: 'library', label: 'EBOOK LIBRARY' },
  ...(currentUser.value?.is_og_admin ? [{ id: 'cms', label: 'OG CMS [ADMIN]' }] : [])
])
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Top Brand Bar & Search -->
    <header class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 border-b border-zinc-800 pb-5">
      <button @click="activeSection = 'explore'" class="flex items-center gap-4 text-left cursor-pointer group">
        <img
          src="/brand/logo-large.png"
          alt="TRIGAMEDEV"
          class="h-11 sm:h-12 w-auto object-contain bg-black border border-zinc-800 px-2 py-1 group-hover:border-white transition"
        />
        <div class="hidden sm:block">
          <div class="font-mono text-[11px] uppercase tracking-widest text-zinc-400">
            Northeast TN &amp; Beyond // Indie Game Collective
          </div>
          <div class="font-display font-bold text-sm text-white tracking-wide uppercase">
            Solo game dev can be lonely. Let's be alone together.
          </div>
        </div>
      </button>

      <div class="relative w-full md:w-80">
        <Search class="w-4 h-4 text-zinc-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="FILTER GAMES, KEYS, EBOOKS..."
          class="w-full pl-10 pr-4 py-2.5 bg-[#09090b] border border-zinc-700 text-xs font-mono text-white placeholder-zinc-500 focus:outline-none focus:border-white transition"
        />
      </div>
    </header>

    <!-- Horizontal Category Tabs (Sharp industrial tabs, no gradients) -->
    <div class="flex items-center justify-between border-b border-zinc-800 overflow-x-auto">
      <div class="flex items-center gap-1 shrink-0">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeSection = tab.id"
          class="px-4 py-2.5 font-mono text-xs tracking-wider uppercase transition cursor-pointer whitespace-nowrap border-b-2"
          :class="
            activeSection === tab.id
              ? 'border-white text-white bg-zinc-900 font-bold'
              : 'border-transparent text-zinc-400 hover:text-white hover:bg-zinc-900/50'
          "
        >
          {{ tab.label }}
        </button>
      </div>
    </div>
  </div>
</template>
