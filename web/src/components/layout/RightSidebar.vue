<script setup>
import { LogOut, Radio } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { activeSection, currentUser, overview, nextEvent, logout } = useTriGameDevStore()
</script>

<template>
  <aside
    class="bg-gradient-to-b from-[#292643] via-[#211f35] to-[#1d1c2d] border-t lg:border-t-0 lg:border-l border-white/[0.07] p-6 flex flex-col justify-between gap-6"
  >
    <div class="flex flex-col gap-6">
      <!-- Top User Profile & Discord Status Header -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3 min-w-0">
          <div v-if="currentUser" class="min-w-0">
            <div class="text-xs font-bold text-white truncate">{{ currentUser.display_name }}</div>
            <div class="text-[11px] font-semibold flex items-center gap-1">
              <span
                class="w-2 h-2 rounded-full"
                :class="currentUser.is_og_admin ? 'bg-[#ff9f43]' : 'bg-[#00f5d4]'"
              ></span>
              <span :class="currentUser.is_og_admin ? 'text-[#ff9f43]' : 'text-[#00f5d4]'">
                {{ currentUser.is_og_admin ? 'OG Role Admin' : 'TriGameDev Member' }}
              </span>
            </div>
          </div>
          <div v-else>
            <div class="text-xs font-bold text-white">Guest Visitor</div>
            <div class="text-[11px] text-[#8e8ea9]">Sign in via Discord</div>
          </div>
        </div>

        <div v-if="currentUser" class="flex items-center gap-2">
          <img
            :src="currentUser.avatar_url"
            alt="Avatar"
            class="w-10 h-10 rounded-2xl border-2 border-[#8b5cf6] bg-[#14141d] object-cover"
          />
          <button
            @click="logout"
            class="p-2 rounded-xl bg-white/[0.06] hover:bg-white/10 text-[#9e9eb8] hover:text-white cursor-pointer"
            title="Sign Out"
          >
            <LogOut class="w-4 h-4" />
          </button>
        </div>
        <a
          v-else
          href="/api/go/auth/discord/login"
          class="px-3.5 py-2 rounded-xl bg-[#5865F2] hover:bg-[#4752c4] text-white text-xs font-bold transition"
        >
          Discord Login
        </a>
      </div>

      <!-- Featured Event / Jam Spotlight Card -->
      <div class="flex flex-col gap-3">
        <div class="flex items-center justify-between">
          <h3 class="font-bold text-sm text-white">Featured Spotlight</h3>
          <button
            @click="activeSection = 'events'"
            class="text-xs text-[#c4c4dc] hover:text-white underline cursor-pointer"
          >
            Events
          </button>
        </div>

        <div class="bg-[#1b1a29] rounded-3xl p-4 border border-white/[0.08] flex flex-col gap-4 shadow-xl">
          <div class="relative h-44 rounded-2xl overflow-hidden">
            <img
              :src="nextEvent?.cover_image_url || '/inspo/rooftop-twilight.jpg'"
              alt="Featured Event"
              class="w-full h-full object-cover"
            />
            <span
              class="absolute bottom-2.5 right-2.5 px-2.5 py-1 rounded-lg bg-white text-[#14141d] text-[10px] font-extrabold shadow"
            >
              TriGameDev ✓
            </span>
          </div>

          <div
            class="p-3 rounded-2xl bg-[#27263b] border border-white/10 flex items-center justify-between gap-2"
          >
            <div class="min-w-0">
              <div class="font-bold text-xs text-white truncate">
                {{ nextEvent?.title || 'TriGameDev Community Hub' }}
              </div>
              <div class="text-[10px] text-[#9e9eb8] truncate">
                {{ nextEvent?.location || 'Join us on Discord' }}
              </div>
            </div>
            <span
              class="w-8 h-8 rounded-full bg-[#ff2e7e] flex items-center justify-center text-white shrink-0 shadow-md shadow-[#ff2e7e]/40"
            >
              <Radio class="w-4 h-4" />
            </span>
          </div>

          <div class="flex flex-col gap-2 text-xs px-1">
            <div class="flex items-center justify-between">
              <span class="text-[#8e8ea9]">Available Vault Keys:</span>
              <span class="font-bold text-[#00f5d4]">{{ overview.stats?.available_keys || 0 }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-[#8e8ea9]">Library Ebooks:</span>
              <span class="font-bold text-[#00f5d4]">{{ overview.stats?.total_ebooks || 0 }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-[#8e8ea9]">Showcase Games:</span>
              <span class="font-bold text-[#00f5d4]">{{ overview.stats?.total_showcase || 0 }}</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2.5 pt-1">
            <button
              @click="activeSection = 'keys'"
              class="py-2 rounded-xl bg-gradient-to-r from-[#6c5ce7] to-[#8b5cf6] text-white text-xs font-bold cursor-pointer shadow-md shadow-[#6c5ce7]/25"
            >
              Key Vault
            </button>
            <button
              @click="activeSection = 'library'"
              class="py-2 rounded-xl border border-white/20 hover:border-white/40 text-white text-xs font-bold cursor-pointer"
            >
              Ebook Library
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Pulse Decorative Graphic -->
    <div class="flex flex-col gap-2 pt-2 border-t border-white/[0.06]">
      <div class="flex items-center justify-between">
        <span class="text-xs font-bold text-white">TriGameDev Signal</span>
        <span class="text-xs font-bold text-[#00f5d4]">ONLINE</span>
      </div>
      <div class="relative pt-2 pb-1">
        <svg viewBox="0 0 260 70" class="w-full h-14 overflow-visible">
          <defs>
            <linearGradient id="pulseGrad" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#8b5cf6" stop-opacity="0.45" />
              <stop offset="100%" stop-color="#8b5cf6" stop-opacity="0.0" />
            </linearGradient>
          </defs>
          <path
            d="M0,55 C35,50 55,20 90,38 C125,55 145,15 180,22 C215,28 235,8 260,12 L260,70 L0,70 Z"
            fill="url(#pulseGrad)"
          />
          <path
            d="M0,55 C35,50 55,20 90,38 C125,55 145,15 180,22 C215,28 235,8 260,12"
            fill="none"
            stroke="#8b5cf6"
            stroke-width="2.5"
          />
          <circle cx="90" cy="38" r="4" fill="#00f5d4" />
          <circle cx="180" cy="22" r="4.5" fill="#00f5d4" stroke="#ffffff" stroke-width="1.5" />
          <circle cx="235" cy="11" r="4" fill="#00f5d4" />
        </svg>
      </div>
    </div>
  </aside>
</template>
