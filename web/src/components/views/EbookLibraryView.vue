<script setup>
import { BookOpen, Plus, Download } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useEbookLibrary } from '../../composables/useEbookLibrary'

const { currentUser, filteredEbooks, showUploadEbookModal } = useTriGameDevStore()
const { downloadEbook } = useEbookLibrary()
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-extrabold text-white flex items-center gap-2.5">
          <BookOpen class="w-6 h-6 text-[#ff9f43]" />
          <span>TriGameDev Ebook Library</span>
        </h1>
        <p class="text-xs sm:text-sm text-[#9e9eb8] mt-1">
          Direct-download game design books and reference manuals for TriGameDev members.
        </p>
      </div>

      <button
        v-if="currentUser"
        @click="showUploadEbookModal = true"
        class="px-4 py-2 rounded-xl bg-gradient-to-r from-[#ff9f43] to-[#ff4d8d] text-[#14141d] font-bold text-xs flex items-center gap-1.5 shadow-lg cursor-pointer w-fit"
      >
        <Plus class="w-4 h-4" />
        <span>Contribute Ebook</span>
      </button>
    </div>

    <div
      v-if="filteredEbooks.length === 0"
      class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
    >
      No ebooks have been uploaded yet.
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="book in filteredEbooks"
        :key="book.id"
        class="bg-[#242336] rounded-3xl p-4 border border-white/[0.07] flex flex-col sm:flex-row gap-4 justify-between"
      >
        <img
          :src="book.cover_image_url || '/inspo/rooftop-twilight.jpg'"
          :alt="book.title"
          class="w-full sm:w-32 h-44 rounded-2xl object-cover shrink-0 border border-white/10"
        />
        <div class="flex-1 flex flex-col justify-between gap-3">
          <div>
            <div class="flex items-center justify-between gap-2">
              <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-[#ff9f43]/20 text-[#ff9f43]">
                {{ book.category }}
              </span>
              <span class="text-xs font-mono text-[#8e8ea9]">
                {{ book.file_format }} {{ book.file_size_label ? `• ${book.file_size_label}` : '' }}
              </span>
            </div>
            <h3 class="font-bold text-base text-white mt-1.5">{{ book.title }}</h3>
            <p class="text-xs text-[#a78bfa] font-medium">by {{ book.author }}</p>
            <p class="text-xs text-[#9e9eb8] mt-1.5 line-clamp-2">{{ book.description }}</p>
          </div>

          <div class="flex items-center justify-between pt-2 border-t border-white/[0.06]">
            <span class="text-[11px] text-[#7c7c9a]">{{ book.download_count }} downloads</span>
            <button
              @click="downloadEbook(book)"
              class="px-4 py-2 rounded-xl bg-gradient-to-r from-[#6c5ce7] to-[#8b5cf6] hover:opacity-95 text-white text-xs font-bold flex items-center gap-1.5 shadow-md shadow-[#6c5ce7]/25 cursor-pointer"
            >
              <Download class="w-3.5 h-3.5" />
              <span>Direct Download</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
