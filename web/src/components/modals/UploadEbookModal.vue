<script setup>
import { computed } from 'vue'
import { X, Upload, Search, FileText } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useEbookLibrary } from '../../composables/useEbookLibrary'

const { showUploadEbookModal, ebooks } = useTriGameDevStore()

const existingCategories = computed(() => {
  const set = new Set(['Game Design', 'Programming', 'Art & Audio', 'History'])
  for (const b of ebooks.value || []) {
    if (b.category && b.category.trim()) {
      set.add(b.category.trim())
    }
  }
  return Array.from(set).sort((a, b) => a.localeCompare(b))
})
const {
  ebookForm,
  selectedBookFile,
  coverCandidates,
  scrapingBook,
  uploadingBook,
  uploadingCustomCover,
  scrapeBookByQuery,
  handleBookFileSelect,
  handleCustomCoverUpload,
  submitEbookUpload
} = useEbookLibrary()
</script>

<template>
  <div
    v-if="showUploadEbookModal"
    class="fixed inset-0 z-50 bg-black/85 flex items-center justify-center p-4 overflow-y-auto"
  >
    <div class="bg-[#121215] border border-zinc-700 max-w-xl w-full p-6 flex flex-col gap-4 my-auto">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <h3 class="text-base font-extrabold text-white uppercase tracking-wider">Contribute Ebook to Library</h3>
        <button @click="showUploadEbookModal = false" class="text-zinc-400 hover:text-white cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Step 1: Select Book File to Upload -->
      <div class="p-4 bg-[#09090b] border border-zinc-700 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div class="flex items-center gap-3 min-w-0">
          <FileText class="w-5 h-5 text-[#ef4444] shrink-0" />
          <div class="min-w-0">
            <div class="text-xs font-bold text-white uppercase truncate">
              {{ selectedBookFile ? selectedBookFile.name : 'Select Ebook File (PDF, EPUB, ZIP)' }}
            </div>
            <div class="text-[11px] font-mono text-zinc-400">
              {{
                selectedBookFile
                  ? `${ebookForm.file_format} // ${ebookForm.file_size_label} (Auto-scraped title & covers)`
                  : 'Selecting a file automatically scrapes its title & cover options'
              }}
            </div>
          </div>
        </div>
        <label class="px-4 py-2 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider flex items-center gap-1.5 cursor-pointer shrink-0 w-fit">
          <Upload class="w-3.5 h-3.5" />
          <span>Choose Book File</span>
          <input
            type="file"
            accept=".pdf,.epub,.zip,.mobi,.txt"
            @change="handleBookFileSelect"
            class="hidden"
          />
        </label>
      </div>

      <!-- Step 2: Book Title (with Search Covers button) & Author -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="sm:col-span-2">
          <label class="text-xs font-mono uppercase text-zinc-400">Book Title</label>
          <div class="flex gap-1.5 mt-1">
            <input
              v-model="ebookForm.title"
              type="text"
              placeholder="Auto-filled from file..."
              class="flex-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
            <button
              type="button"
              @click="scrapeBookByQuery(ebookForm.title)"
              :disabled="scrapingBook"
              class="px-3 py-2 bg-zinc-900 hover:bg-zinc-800 border border-zinc-700 text-white text-xs font-mono uppercase flex items-center gap-1 cursor-pointer"
              title="Search Featured Cover Images"
            >
              <Search class="w-3.5 h-3.5 text-[#ef4444]" />
              <span>{{ scrapingBook ? '...' : 'Covers' }}</span>
            </button>
          </div>
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

      <!-- Step 3: Featured Image Selection (Scraped Candidates + Custom Upload) -->
      <div class="p-3.5 bg-[#09090b] border border-zinc-800 flex flex-col gap-2.5">
        <div class="flex items-center justify-between">
          <label class="text-xs font-mono uppercase text-zinc-300">
            Select Featured Cover Image
            <span v-if="scrapingBook" class="text-[#ef4444] ml-1">(Searching covers...)</span>
          </label>
          <label class="px-2.5 py-1 bg-zinc-900 hover:bg-zinc-800 border border-zinc-700 text-white text-[10px] font-mono uppercase flex items-center gap-1 cursor-pointer">
            <Upload class="w-3 h-3" />
            <span>{{ uploadingCustomCover ? 'Uploading...' : 'Custom Image' }}</span>
            <input type="file" accept="image/*" @change="handleCustomCoverUpload" class="hidden" />
          </label>
        </div>

        <div v-if="coverCandidates.length > 0" class="grid grid-cols-4 sm:grid-cols-6 gap-2 pt-1">
          <button
            v-for="imgUrl in coverCandidates"
            :key="imgUrl"
            type="button"
            @click="ebookForm.cover_image_url = imgUrl"
            class="relative h-24 bg-black border-2 overflow-hidden cursor-pointer transition"
            :class="ebookForm.cover_image_url === imgUrl ? 'border-[#ef4444]' : 'border-zinc-800 hover:border-white'"
          >
            <img :src="imgUrl" class="w-full h-full object-cover" />
            <span
              v-if="ebookForm.cover_image_url === imgUrl"
              class="absolute bottom-0 inset-x-0 bg-[#ef4444] text-white text-[9px] font-mono uppercase text-center py-0.5"
            >
              Selected
            </span>
          </button>
        </div>
        <div v-else class="text-[11px] font-mono text-zinc-500">
          Select a book file or enter a title and click "Covers" to choose from scraped book covers.
        </div>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="text-xs font-mono uppercase text-zinc-400">Category</label>
          <input
            v-model="ebookForm.category"
            type="text"
            list="ebook-categories-datalist"
            placeholder="e.g. Game Design, Programming..."
            class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
          />
          <datalist id="ebook-categories-datalist">
            <option v-for="cat in existingCategories" :key="cat" :value="cat" />
          </datalist>
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
            <option>MOBI</option>
          </select>
        </div>
      </div>

      <div>
        <label class="text-xs font-mono uppercase text-zinc-400">Synopsis / Notes</label>
        <textarea
          v-model="ebookForm.description"
          rows="2"
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
          :disabled="uploadingBook"
          class="px-5 py-2.5 bg-white hover:bg-zinc-200 disabled:opacity-50 text-black text-xs font-extrabold uppercase tracking-wider cursor-pointer"
        >
          {{ uploadingBook ? 'Uploading Book...' : 'Upload & Publish' }}
        </button>
      </div>
    </div>
  </div>
</template>
