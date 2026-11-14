<script setup>
import { LogOut, ExternalLink } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { currentUser, staySignedIn, overview, loginWithDiscord, logout } = useTriGameDevStore()
</script>

<template>
  <aside
    class="bg-[#09090b] border-t lg:border-t-0 lg:border-l border-zinc-800 p-6 flex flex-col justify-between gap-8"
  >
    <div class="flex flex-col gap-6">
      <!-- Top User Profile & Discord Status Header -->
      <div class="p-4 bg-[#141417] border border-zinc-800 flex flex-col gap-3">
        <div class="flex items-center justify-between gap-3">
          <div class="min-w-0">
            <div class="font-mono text-[10px] uppercase tracking-widest text-zinc-400">
              AUTHENTICATION
            </div>
            <div v-if="currentUser" class="mt-1">
              <div class="text-xs font-bold text-white truncate uppercase">{{ currentUser.display_name }}</div>
              <!-- User Roles List (can have multiple roles) -->
              <div class="flex flex-wrap gap-1 mt-1.5">
                <span
                  v-for="role in (currentUser.roles && currentUser.roles.length ? currentUser.roles : ['Member'])"
                  :key="role"
                  class="px-1.5 py-0.5 text-[10px] font-mono uppercase border"
                  :class="
                    role === 'OG'
                      ? 'bg-[#dc2626]/20 border-[#dc2626] text-[#ef4444] font-bold'
                      : 'bg-zinc-900 border-zinc-700 text-zinc-300'
                  "
                >
                  {{ role }}
                </span>
              </div>
            </div>
            <div v-else class="mt-0.5">
              <div class="text-xs font-bold text-white uppercase">GUEST</div>
              <div class="font-mono text-[11px] text-zinc-500">NOT SIGNED IN</div>
            </div>
          </div>

          <div v-if="currentUser" class="flex items-center gap-2 shrink-0">
            <img
              :src="currentUser.avatar_url"
              alt="Avatar"
              class="w-10 h-10 bg-black object-cover"
            />
            <button
              @click="logout"
              class="p-2.5 bg-black border border-zinc-700 hover:border-white text-zinc-400 hover:text-white cursor-pointer"
              title="Sign Out"
            >
              <LogOut class="w-4 h-4" />
            </button>
          </div>
          <button
            v-else
            @click="loginWithDiscord"
            class="px-3.5 py-2 bg-white hover:bg-zinc-200 text-black font-mono text-xs font-bold uppercase tracking-wider transition shrink-0 cursor-pointer"
          >
            DISCORD LOGIN
          </button>
        </div>

        <!-- Stay Signed In Option for Guests -->
        <label
          v-if="!currentUser"
          class="flex items-center gap-2 pt-2 border-t border-zinc-800 text-[11px] font-mono text-zinc-400 cursor-pointer select-none"
        >
          <input
            v-model="staySignedIn"
            type="checkbox"
            class="accent-[#ef4444] w-3.5 h-3.5 cursor-pointer"
          />
          <span>Stay signed in on this device</span>
        </label>
      </div>

      <!-- Official Mixtape / QR Card & Community Stats -->
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

          <!-- Community Metrics: Total Members (Discord), Number of Jams, Released Games -->
          <div class="flex flex-col divide-y divide-zinc-800 border-t border-b border-zinc-800 font-mono text-xs">
            <div class="py-2 flex items-center justify-between">
              <span class="text-zinc-400">TOTAL MEMBERS:</span>
              <span class="font-bold text-white">{{ overview.stats?.total_members || 162 }}</span>
            </div>
            <div class="py-2 flex items-center justify-between">
              <span class="text-zinc-400">NUMBER OF JAMS:</span>
              <span class="font-bold text-white">{{ overview.stats?.total_jams || 0 }}</span>
            </div>
            <div class="py-2 flex items-center justify-between">
              <span class="text-zinc-400">RELEASED GAMES:</span>
              <span class="font-bold text-white">{{ overview.stats?.released_games || 0 }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Origins / First Meetup Snapshot -->
      <div class="bg-[#141417] border border-zinc-800 p-4 flex flex-col gap-3">
        <div class="font-mono text-[11px] uppercase tracking-wider text-zinc-400">
          // OUR FIRST MEETUP: THE WILLOW TREE
        </div>
        <img
          src="/brand/first-meetup.jpg"
          alt="First TriGameDev Meetup at The Willow Tree"
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
