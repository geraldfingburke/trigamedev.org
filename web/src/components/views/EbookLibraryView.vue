<script setup>
import { BookOpen, Plus, Download, Trash2 } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useEbookLibrary } from '../../composables/useEbookLibrary'

const { currentUser, filteredEbooks, showUploadEbookModal } = useTriGameDevStore()
const { downloadEbook, deleteEbook } = useEbookLibrary()
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
      <div>
        <h1 class="text-2xl font-extrabold text-white uppercase flex items-center gap-2.5">
          <BookOpen class="w-6 h-6 text-[#ef4444]" />
          <span>TriGameDev Ebook Library</span>
        </h1>
        <p class="text-xs sm:text-sm text-zinc-400 mt-1">
          Direct-download game design books and reference manuals for TriGameDev members.
        </p>
      </div>

      <button
        v-if="currentUser"
        @click="showUploadEbookModal = true"
        class="px-4 py-2 bg-white hover:bg-zinc-200 text-black font-bold text-xs uppercase tracking-wider flex items-center gap-1.5 cursor-pointer w-fit transition"
      >
        <Plus class="w-4 h-4" />
        <span>Contribute Ebook</span>
      </button>
    </div>

    <div
      v-if="filteredEbooks.length === 0"
      class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
    >
      No ebooks have been uploaded yet.
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="book in filteredEbooks"
        :key="book.id"
        class="bg-[#121215] p-4 border border-zinc-800 flex flex-col sm:flex-row gap-4 justify-between"
      >
        <img
          :src="book.cover_image_url || '/brand/filler.png'"
          :alt="book.title"
          class="w-full sm:w-32 h-44 object-cover bg-black shrink-0 border border-zinc-800"
        />
        <div class="flex-1 flex flex-col justify-between gap-3">
          <div>
            <div class="flex items-center justify-between gap-2">
              <span class="px-2 py-0.5 text-[11px] font-mono uppercase bg-zinc-900 border border-zinc-700 text-[#ef4444]">
                {{ book.category }}
              </span>
              <span class="text-xs font-mono text-zinc-400">
                {{ book.file_format }} {{ book.file_size_label ? `// ${book.file_size_label}` : '' }}
              </span>
            </div>
            <h3 class="font-bold text-base text-white uppercase mt-1.5">{{ book.title }}</h3>
            <p class="text-xs font-mono text-zinc-400">BY {{ book.author }}</p>
            <p class="text-xs text-zinc-400 mt-1.5 line-clamp-2">{{ book.description }}</p>
          </div>

          <div class="flex items-center justify-between pt-2 border-t border-zinc-800 gap-2">
            <div class="flex flex-col">
              <span class="text-[11px] font-mono text-zinc-500">{{ book.download_count }} DOWNLOADS</span>
              <span v-if="book.uploaded_by_name" class="text-[10px] font-mono text-zinc-600 uppercase">
                SHARED BY {{ book.uploaded_by_name }}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <button
                v-if="book.can_delete"
                @click="deleteEbook(book)"
                class="px-2.5 py-2 bg-zinc-900 hover:bg-[#dc2626] text-zinc-300 hover:text-white border border-zinc-700 hover:border-[#dc2626] text-[11px] font-mono uppercase flex items-center gap-1 cursor-pointer transition"
                title="Remove Ebook"
              >
                <Trash2 class="w-3.5 h-3.5" />
                <span>Remove</span>
              </button>
              <button
                @click="downloadEbook(book)"
                class="px-4 py-2 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider flex items-center gap-1.5 cursor-pointer transition"
              >
                <Download class="w-3.5 h-3.5" />
                <span>Direct Download</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
