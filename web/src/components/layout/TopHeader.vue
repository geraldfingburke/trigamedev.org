<script setup>
import { computed } from 'vue'
import { Search } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { activeSection, searchQuery, currentUser } = useTriGameDevStore()

const tabs = computed(() => [
  { id: 'explore', label: 'EXPLORE' },
  { id: 'jams', label: 'JAMS' },
  { id: 'events', label: 'EVENTS' },
  { id: 'news', label: 'ANNOUNCEMENTS' },
  { id: 'showcase', label: 'GAMES BY OUR PEOPLE' },
  { id: 'media', label: 'GALLERY & VIDEO' },
  ...(currentUser.value
    ? [
        { id: 'keys', label: 'KEY VAULT' },
        { id: 'library', label: 'EBOOK LIBRARY' }
      ]
    : []),
  ...(currentUser.value?.is_og_admin ? [{ id: 'cms', label: 'CMS' }] : [])
])
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Top Header & Search (No duplicate brand image; circle logo in NavRail stands on its own) -->
    <header class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 border-b border-zinc-800 pb-5">
      <button @click="activeSection = 'explore'" class="flex flex-col text-left cursor-pointer group">
        <div class="font-mono text-xs uppercase tracking-widest text-[#ef4444]">
          Indie Dev Collective
        </div>
        <div class="font-display font-bold text-base sm:text-lg text-white tracking-wide uppercase mt-0.5">
          Solo game dev can be lonely. Let's be alone together.
        </div>
      </button>

      <div class="relative w-full md:w-80">
        <Search class="w-4 h-4 text-zinc-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="What're you lookin' for?"
          class="w-full pl-10 pr-4 py-2.5 bg-[#09090b] border border-zinc-700 text-xs font-mono text-white placeholder-zinc-500 focus:outline-none focus:border-white transition"
        />
      </div>
    </header>

    <!-- Horizontal Category Tabs -->
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
