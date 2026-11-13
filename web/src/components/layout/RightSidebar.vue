<script setup>
import { LogOut, ExternalLink } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { activeSection, currentUser, overview, logout } = useTriGameDevStore()
</script>

<template>
  <aside
    class="bg-[#09090b] border-t lg:border-t-0 lg:border-l border-zinc-800 p-6 flex flex-col justify-between gap-8"
  >
    <div class="flex flex-col gap-6">
      <!-- Top User Profile & Discord Status Header -->
      <div class="p-4 bg-[#141417] border border-zinc-800 flex items-center justify-between gap-3">
        <div class="min-w-0">
          <div class="font-mono text-[10px] uppercase tracking-widest text-zinc-400">
            AUTHENTICATION
          </div>
          <div v-if="currentUser" class="truncate mt-0.5">
            <div class="text-xs font-bold text-white truncate uppercase">{{ currentUser.display_name }}</div>
            <div class="font-mono text-[11px] mt-0.5 flex items-center gap-1.5">
              <span
                class="w-2 h-2 inline-block"
                :class="currentUser.is_og_admin ? 'bg-[#dc2626]' : 'bg-white'"
              ></span>
              <span :class="currentUser.is_og_admin ? 'text-[#dc2626] font-bold' : 'text-zinc-300'">
                {{ currentUser.is_og_admin ? 'OG ADMIN ROLE' : 'SERVER MEMBER' }}
              </span>
            </div>
          </div>
          <div v-else class="mt-0.5">
            <div class="text-xs font-bold text-white uppercase">PUBLIC VISITOR</div>
            <div class="font-mono text-[11px] text-zinc-500">UNAUTHENTICATED</div>
          </div>
        </div>

        <div v-if="currentUser" class="flex items-center gap-2 shrink-0">
          <img
            :src="currentUser.avatar_url"
            alt="Avatar"
            class="w-10 h-10 border border-white bg-black object-cover"
          />
          <button
            @click="logout"
            class="p-2.5 bg-black border border-zinc-700 hover:border-white text-zinc-400 hover:text-white cursor-pointer"
            title="Sign Out"
          >
            <LogOut class="w-4 h-4" />
          </button>
        </div>
        <a
          v-else
          href="/api/go/auth/discord/login"
          class="px-3.5 py-2 bg-white hover:bg-zinc-200 text-black font-mono text-xs font-bold uppercase tracking-wider transition shrink-0"
        >
          DISCORD LOGIN
        </a>
      </div>

      <!-- Official Mixtape / QR Card from tgd-brandkit -->
      <div class="flex flex-col gap-3">
        <div class="flex items-center justify-between border-b border-zinc-800 pb-2">
          <h3 class="font-mono text-xs font-bold uppercase tracking-wider text-white">
            // COLLECTIVE DISCORD
          </h3>
          <a
            href="https://discord.gg/6jdhMV6sqt"
            target="_blank"
            class="font-mono text-[11px] text-zinc-400 hover:text-white flex items-center gap-1"
          >
            <span>JOIN SERVER</span>
            <ExternalLink class="w-3 h-3" />
          </a>
        </div>

        <div class="bg-[#141417] border border-zinc-800 p-4 flex flex-col gap-4">
          <a href="https://discord.gg/6jdhMV6sqt" target="_blank" class="block">
            <img
              src="/brand/mixtape-logo.png"
              alt="TriGameDev Mixtape Discord QR"
              class="w-full h-auto object-cover"
            />
          </a>

          <div class="font-mono text-xs text-zinc-300 leading-relaxed border-l-2 border-white pl-3">
            Pros, Hobbyists, Students, Curious Non-Developers, and Game Playing People.
          </div>

          <!-- Metrics Table -->
          <div class="flex flex-col divide-y divide-zinc-800 border-t border-b border-zinc-800 font-mono text-xs">
            <div class="py-2 flex items-center justify-between">
              <span class="text-zinc-400">VAULT KEYS:</span>
              <span class="font-bold text-white">{{ overview.stats?.available_keys || 0 }} AVAILABLE</span>
            </div>
            <div class="py-2 flex items-center justify-between">
              <span class="text-zinc-400">EBOOK ARCHIVE:</span>
              <span class="font-bold text-white">{{ overview.stats?.total_ebooks || 0 }} VOLUMES</span>
            </div>
            <div class="py-2 flex items-center justify-between">
              <span class="text-zinc-400">SHOWCASE BUILDS:</span>
              <span class="font-bold text-white">{{ overview.stats?.total_showcase || 0 }} LISTED</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2">
            <button
              @click="activeSection = 'keys'"
              class="py-2.5 px-3 bg-white hover:bg-zinc-200 text-black font-mono text-xs font-bold uppercase cursor-pointer transition"
            >
              KEY VAULT
            </button>
            <button
              @click="activeSection = 'library'"
              class="py-2.5 px-3 bg-black border border-zinc-700 hover:border-white text-white font-mono text-xs font-bold uppercase cursor-pointer transition"
            >
              EBOOKS
            </button>
          </div>
        </div>
      </div>

      <!-- Origins / First Meetup Snapshot from tgd-brandkit -->
      <div class="bg-[#141417] border border-zinc-800 p-4 flex flex-col gap-3">
        <div class="font-mono text-[11px] uppercase tracking-wider text-zinc-400">
          // ARCHIVE 001: THE WILLOW TREE
        </div>
        <img
          src="/brand/first-meetup.jpg"
          alt="First TriGameDev Meetup"
          class="w-full h-36 object-cover border border-zinc-800 grayscale hover:grayscale-0 transition duration-300"
        />
        <p class="text-xs text-zinc-400 leading-relaxed">
          TriGameDev started in the Tri-Cities region of Northeast Tennessee. While many of our members are still here, plenty of us are not. All are welcome.
        </p>
      </div>
    </div>

    <!-- Bottom Footer -->
    <div class="pt-4 border-t border-zinc-800 flex items-center justify-between font-mono text-[10px] text-zinc-500 uppercase">
      <span>TRIGAMEDEV.ORG</span>
      <span>EST. NORTHEAST TN</span>
    </div>
  </aside>
</template>
