<script setup>
import { X } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useKeyVault } from '../../composables/useKeyVault'

const { showDonateModal } = useTriGameDevStore()
const { donateForm, submitKeyDonation } = useKeyVault()
</script>

<template>
  <div
    v-if="showDonateModal"
    class="fixed inset-0 z-50 bg-black/85 flex items-center justify-center p-4"
  >
    <div class="bg-[#121215] border border-zinc-700 max-w-lg w-full p-6 flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <h3 class="text-base font-extrabold text-white uppercase tracking-wider">Donate Game Keys to the Pool</h3>
        <button @click="showDonateModal = false" class="text-zinc-400 hover:text-white cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="sm:col-span-2">
          <label class="text-xs font-mono uppercase text-zinc-400">Game / Tool Title</label>
          <input
            v-model="donateForm.title"
            type="text"
            class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
        </div>
        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Platform</label>
          <select
            v-model="donateForm.platform"
            class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          >
            <option>Steam</option>
            <option>Itch.io</option>
            <option>GOG</option>
            <option>Epic</option>
          </select>
        </div>
      </div>
      <div>
        <label class="text-xs font-mono uppercase text-zinc-400">Short Description</label>
        <input
          v-model="donateForm.description"
          type="text"
          class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
        />
      </div>
      <div>
        <label class="text-xs font-mono uppercase text-zinc-400">Paste Key Codes (One key per line)</label>
        <textarea
          v-model="donateForm.raw_keys"
          rows="4"
          class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-xs font-mono text-white focus:outline-none focus:border-white"
        ></textarea>
      </div>
      <div class="flex justify-end gap-2 pt-2">
        <button
          @click="showDonateModal = false"
          class="px-4 py-2 border border-zinc-700 text-xs font-bold uppercase text-zinc-400 hover:text-white cursor-pointer"
        >
          Cancel
        </button>
        <button
          @click="submitKeyDonation"
          class="px-5 py-2.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider cursor-pointer"
        >
          Add Keys to Pool
        </button>
      </div>
    </div>
  </div>
</template>
