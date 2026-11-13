<script setup>
import { computed } from 'vue'
import {
  Compass,
  Flame,
  Calendar,
  Newspaper,
  Sparkles,
  Image as ImageIcon,
  KeyRound,
  BookOpen,
  ShieldAlert
} from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { activeSection, overview, currentUser } = useTriGameDevStore()

const navItems = computed(() => [
  { id: 'explore', label: 'Explore Hub', icon: Compass },
  { id: 'jams', label: 'Monthly Game Jam', icon: Flame },
  { id: 'events', label: 'Events & Meetups', icon: Calendar },
  { id: 'news', label: 'News & Announcements', icon: Newspaper },
  { id: 'showcase', label: 'Member Showcase', icon: Sparkles },
  { id: 'media', label: 'Gallery & Videos', icon: ImageIcon },
  { id: 'keys', label: 'Member Key Vault', icon: KeyRound, badge: overview.value.stats?.available_keys },
  { id: 'library', label: 'Ebook Library', icon: BookOpen }
])
</script>

<template>
  <aside
    class="bg-[#09090b] border-b lg:border-b-0 lg:border-r border-zinc-800 flex lg:flex-col items-center justify-between px-3 py-3 lg:py-6 gap-2"
  >
    <!-- Official TriGameDev Brand Seal -->
    <button
      @click="activeSection = 'explore'"
      class="w-12 h-12 bg-black flex items-center justify-center cursor-pointer"
      title="TriGameDev.org Home"
    >
      <img src="/brand/circle-logo-small.png" alt="TriGameDev Seal" class="w-full h-full object-contain" />
    </button>

    <!-- Primary Icon Rail -->
    <nav class="flex lg:flex-col items-center gap-1.5 overflow-x-auto">
      <button
        v-for="item in navItems"
        :key="item.id"
        @click="activeSection = item.id"
        class="relative group w-11 h-11 flex items-center justify-center transition-colors cursor-pointer border"
        :class="
          activeSection === item.id
            ? 'bg-white text-black border-white'
            : 'bg-transparent text-zinc-400 border-transparent hover:text-white hover:border-zinc-700'
        "
        :title="item.label"
      >
        <component :is="item.icon" class="w-5 h-5" />
        <span
          v-if="item.badge"
          class="absolute -top-1.5 -right-1.5 px-1 text-[10px] font-mono font-bold bg-[#dc2626] text-white border border-black"
        >
          {{ item.badge }}
        </span>
      </button>
    </nav>

    <!-- OG Admin CMS Button at Bottom of Rail -->
    <div class="flex items-center">
      <button
        v-if="currentUser?.is_og_admin"
        @click="activeSection = 'cms'"
        class="relative w-11 h-11 flex items-center justify-center transition cursor-pointer border"
        :class="
          activeSection === 'cms'
            ? 'bg-[#dc2626] text-white border-[#dc2626]'
            : 'text-zinc-400 border-zinc-800 hover:text-white hover:border-white'
        "
        title="OG Role Admin CMS"
      >
        <ShieldAlert class="w-5 h-5" />
      </button>
    </div>
  </aside>
</template>
