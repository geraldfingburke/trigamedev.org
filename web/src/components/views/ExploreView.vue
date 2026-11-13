<script setup>
import {
  Radio,
  KeyRound,
  ChevronRight,
  Eye,
  Sparkles,
  Newspaper,
  BookOpen
} from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const {
  activeSection,
  overview,
  currentUser,
  activeJam,
  filteredShowcase,
  selectedLightbox
} = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-8">
    <!-- Full-Width Lofi Sunset Hero Banner -->
    <section
      class="relative rounded-3xl overflow-hidden border border-white/[0.1] min-h-[230px] flex items-center p-6 sm:p-10 shadow-2xl"
    >
      <img
        src="/inspo/sunset-studio.jpg"
        alt="TriGameDev Lofi Sunset Banner"
        class="absolute inset-0 w-full h-full object-cover object-center scale-105 filter brightness-90"
      />
      <div
        class="absolute inset-0 bg-gradient-to-r from-[#1b1838]/95 via-[#433089]/80 to-[#ff4d8d]/35"
      ></div>

      <div class="relative z-10 max-w-xl flex flex-col gap-4">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-black/40 backdrop-blur-md border border-white/15 w-fit text-xs font-semibold text-[#ff7eb3]">
          <Radio class="w-3.5 h-3.5 text-[#ff4d8d]" />
          <span>{{ activeJam ? `${activeJam.month_label} • LIVE JAM` : 'TRIGAMEDEV COLLECTIVE' }}</span>
        </div>

        <h1 class="text-2xl sm:text-4xl font-extrabold tracking-tight leading-tight text-white">
          {{ activeJam ? activeJam.title : 'TriGameDev Community Hub' }}
        </h1>

        <p class="text-sm sm:text-base text-white/85 leading-relaxed">
          <template v-if="activeJam">
            Theme: <span class="font-mono font-semibold text-[#ffd166]">{{ activeJam.theme }}</span> —
            {{ activeJam.description }}
          </template>
          <template v-else>
            Connect with local & indie game developers on Discord. Share builds, join monthly jams, swap keys, and explore community resources.
          </template>
        </p>

        <div class="flex flex-wrap items-center gap-3 pt-1">
          <button
            @click="activeSection = 'jams'"
            class="px-6 py-2.5 rounded-xl bg-white text-[#18162b] font-bold text-sm shadow-lg hover:bg-[#f0ecff] transition cursor-pointer"
          >
            Monthly Jams
          </button>
          <button
            @click="activeSection = 'keys'"
            class="px-5 py-2.5 rounded-xl bg-black/45 hover:bg-black/65 text-white border border-white/20 backdrop-blur-md font-semibold text-sm flex items-center gap-2 transition cursor-pointer"
          >
            <KeyRound class="w-4 h-4 text-[#00f5d4]" />
            <span>Key Vault ({{ overview.stats?.available_keys || 0 }} Available)</span>
          </button>
        </div>
      </div>
    </section>

    <!-- 3-Column Featured Showcase Cards -->
    <section class="flex flex-col gap-4">
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-bold text-white flex items-center gap-2">
          <span>Member Showcase</span>
          <span class="text-xs px-2.5 py-0.5 rounded-full bg-[#6c5ce7]/25 text-[#a78bfa] font-semibold">
            Community Builds
          </span>
        </h2>
        <button
          @click="activeSection = 'showcase'"
          class="px-4 py-2 rounded-xl bg-[#6c5ce7] hover:bg-[#5b4bc4] text-xs font-semibold text-white flex items-center gap-1.5 transition cursor-pointer shadow-md shadow-[#6c5ce7]/25"
        >
          <span>View All</span>
          <ChevronRight class="w-4 h-4" />
        </button>
      </div>

      <div
        v-if="filteredShowcase.length === 0"
        class="p-8 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
      >
        No showcase games published yet. OG Admins can add projects from the CMS.
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <article
          v-for="game in filteredShowcase.slice(0, 3)"
          :key="game.id"
          class="group bg-[#242336] rounded-3xl p-3.5 border border-white/[0.06] hover:border-[#8b5cf6]/50 transition-all flex flex-col justify-between shadow-xl"
        >
          <div class="relative h-56 rounded-2xl overflow-hidden bg-[#181726]">
            <img
              :src="game.thumbnail_url || '/inspo/sunset-studio.jpg'"
              :alt="game.title"
              class="w-full h-full object-cover group-hover:scale-105 transition duration-500"
            />
            <div
              class="absolute top-3 right-3 px-2.5 py-1 rounded-full bg-white/90 text-[#14141d] text-[11px] font-bold flex items-center gap-1 shadow"
            >
              <Eye class="w-3.5 h-3.5" />
              <span>{{ game.views_count }}</span>
            </div>
            <div
              class="absolute bottom-3 left-3 right-3 p-3 rounded-2xl bg-[#1c1b29]/80 backdrop-blur-md border border-white/15 flex items-center justify-between gap-2"
            >
              <div class="min-w-0">
                <h3 class="font-bold text-sm text-white truncate">{{ game.title }}</h3>
                <p class="text-[11px] text-[#b4b4cc] truncate">by {{ game.developer_name }}</p>
              </div>
              <a
                :href="game.play_url || '#'"
                target="_blank"
                class="w-9 h-9 rounded-full bg-gradient-to-tr from-[#ff2e7e] to-[#ff6b9d] flex items-center justify-center text-white shrink-0 shadow-lg shadow-[#ff2e7e]/40 hover:scale-105 transition"
              >
                <Radio class="w-4 h-4" />
              </a>
            </div>
          </div>

          <div class="pt-4 px-1 flex flex-col gap-3">
            <div class="flex items-center justify-between text-xs">
              <span class="text-[#8e8ea9]">Engine:</span>
              <span class="font-semibold text-[#00f5d4] flex items-center gap-1">
                <Sparkles class="w-3.5 h-3.5 text-[#8b5cf6]" />
                {{ game.engine }}
              </span>
            </div>
            <div class="grid grid-cols-2 gap-2.5">
              <a
                :href="game.play_url || '#'"
                target="_blank"
                class="py-2 px-3 rounded-xl bg-gradient-to-r from-[#6c5ce7] to-[#8b5cf6] hover:opacity-95 text-white font-semibold text-xs text-center shadow-md shadow-[#6c5ce7]/25 transition"
              >
                Play Build
              </a>
              <a
                :href="game.source_url || '#'"
                target="_blank"
                class="py-2 px-3 rounded-xl bg-transparent border border-white/20 hover:border-white/40 text-white/90 font-semibold text-xs text-center transition"
              >
                Source / Info
              </a>
            </div>
          </div>
        </article>
      </div>
    </section>

    <!-- 7-Column Compact Gallery Strip -->
    <section v-if="overview.gallery && overview.gallery.length > 0" class="flex flex-col gap-3">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-bold text-white">Recent Gallery Snaps</h2>
        <button
          @click="activeSection = 'media'"
          class="text-xs font-semibold text-[#a78bfa] hover:underline cursor-pointer"
        >
          Open Gallery & Videos →
        </button>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-7 gap-3">
        <button
          v-for="snap in overview.gallery.slice(0, 7)"
          :key="snap.id"
          @click="selectedLightbox = snap"
          class="group relative aspect-square rounded-2xl overflow-hidden bg-[#242336] border border-white/[0.07] hover:border-[#ff4d8d] transition cursor-pointer"
        >
          <img
            :src="snap.image_url"
            :alt="snap.title"
            class="w-full h-full object-cover group-hover:scale-110 transition duration-300"
          />
        </button>
      </div>
    </section>

    <!-- Two-Column Quick Feed: Latest News & Member Portals -->
    <section class="grid grid-cols-1 xl:grid-cols-2 gap-5">
      <!-- Latest News & Discord Announcements -->
      <div class="bg-[#242336] rounded-3xl p-5 border border-white/[0.06] flex flex-col gap-4">
        <div class="flex items-center justify-between">
          <h3 class="font-bold text-base text-white flex items-center gap-2">
            <Newspaper class="w-4 h-4 text-[#ff9f43]" />
            <span>News & Announcements</span>
          </h3>
          <button
            @click="activeSection = 'news'"
            class="text-xs text-[#a78bfa] hover:underline cursor-pointer"
          >
            View All
          </button>
        </div>

        <div v-if="!overview.news || overview.news.length === 0" class="text-xs text-[#8e8ea9] py-6 text-center">
          No announcements posted yet.
        </div>

        <div v-else class="flex flex-col gap-3">
          <div
            v-for="post in overview.news.slice(0, 3)"
            :key="post.id"
            class="p-3.5 rounded-2xl bg-[#1b1a29] border border-white/[0.04] flex flex-col gap-1"
          >
            <div class="flex items-center justify-between text-[11px]">
              <span class="px-2 py-0.5 rounded-md bg-[#6c5ce7]/20 text-[#a78bfa] font-semibold">
                {{ post.tag }}
              </span>
              <span class="text-[#7c7c9a]">{{ post.published_at }} • {{ post.author_name }}</span>
            </div>
            <h4 class="font-bold text-sm text-white mt-1">{{ post.title }}</h4>
            <p class="text-xs text-[#9e9eb8] line-clamp-2">{{ post.summary }}</p>
          </div>
        </div>
      </div>

      <!-- Member Vault Quick Access -->
      <div class="bg-[#242336] rounded-3xl p-5 border border-white/[0.06] flex flex-col justify-between gap-4">
        <div class="flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span class="px-2.5 py-1 rounded-full text-[11px] font-bold bg-[#00f5d4]/15 text-[#00f5d4]">
              DISCORD MEMBER PERKS
            </span>
          </div>
          <h3 class="text-lg font-bold text-white mt-1">
            Key Sharing Portal & Ebook Library
          </h3>
          <p class="text-xs text-[#9e9eb8] leading-relaxed">
            Sign in with your TriGameDev Discord account to donate or claim game keys into your permanent
            <strong class="text-white">My Keys</strong> vault and download community ebooks.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <button
            @click="activeSection = 'keys'"
            class="p-4 rounded-2xl bg-gradient-to-br from-[#2d294e] to-[#1f1d36] border border-[#8b5cf6]/30 hover:border-[#8b5cf6] text-left transition cursor-pointer group"
          >
            <KeyRound class="w-5 h-5 text-[#ff4d8d] mb-2 group-hover:scale-110 transition" />
            <div class="font-bold text-sm text-white">Key Sharing Vault</div>
            <div class="text-xs text-[#00f5d4] mt-0.5">
              {{ overview.stats?.available_keys || 0 }} keys available →
            </div>
          </button>

          <button
            @click="activeSection = 'library'"
            class="p-4 rounded-2xl bg-gradient-to-br from-[#2b2540] to-[#1d1b2f] border border-[#ff4d8d]/30 hover:border-[#ff4d8d] text-left transition cursor-pointer group"
          >
            <BookOpen class="w-5 h-5 text-[#ff9f43] mb-2 group-hover:scale-110 transition" />
            <div class="font-bold text-sm text-white">Ebook Library</div>
            <div class="text-xs text-[#ff9f43] mt-0.5">
              {{ overview.stats?.total_ebooks || 0 }} direct-download books →
            </div>
          </button>
        </div>
      </div>
    </section>

    <!-- Bottom Discord Connect Bar -->
    <section
      v-if="!currentUser"
      class="p-4 sm:px-6 rounded-2xl bg-[#242336] border border-white/[0.07] flex flex-col sm:flex-row items-center justify-between gap-4"
    >
      <div class="flex items-center gap-3">
        <div class="w-3 h-3 rounded-full bg-[#ff4d8d]"></div>
        <p class="text-xs sm:text-sm text-[#c4c4dc]">
          Member of the <strong class="text-white">TriGameDev Discord</strong>? Sign in to unlock the Key Vault, Ebook Library, and OG CMS.
        </p>
      </div>
      <a
        href="/api/go/auth/discord/login"
        class="px-5 py-2.5 rounded-xl bg-[#5865F2] hover:bg-[#4752c4] text-white text-xs font-bold whitespace-nowrap shadow-md transition"
      >
        Sign In with Discord
      </a>
    </section>
  </div>
</template>
