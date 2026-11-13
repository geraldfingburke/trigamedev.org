<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useKeyVault } from '../../composables/useKeyVault'

const { showDonateModal } = useTriGameDevStore()
const { donateForm, submitKeyDonation } = useKeyVault()
</script>

<template>
  <div
    v-if="showDonateModal"
    class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
  >
    <div class="bg-[#232234] border border-white/15 rounded-3xl max-w-lg w-full p-6 flex flex-col gap-4 shadow-2xl">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-extrabold text-white">Donate Game Keys to the Pool</h3>
        <button @click="showDonateModal = false" class="text-[#8e8ea9] hover:text-white cursor-pointer">✕</button>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="sm:col-span-2">
          <label class="text-xs text-[#8e8ea9]">Game / Tool Title</label>
          <input
            v-model="donateForm.title"
            type="text"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          />
        </div>
        <div>
          <label class="text-xs text-[#8e8ea9]">Platform</label>
          <select
            v-model="donateForm.platform"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          >
            <option>Steam</option>
            <option>Itch.io</option>
            <option>GOG</option>
            <option>Epic</option>
          </select>
        </div>
      </div>
      <div>
        <label class="text-xs text-[#8e8ea9]">Short Description</label>
        <input
          v-model="donateForm.description"
          type="text"
          class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
        />
      </div>
      <div>
        <label class="text-xs text-[#8e8ea9]">Paste Key Codes (One key per line)</label>
        <textarea
          v-model="donateForm.raw_keys"
          rows="4"
          class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-xs font-mono text-[#00f5d4]"
        ></textarea>
      </div>
      <div class="flex justify-end gap-2 pt-2">
        <button @click="showDonateModal = false" class="px-4 py-2 rounded-xl text-xs text-[#8e8ea9] cursor-pointer">
          Cancel
        </button>
        <button
          @click="submitKeyDonation"
          class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-[#ff4d8d] to-[#ff7eb3] text-white text-xs font-bold cursor-pointer"
        >
          Add Keys to Pool
        </button>
      </div>
    </div>
  </div>
</template>
