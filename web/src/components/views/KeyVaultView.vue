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
  showDonateModal
} = useTriGameDevStore()

const { copiedKeyId, claimKeyFromPool, copyCode } = useKeyVault()
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-extrabold text-white flex items-center gap-2.5">
          <KeyRound class="w-6 h-6 text-[#ff4d8d]" />
          <span>Member Key Sharing Portal</span>
        </h1>
        <p class="text-xs sm:text-sm text-[#9e9eb8] mt-1">
          Single redemption per pool — claiming a key permanently ties it to your account and moves it to
          <strong class="text-white">My Keys</strong>.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <div class="bg-[#232234] p-1 rounded-xl border border-white/10 flex items-center">
          <button
            @click="keySubTab = 'pools'"
            class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer"
            :class="keySubTab === 'pools' ? 'bg-[#6c5ce7] text-white' : 'text-[#8e8ea9] hover:text-white'"
          >
            Available Pools ({{ keyPools.length }})
          </button>
          <button
            @click="keySubTab = 'my-keys'"
            class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer flex items-center gap-1.5"
            :class="keySubTab === 'my-keys' ? 'bg-[#ff4d8d] text-white' : 'text-[#8e8ea9] hover:text-white'"
          >
            <span>My Keys</span>
            <span class="px-1.5 py-0.2 rounded-full bg-black/30 text-[10px]">{{ myKeys.length }}</span>
          </button>
        </div>

        <button
          v-if="currentUser"
          @click="showDonateModal = true"
          class="px-4 py-2 rounded-xl bg-gradient-to-r from-[#ff4d8d] to-[#ff7eb3] text-white text-xs font-bold flex items-center gap-1.5 shadow-lg shadow-[#ff4d8d]/25 cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Donate Keys</span>
        </button>
      </div>
    </div>

    <div
      v-if="!currentUser"
      class="p-6 rounded-3xl bg-[#242336] border border-[#ff4d8d]/40 flex flex-col sm:flex-row items-center justify-between gap-4"
    >
      <div class="flex items-center gap-4">
        <Lock class="w-8 h-8 text-[#ff4d8d] shrink-0" />
        <div>
          <h3 class="font-bold text-white">Discord Member Login Required</h3>
          <p class="text-xs text-[#9e9eb8]">
            Sign in with your TriGameDev Discord account to donate keys, redeem keys, or view your claimed keys.
          </p>
        </div>
      </div>
      <a
        href="/api/go/auth/discord/login"
        class="px-4 py-2 rounded-xl bg-[#5865F2] text-white text-xs font-bold shrink-0"
      >
        Sign In with Discord
      </a>
    </div>

    <div v-if="keySubTab === 'pools'">
      <div
        v-if="filteredKeyPools.length === 0"
        class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
      >
        No key pools have been added yet.
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <div
          v-for="pool in filteredKeyPools"
          :key="pool.id"
          class="bg-[#242336] rounded-3xl p-4 border border-white/[0.07] flex flex-col sm:flex-row gap-4 justify-between"
        >
          <img
            :src="pool.cover_image_url || '/inspo/sunset-studio.jpg'"
            :alt="pool.title"
            class="w-full sm:w-36 h-36 rounded-2xl object-cover shrink-0"
          />
          <div class="flex-1 flex flex-col justify-between gap-3">
            <div>
              <div class="flex items-center justify-between gap-2">
                <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-[#6c5ce7]/25 text-[#a78bfa]">
                  {{ pool.platform }}
                </span>
                <span
                  class="text-xs font-mono font-semibold"
                  :class="pool.available_count > 0 ? 'text-[#00f5d4]' : 'text-[#ff4d8d]'"
                >
                  {{ pool.available_count }} / {{ pool.total_count }} Keys Left
                </span>
              </div>
              <h3 class="font-bold text-base text-white mt-1.5">{{ pool.title }}</h3>
              <p class="text-xs text-[#9e9eb8] line-clamp-2 mt-1">{{ pool.description }}</p>
            </div>

            <div class="flex items-center justify-between pt-2 border-t border-white/[0.06]">
              <span class="text-[11px] text-[#7c7c9a]">Donated by {{ pool.donated_by_name }}</span>

              <button
                v-if="pool.claimed_by_user"
                @click="keySubTab = 'my-keys'"
                class="px-3.5 py-1.5 rounded-xl bg-[#00f5d4]/15 text-[#00f5d4] border border-[#00f5d4]/30 text-xs font-bold cursor-pointer"
              >
                ✓ Claimed (View in My Keys)
              </button>
              <button
                v-else-if="pool.available_count > 0"
                @click="claimKeyFromPool(pool)"
                class="px-4 py-1.5 rounded-xl bg-gradient-to-r from-[#6c5ce7] to-[#8b5cf6] hover:opacity-95 text-white text-xs font-bold shadow-md shadow-[#6c5ce7]/30 cursor-pointer"
              >
                Claim 1 Key
              </button>
              <span v-else class="text-xs font-semibold text-[#7c7c9a]">Pool Depleted</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col gap-4">
      <div
        v-if="myKeys.length === 0"
        class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center flex flex-col items-center gap-3"
      >
        <Gift class="w-10 h-10 text-[#ff4d8d]" />
        <h3 class="font-bold text-base text-white">You haven't claimed any keys yet</h3>
        <p class="text-xs text-[#9e9eb8] max-w-md">
          Keys you redeem from community pools are permanently bound to your account and displayed here.
        </p>
      </div>

      <div
        v-for="item in myKeys"
        :key="item.id"
        class="p-4 rounded-3xl bg-[#242336] border border-[#00f5d4]/30 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
      >
        <div class="flex items-center gap-4">
          <img
            :src="item.cover_image_url || '/inspo/sunset-studio.jpg'"
            class="w-16 h-16 rounded-2xl object-cover"
          />
          <div>
            <div class="flex items-center gap-2">
              <span class="px-2 py-0.5 rounded-md bg-[#6c5ce7]/25 text-[#a78bfa] text-[10px] font-bold">
                {{ item.platform }}
              </span>
              <span class="text-[11px] text-[#7c7c9a]">Claimed {{ item.claimed_at }}</span>
            </div>
            <h4 class="font-bold text-base text-white mt-0.5">{{ item.title }}</h4>
            <p class="text-xs text-[#8e8ea9]">Donated by {{ item.donated_by_name }}</p>
          </div>
        </div>

        <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
          <code
            class="px-3.5 py-2 rounded-xl bg-[#14141d] border border-white/10 text-xs sm:text-sm font-mono text-[#00f5d4] tracking-wider"
          >
            {{ item.key_code }}
          </code>
          <button
            @click="copyCode(item.key_code, item.id)"
            class="p-2.5 rounded-xl bg-[#6c5ce7] hover:bg-[#5a4ad1] text-white transition cursor-pointer"
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
