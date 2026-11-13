<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useKeyVault } from '../../composables/useKeyVault'

const { revealedKeyModal, keySubTab } = useTriGameDevStore()
const { copyCode } = useKeyVault()
</script>

<template>
  <div
    v-if="revealedKeyModal"
    class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
  >
    <div class="bg-[#232234] border border-[#00f5d4]/40 rounded-3xl max-w-md w-full p-6 flex flex-col gap-4 shadow-2xl">
      <div class="flex items-center justify-between">
        <span class="px-2.5 py-1 rounded-full bg-[#00f5d4]/15 text-[#00f5d4] text-xs font-bold">
          KEY CLAIMED & BOUND TO YOUR ACCOUNT
        </span>
        <button @click="revealedKeyModal = null" class="text-[#8e8ea9] hover:text-white text-sm cursor-pointer">✕</button>
      </div>
      <h3 class="text-xl font-extrabold text-white">{{ revealedKeyModal.title }}</h3>
      <p class="text-xs text-[#9e9eb8]">
        This key has been removed from the community pool and permanently saved under your
        <strong class="text-white">My Keys</strong> tab.
      </p>
      <div class="p-4 rounded-2xl bg-[#14141d] border border-white/10 flex items-center justify-between gap-2">
        <code class="text-sm font-mono text-[#00f5d4] font-bold">{{ revealedKeyModal.key_code }}</code>
        <button
          @click="copyCode(revealedKeyModal.key_code, 'modal')"
          class="px-3 py-1.5 rounded-xl bg-[#6c5ce7] text-white text-xs font-bold cursor-pointer"
        >
          Copy
        </button>
      </div>
      <button
        @click="revealedKeyModal = null; keySubTab = 'my-keys'"
        class="w-full py-2.5 rounded-xl bg-white text-[#14141d] font-bold text-xs cursor-pointer"
      >
        Go to My Keys Vault
      </button>
    </div>
  </div>
</template>
