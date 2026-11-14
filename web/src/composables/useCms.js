import { ref } from 'vue'
import { useTriGameDevStore } from './useTriGameDevStore'

const scrapingItch = ref(false)
const uploadingCmsImage = ref(false)
const cmsForm = ref({
  section: 'jams', // jams, showcase, news, events, gallery, videos
  jam_type: 'monthly', // monthly, other
  title: '',
  subtitle: '',
  description: '',
  tag: '',
  action_label: 'Play Now on Itch.io',
  url: '',
  secondary_url: '',
  image_url: '',
  start_date: '',
  end_date: '',
  status: 'active',
  participants_count: 0,
  submissions_count: 0,
  winners: []
})

export function useCms() {
  const { showToast, apiFetch, uploadImageFile, loadAllData } = useTriGameDevStore()

  async function autofillFromItch() {
    const rawUrl = (cmsForm.value.url || '').trim()
    if (!rawUrl) {
      showToast('Paste a game or jam URL first.')
      return
    }
    scrapingItch.value = true
    try {
      const data = await apiFetch(`/tgd/cms/itch-scrape?url=${encodeURIComponent(rawUrl)}`)
      cmsForm.value.section = data.kind === 'jam' ? 'jams' : 'showcase'
      cmsForm.value.title = data.title || ''
      cmsForm.value.subtitle = data.subtitle || ''
      cmsForm.value.description = data.description || ''
      cmsForm.value.tag = data.tag || ''
      cmsForm.value.action_label = data.action_label || 'Play Now on Itch.io'
      cmsForm.value.url = data.url || rawUrl
      cmsForm.value.secondary_url = data.secondary_url || ''
      cmsForm.value.image_url = data.image_url || ''
      cmsForm.value.start_date = data.start_date || ''
      cmsForm.value.end_date = data.end_date || ''
      cmsForm.value.status = data.status || 'active'
      cmsForm.value.participants_count = data.participants_count || 0
      cmsForm.value.submissions_count = data.submissions_count || 0
      cmsForm.value.winners = data.winners || []
      showToast(`Imported "${data.title}" (${data.platform || 'Web'})`)
    } catch (e) {
      showToast(`Import error: ${e.message}`)
    } finally {
      scrapingItch.value = false
    }
  }

  async function handleCmsImageUpload(event) {
    const file = event.target.files?.[0]
    if (!file) return
    uploadingCmsImage.value = true
    try {
      const url = await uploadImageFile(file)
      cmsForm.value.image_url = url
      showToast('Image uploaded!')
    } catch (e) {
      showToast(e.message)
    } finally {
      uploadingCmsImage.value = false
    }
  }

  function handleVideoUrlInput() {
    const u = cmsForm.value.url.trim()
    if (!u || cmsForm.value.image_url) return
    const ytMatch = u.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([A-Za-z0-9_-]{11})/)
    if (ytMatch) {
      cmsForm.value.image_url = `https://img.youtube.com/vi/${ytMatch[1]}/hqdefault.jpg`
    }
  }

  async function publishCMSItem() {
    if (!cmsForm.value.title && !cmsForm.value.url) {
      showToast('Please enter a title or URL.')
      return
    }
    try {
      await apiFetch('/tgd/cms/item', {
        method: 'POST',
        body: JSON.stringify(cmsForm.value)
      })
      showToast(`Published to ${cmsForm.value.section.toUpperCase()}!`)
      cmsForm.value.title = ''
      cmsForm.value.subtitle = ''
      cmsForm.value.description = ''
      cmsForm.value.tag = ''
      cmsForm.value.action_label = 'Play Now on Itch.io'
      cmsForm.value.url = ''
      cmsForm.value.secondary_url = ''
      cmsForm.value.image_url = ''
      cmsForm.value.winners = []
      await loadAllData()
    } catch (e) {
      showToast(e.message)
    }
  }

  async function deleteCMSItem(section, id) {
    try {
      await apiFetch(`/tgd/cms/item?section=${section}&id=${id}`, { method: 'DELETE' })
      showToast(`Removed item #${id} from ${section}.`)
      await loadAllData()
    } catch (e) {
      showToast(e.message)
    }
  }

  async function triggerDiscordSync() {
    try {
      const res = await apiFetch('/tgd/cms/discord-sync', { method: 'POST' })
      showToast(`Discord Sync complete (${res.synced_events} live scheduled events pulled).`)
      await loadAllData()
    } catch (e) {
      showToast(e.message)
    }
  }

  return {
    scrapingItch,
    uploadingCmsImage,
    cmsForm,
    autofillFromItch,
    handleCmsImageUpload,
    handleVideoUrlInput,
    publishCMSItem,
    deleteCMSItem,
    triggerDiscordSync
  }
}
