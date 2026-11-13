import { ref } from 'vue'
import { useTriGameDevStore } from './useTriGameDevStore'

export function useEbookLibrary() {
  const {
    currentUser,
    showUploadEbookModal,
    showToast,
    apiFetch,
    loadAllData
  } = useTriGameDevStore()

  const ebookForm = ref({
    title: '',
    author: '',
    category: 'Game Design',
    description: '',
    file_format: 'PDF',
    file_size_label: '',
    cover_image_url: '/inspo/rooftop-twilight.jpg'
  })

  async function submitEbookUpload() {
    if (!ebookForm.value.title || !ebookForm.value.author) {
      showToast('Title and author are required.')
      return
    }
    try {
      await apiFetch('/tgd/library/ebooks', {
        method: 'POST',
        body: JSON.stringify(ebookForm.value)
      })
      showUploadEbookModal.value = false
      ebookForm.value = {
        title: '',
        author: '',
        category: 'Game Design',
        description: '',
        file_format: 'PDF',
        file_size_label: '',
        cover_image_url: '/inspo/rooftop-twilight.jpg'
      }
      showToast('Ebook added to the direct-download library!')
      await loadAllData()
    } catch (e) {
      showToast(e.message)
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
    submitEbookUpload,
    downloadEbook
  }
}
