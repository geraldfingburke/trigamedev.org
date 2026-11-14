import { ref } from 'vue'
import { useTriGameDevStore } from './useTriGameDevStore'

const ebookForm = ref({
  title: '',
  author: '',
  category: 'Game Design',
  description: '',
  file_format: 'PDF',
  file_size_label: '',
  cover_image_url: '/brand/filler.png'
})
const selectedBookFile = ref(null)
const coverCandidates = ref([])
const scrapingBook = ref(false)
const uploadingBook = ref(false)
const uploadingCustomCover = ref(false)

export function useEbookLibrary() {
  const {
    currentUser,
    showUploadEbookModal,
    showToast,
    apiFetch,
    uploadImageFile,
    loadAllData
  } = useTriGameDevStore()

  async function scrapeBookByQuery(rawQuery) {
    const q = (rawQuery || ebookForm.value.title || '').trim()
    if (!q) {
      showToast('Enter a book title or select a file first.')
      return
    }
    scrapingBook.value = true
    try {
      const data = await apiFetch(`/tgd/library/book-scrape?q=${encodeURIComponent(q)}`)
      if (data.title) ebookForm.value.title = data.title
      if (data.author && !ebookForm.value.author) ebookForm.value.author = data.author
      if (data.description && !ebookForm.value.description) ebookForm.value.description = data.description
      coverCandidates.value = data.cover_candidates || []
      if (coverCandidates.value.length > 0) {
        ebookForm.value.cover_image_url = coverCandidates.value[0]
      }
    } catch (e) {
      showToast(`Book lookup: ${e.message}`)
    } finally {
      scrapingBook.value = false
    }
  }

  async function handleBookFileSelect(event) {
    const file = event.target.files?.[0]
    if (!file) return
    selectedBookFile.value = file

    const extMatch = file.name.match(/\.([^.]+)$/)
    if (extMatch) {
      ebookForm.value.file_format = extMatch[1].toUpperCase()
    }
    const mb = file.size / (1024 * 1024)
    ebookForm.value.file_size_label = mb >= 1 ? `${mb.toFixed(1)} MB` : `${Math.max(1, Math.round(file.size / 1024))} KB`

    // Clean filename to initial title guess and scrape cover candidates
    const baseName = file.name.replace(/\.[^.]+$/, '').replace(/[_]+/g, ' ')
    ebookForm.value.title = baseName
    await scrapeBookByQuery(file.name)
  }

  async function handleCustomCoverUpload(event) {
    const file = event.target.files?.[0]
    if (!file) return
    uploadingCustomCover.value = true
    try {
      const url = await uploadImageFile(file)
      ebookForm.value.cover_image_url = url
      if (!coverCandidates.value.includes(url)) {
        coverCandidates.value = [url, ...coverCandidates.value]
      }
      showToast('Custom cover image uploaded!')
    } catch (e) {
      showToast(e.message)
    } finally {
      uploadingCustomCover.value = false
    }
  }

  async function submitEbookUpload() {
    if (!selectedBookFile.value && !ebookForm.value.title) {
      showToast('Please select a book file to upload.')
      return
    }
    uploadingBook.value = true
    try {
      const fd = new FormData()
      if (selectedBookFile.value) {
        fd.append('file', selectedBookFile.value)
      }
      fd.append('title', ebookForm.value.title)
      fd.append('author', ebookForm.value.author || 'Unknown Author')
      fd.append('category', ebookForm.value.category || 'Game Design')
      fd.append('description', ebookForm.value.description || '')
      fd.append('file_format', ebookForm.value.file_format || 'PDF')
      fd.append('cover_image_url', ebookForm.value.cover_image_url || '/brand/filler.png')

      await apiFetch('/tgd/library/ebooks', {
        method: 'POST',
        body: fd
      })
      showUploadEbookModal.value = false
      selectedBookFile.value = null
      coverCandidates.value = []
      ebookForm.value = {
        title: '',
        author: '',
        category: 'Game Design',
        description: '',
        file_format: 'PDF',
        file_size_label: '',
        cover_image_url: '/brand/filler.png'
      }
      showToast('Ebook uploaded to the library!')
      await loadAllData()
    } catch (e) {
      showToast(e.message)
    } finally {
      uploadingBook.value = false
    }
  }

  function downloadEbook(book) {
    if (!currentUser.value) {
      showToast('Sign in with your TriGameDev Discord account to download ebooks.')
      return
    }
    window.open(`/api/go/tgd/library/download?id=${book.id}`, '_blank')
    setTimeout(() => loadAllData(), 800)
  }

  return {
    ebookForm,
    selectedBookFile,
    coverCandidates,
    scrapingBook,
    uploadingBook,
    uploadingCustomCover,
    scrapeBookByQuery,
    handleBookFileSelect,
    handleCustomCoverUpload,
    submitEbookUpload,
    downloadEbook
  }
}
