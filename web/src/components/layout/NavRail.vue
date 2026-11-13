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
  ShieldAlert,
  Gamepad2
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
    class="bg-[#232234] border-b lg:border-b-0 lg:border-r border-white/[0.06] flex lg:flex-col items-center justify-between px-4 py-3 lg:py-6 gap-2"
  >
    <!-- Brand Emblem -->
    <button
      @click="activeSection = 'explore'"
      class="w-11 h-11 rounded-2xl bg-gradient-to-br from-[#6c5ce7] via-[#8b5cf6] to-[#ff4d8d] flex items-center justify-center shadow-lg shadow-[#6c5ce7]/30 hover:scale-105 transition cursor-pointer"
      title="TriGameDev.org Home"
    >
      <Gamepad2 class="w-6 h-6 text-white" />
    </button>

    <!-- Primary Icon Rail -->
    <nav class="flex lg:flex-col items-center gap-2 sm:gap-3 overflow-x-auto">
      <button
        v-for="item in navItems"
        :key="item.id"
        @click="activeSection = item.id"
        class="relative group w-11 h-11 rounded-2xl flex items-center justify-center transition-all cursor-pointer"
        :class="
          activeSection === item.id
            ? 'bg-[#6c5ce7]/20 text-white shadow-inner border border-[#8b5cf6]/40'
            : 'text-[#8e8ea9] hover:text-white hover:bg-white/[0.05]'
        "
        :title="item.label"
      >
        <span
          v-if="activeSection === item.id"
          class="hidden lg:block absolute -left-3 w-1.5 h-6 rounded-r-full bg-gradient-to-b from-[#8b5cf6] to-[#ff4d8d]"
        ></span>
        <component :is="item.icon" class="w-5 h-5" />
        <span
          v-if="item.badge"
          class="absolute -top-1 -right-1 px-1.5 py-0.2 rounded-full text-[10px] font-bold bg-[#ff4d8d] text-white"
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
        class="relative w-11 h-11 rounded-2xl flex items-center justify-center transition cursor-pointer"
        :class="
          activeSection === 'cms'
            ? 'bg-[#ff9f43]/25 text-[#ff9f43] border border-[#ff9f43]/50'
            : 'text-[#8e8ea9] hover:text-[#ff9f43] hover:bg-white/[0.05]'
        "
        title="OG Role Admin CMS"
      >
        <ShieldAlert class="w-5 h-5" />
      </button>
    </div>
  </aside>
</template>
