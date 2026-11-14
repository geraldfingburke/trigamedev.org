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

const navItems = computed(() => {
  const u = currentUser.value
  const unread = u ? (overview.value.unread || {}) : {}
  return [
    { id: 'explore', label: 'Explore', icon: Compass, badge: 0 },
    { id: 'jams', label: 'Jams', icon: Flame, badge: unread.jams || 0 },
    { id: 'events', label: 'Events', icon: Calendar, badge: unread.events || 0 },
    { id: 'news', label: 'News & Announcements', icon: Newspaper, badge: unread.news || 0 },
    { id: 'showcase', label: 'Games by Our People', icon: Sparkles, badge: unread.showcase || 0 },
    { id: 'media', label: 'Photos & Videos', icon: ImageIcon, badge: unread.media || 0 },
    ...(u
      ? [
          { id: 'keys', label: 'Key Vault', icon: KeyRound, badge: unread.keys || 0 },
          { id: 'library', label: 'Ebook Library', icon: BookOpen, badge: unread.library || 0 }
        ]
      : [])
  ]
})
</script>

<template>
  <aside
    class="bg-[#09090b] border-b lg:border-b-0 lg:border-r border-zinc-800 flex lg:flex-col items-center justify-between px-3 py-3 lg:py-6 gap-2"
  >
    <!-- Official TriGameDev Brand Seal (No border) -->
    <button
      @click="activeSection = 'explore'"
      class="w-12 h-12 bg-black flex items-center justify-center cursor-pointer"
      title="TriGameDev.org Home"
    >
      <img src="/brand/circle-logo-small.png" alt="TriGameDev Seal" class="w-full h-full object-contain" />
    </button>

    <!-- Primary Icon Rail with Per-User Unread Notifications -->
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
          v-if="currentUser && item.badge > 0"
          class="absolute -top-1.5 -right-1.5 px-1 min-w-[18px] text-center text-[10px] font-mono font-bold bg-[#dc2626] text-white border border-black"
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
        title="CMS"
      >
        <ShieldAlert class="w-5 h-5" />
      </button>
    </div>
  </aside>
</template>
