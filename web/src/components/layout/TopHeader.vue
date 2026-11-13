<script setup>
import { computed } from 'vue'
import { Search } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { activeSection, searchQuery, currentUser } = useTriGameDevStore()

const tabs = computed(() => [
  { id: 'explore', label: 'Explore' },
  { id: 'events', label: 'Events' },
  { id: 'news', label: 'Announcements' },
  { id: 'jams', label: 'Monthly Jam' },
  { id: 'showcase', label: 'Showcase' },
  { id: 'media', label: 'Gallery & Videos' },
  { id: 'keys', label: 'Key Vault' },
  { id: 'library', label: 'Ebook Library' },
  ...(currentUser.value?.is_og_admin ? [{ id: 'cms', label: '⚡ OG CMS' }] : [])
])
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Top Search Bar -->
    <header class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4">
      <div class="relative flex-1 max-w-xl">
        <Search class="w-4 h-4 text-[#8e8ea9] absolute left-4 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search TriGameDev games, keys, ebooks, or events..."
          class="w-full pl-11 pr-4 py-3 rounded-2xl bg-[#252438] border border-white/[0.06] text-sm text-white placeholder-[#7c7c9a] focus:outline-none focus:border-[#8b5cf6] transition"
        />
      </div>
    </header>

    <!-- Horizontal Category Underline Tabs -->
    <div class="flex items-center justify-between border-b border-white/[0.07] pb-1 overflow-x-auto gap-6">
      <div class="flex items-center gap-7 text-sm font-medium shrink-0">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeSection = tab.id"
          class="pb-3 relative transition cursor-pointer whitespace-nowrap"
          :class="
            activeSection === tab.id
              ? 'text-white font-semibold'
              : 'text-[#8e8ea9] hover:text-white'
          "
        >
          {{ tab.label }}
          <span
            v-if="activeSection === tab.id"
            class="absolute bottom-0 left-0 right-0 h-[3px] rounded-full bg-gradient-to-r from-[#6c5ce7] via-[#8b5cf6] to-[#ff4d8d]"
          ></span>
        </button>
      </div>
    </div>
  </div>
</template>
