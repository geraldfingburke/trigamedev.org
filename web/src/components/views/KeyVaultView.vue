<script setup>
import { KeyRound, Plus, Lock, Gift, Check, Copy } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useKeyVault } from '../../composables/useKeyVault'

const {
  keySubTab,
  currentUser,
  keyPools,
  myKeys,
  filteredKeyPools,
  showDonateModal,
  loginWithDiscord
} = useTriGameDevStore()

const { copiedKeyId, claimKeyFromPool, copyCode } = useKeyVault()
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
      <div>
        <h1 class="text-2xl font-extrabold text-white uppercase flex items-center gap-2.5">
          <KeyRound class="w-6 h-6 text-[#ef4444]" />
          <span>Key Vault</span>
        </h1>
        <p class="text-xs sm:text-sm text-zinc-400 mt-1">
          Claim a game key or donate keys to the community.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <div class="bg-[#121215] p-1 border border-zinc-800 flex items-center">
          <button
            @click="keySubTab = 'pools'"
            class="px-3.5 py-1.5 text-xs font-bold uppercase tracking-wider transition cursor-pointer"
            :class="keySubTab === 'pools' ? 'bg-white text-black' : 'text-zinc-400 hover:text-white'"
          >
            Projects ({{ keyPools.length }})
          </button>
          <button
            @click="keySubTab = 'my-keys'"
            class="px-3.5 py-1.5 text-xs font-bold uppercase tracking-wider transition cursor-pointer flex items-center gap-1.5"
            :class="keySubTab === 'my-keys' ? 'bg-[#ef4444] text-white' : 'text-zinc-400 hover:text-white'"
          >
            <span>My Keys</span>
            <span class="px-1.5 py-0.2 bg-black/40 text-[10px] font-mono">{{ myKeys.length }}</span>
          </button>
        </div>

        <button
          v-if="currentUser"
          @click="showDonateModal = true"
          class="px-4 py-2 bg-[#ef4444] hover:bg-red-500 text-white text-xs font-bold uppercase tracking-wider flex items-center gap-1.5 cursor-pointer transition"
        >
          <Plus class="w-4 h-4" />
          <span>Donate Keys</span>
        </button>
      </div>
    </div>

    <div
      v-if="!currentUser"
      class="p-6 bg-[#121215] border border-zinc-700 flex flex-col sm:flex-row items-center justify-between gap-4"
    >
      <div class="flex items-center gap-4">
        <Lock class="w-7 h-7 text-[#ef4444] shrink-0" />
        <div>
          <h3 class="font-bold text-white uppercase">Discord Member Login Required</h3>
          <p class="text-xs text-zinc-400">
            Sign in with Discord to donate or claim keys.
          </p>
        </div>
      </div>
      <button
        @click="loginWithDiscord"
        class="px-4 py-2 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider shrink-0 cursor-pointer"
      >
        Sign In with Discord
      </button>
    </div>

    <div v-if="keySubTab === 'pools'">
      <div
        v-if="filteredKeyPools.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No projects have been added yet.
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <div
          v-for="project in filteredKeyPools"
          :key="project.id"
          class="bg-[#121215] p-4 border border-zinc-800 flex flex-col sm:flex-row gap-4 justify-between"
        >
          <img
            :src="project.cover_image_url || '/brand/filler.png'"
            :alt="project.title"
            class="w-full sm:w-36 h-36 object-cover bg-black border border-zinc-800 shrink-0"
          />
          <div class="flex-1 flex flex-col justify-between gap-3">
            <div>
              <div class="flex items-center justify-between gap-2">
                <span class="px-2 py-0.5 text-[11px] font-mono uppercase bg-zinc-900 border border-zinc-700 text-zinc-300">
                  {{ project.platform }}
                </span>
                <span
                  class="text-xs font-mono font-semibold"
                  :class="project.available_count > 0 ? 'text-white' : 'text-[#ef4444]'"
                >
                  {{ project.available_count }} / {{ project.total_count }} KEYS LEFT
                </span>
              </div>
              <h3 class="font-bold text-base text-white uppercase mt-1.5">{{ project.title }}</h3>
              <p class="text-xs text-zinc-400 line-clamp-2 mt-1">{{ project.description }}</p>
            </div>

            <div class="flex items-center justify-between pt-2 border-t border-zinc-800">
              <span class="text-[11px] font-mono text-zinc-500">BY {{ project.donated_by_name }}</span>

              <button
                v-if="project.claimed_by_user"
                @click="keySubTab = 'my-keys'"
                class="px-3.5 py-1.5 bg-zinc-900 text-white border border-zinc-600 text-xs font-bold uppercase cursor-pointer"
              >
                CLAIMED (VIEW MY KEYS)
              </button>
              <button
                v-else-if="project.available_count > 0"
                @click="claimKeyFromPool(project)"
                class="px-4 py-1.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider cursor-pointer transition"
              >
                Claim Key
              </button>
              <span v-else class="text-xs font-mono uppercase text-zinc-500">Depleted</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col gap-4">
      <div
        v-if="myKeys.length === 0"
        class="p-10 bg-[#121215] border border-zinc-800 text-center flex flex-col items-center gap-3"
      >
        <Gift class="w-9 h-9 text-[#ef4444]" />
        <h3 class="font-bold text-base text-white uppercase">You haven't claimed any keys yet</h3>
        <p class="text-xs text-zinc-400 max-w-md">
          Keys you redeem from projects are permanently bound to your account and displayed here.
        </p>
      </div>

      <div
        v-for="item in myKeys"
        :key="item.id"
        class="p-4 bg-[#121215] border border-zinc-700 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
      >
        <div class="flex items-center gap-4">
          <img
            :src="item.cover_image_url || '/brand/filler.png'"
            class="w-16 h-16 object-cover bg-black border border-zinc-800"
          />
          <div>
            <div class="flex items-center gap-2">
              <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-zinc-300 text-[10px] font-mono uppercase">
                {{ item.platform }}
              </span>
              <span class="text-[11px] font-mono text-zinc-500">CLAIMED {{ item.claimed_at }}</span>
            </div>
            <h4 class="font-bold text-base text-white uppercase mt-0.5">{{ item.title }}</h4>
            <p class="text-xs text-zinc-400">Donated by {{ item.donated_by_name }}</p>
          </div>
        </div>

        <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
          <code
            class="px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-xs sm:text-sm font-mono text-white tracking-wider"
          >
            {{ item.key_code }}
          </code>
          <button
            @click="copyCode(item.key_code, item.id)"
            class="p-2.5 bg-white hover:bg-zinc-200 text-black transition cursor-pointer"
            title="Copy Key"
          >
            <Check v-if="copiedKeyId === item.id" class="w-4 h-4" />
            <Copy v-else class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
