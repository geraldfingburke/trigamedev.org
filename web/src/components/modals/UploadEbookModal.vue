<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useEbookLibrary } from '../../composables/useEbookLibrary'

const { showUploadEbookModal } = useTriGameDevStore()
const { ebookForm, submitEbookUpload } = useEbookLibrary()
</script>

<template>
  <div
    v-if="showUploadEbookModal"
    class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
  >
    <div class="bg-[#232234] border border-white/15 rounded-3xl max-w-lg w-full p-6 flex flex-col gap-4 shadow-2xl">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-extrabold text-white">Add Ebook to TriGameDev Library</h3>
        <button @click="showUploadEbookModal = false" class="text-[#8e8ea9] hover:text-white cursor-pointer">✕</button>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label class="text-xs text-[#8e8ea9]">Book Title</label>
          <input
            v-model="ebookForm.title"
            type="text"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          />
        </div>
        <div>
          <label class="text-xs text-[#8e8ea9]">Author</label>
          <input
            v-model="ebookForm.author"
            type="text"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="text-xs text-[#8e8ea9]">Category</label>
          <input
            v-model="ebookForm.category"
            type="text"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          />
        </div>
        <div>
          <label class="text-xs text-[#8e8ea9]">Format</label>
          <select
            v-model="ebookForm.file_format"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          >
            <option>PDF</option>
            <option>EPUB</option>
            <option>ZIP</option>
          </select>
        </div>
      </div>
      <div>
        <label class="text-xs text-[#8e8ea9]">Synopsis / Notes</label>
        <textarea
          v-model="ebookForm.description"
          rows="3"
          class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
        ></textarea>
      </div>
      <div class="flex justify-end gap-2">
        <button @click="showUploadEbookModal = false" class="px-4 py-2 rounded-xl text-xs text-[#8e8ea9] cursor-pointer">
          Cancel
        </button>
        <button
          @click="submitEbookUpload"
          class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-[#ff9f43] to-[#ff4d8d] text-[#14141d] text-xs font-extrabold cursor-pointer"
        >
          Publish to Library
        </button>
      </div>
    </div>
  </div>
</template>
