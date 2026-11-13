<script setup>
import { X } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useKeyVault } from '../../composables/useKeyVault'

const { revealedKeyModal, keySubTab } = useTriGameDevStore()
const { copyCode } = useKeyVault()
</script>

<template>
  <div
    v-if="revealedKeyModal"
    class="fixed inset-0 z-50 bg-black/85 flex items-center justify-center p-4"
  >
    <div class="bg-[#121215] border border-zinc-700 max-w-md w-full p-6 flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[#ef4444] text-[10px] font-mono uppercase tracking-widest">
          KEY CLAIMED & BOUND TO YOUR ACCOUNT
        </span>
        <button @click="revealedKeyModal = null" class="text-zinc-400 hover:text-white cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>
      <h3 class="text-xl font-extrabold text-white uppercase">{{ revealedKeyModal.title }}</h3>
      <p class="text-xs text-zinc-400">
        This key has been removed from the community pool and permanently saved under your
        <strong class="text-white">My Keys</strong> tab.
      </p>
      <div class="p-4 bg-[#09090b] border border-zinc-700 flex items-center justify-between gap-2">
        <code class="text-sm font-mono text-white font-bold">{{ revealedKeyModal.key_code }}</code>
        <button
          @click="copyCode(revealedKeyModal.key_code, 'modal')"
          class="px-3 py-1.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase cursor-pointer"
        >
          Copy
        </button>
      </div>
      <button
        @click="revealedKeyModal = null; keySubTab = 'my-keys'"
        class="w-full py-2.5 bg-[#ef4444] hover:bg-red-500 text-white font-bold text-xs uppercase tracking-wider cursor-pointer"
      >
        Go to My Keys Vault
      </button>
    </div>
  </div>
</template>
