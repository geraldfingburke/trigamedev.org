<script setup>
import { X } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useEbookLibrary } from '../../composables/useEbookLibrary'

const { showUploadEbookModal } = useTriGameDevStore()
const { ebookForm, submitEbookUpload } = useEbookLibrary()
</script>

<template>
  <div
    v-if="showUploadEbookModal"
    class="fixed inset-0 z-50 bg-black/85 flex items-center justify-center p-4"
  >
    <div class="bg-[#121215] border border-zinc-700 max-w-lg w-full p-6 flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <h3 class="text-base font-extrabold text-white uppercase tracking-wider">Add Ebook to TriGameDev Library</h3>
        <button @click="showUploadEbookModal = false" class="text-zinc-400 hover:text-white cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Book Title</label>
          <input
            v-model="ebookForm.title"
            type="text"
            class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
        </div>
        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Author</label>
          <input
            v-model="ebookForm.author"
            type="text"
            class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Category</label>
          <input
            v-model="ebookForm.category"
            type="text"
            class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
        </div>
        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Format</label>
          <select
            v-model="ebookForm.file_format"
            class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          >
            <option>PDF</option>
            <option>EPUB</option>
            <option>ZIP</option>
          </select>
        </div>
      </div>
      <div>
        <label class="text-xs font-mono uppercase text-zinc-400">Synopsis / Notes</label>
        <textarea
          v-model="ebookForm.description"
          rows="3"
          class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
        ></textarea>
      </div>
      <div class="flex justify-end gap-2 pt-2">
        <button
          @click="showUploadEbookModal = false"
          class="px-4 py-2 border border-zinc-700 text-xs font-bold uppercase text-zinc-400 hover:text-white cursor-pointer"
        >
          Cancel
        </button>
        <button
          @click="submitEbookUpload"
          class="px-5 py-2.5 bg-white hover:bg-zinc-200 text-black text-xs font-extrabold uppercase tracking-wider cursor-pointer"
        >
          Publish to Library
        </button>
      </div>
    </div>
  </div>
</template>
