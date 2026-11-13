import { ref } from 'vue'
import { useTriGameDevStore } from './useTriGameDevStore'

export function useCms() {
  const { showToast, apiFetch, loadAllData } = useTriGameDevStore()

  const scrapingItch = ref(false)
  const cmsForm = ref({
    section: 'jams',
    title: '',
    subtitle: '',
    description: '',
    tag: '',
    url: '',
    secondary_url: '',
    image_url: '',
    start_date: '',
    end_date: '',
    status: 'active',
    participants_count: 0,
    submissions_count: 0
  })

  async function autofillFromItch() {
    if (!cmsForm.value.url || !cmsForm.value.url.includes('itch.io')) {
      showToast('Please paste an itch.io Jam or Game URL first.')
      return
    }
    scrapingItch.value = true
    try {
      const data = await apiFetch(`/tgd/cms/itch-scrape?url=${encodeURIComponent(cmsForm.value.url)}`)
      cmsForm.value.section = data.kind === 'jam' ? 'jams' : 'showcase'
      cmsForm.value.title = data.title || ''
      cmsForm.value.subtitle = data.subtitle || ''
      cmsForm.value.description = data.description || ''
      cmsForm.value.tag = data.tag || ''
      cmsForm.value.url = data.url || cmsForm.value.url
      cmsForm.value.secondary_url = data.secondary_url || ''
      cmsForm.value.image_url = data.image_url || ''
      cmsForm.value.start_date = data.start_date || ''
      cmsForm.value.end_date = data.end_date || ''
      cmsForm.value.status = data.status || 'active'
      cmsForm.value.participants_count = data.participants_count || 0
      cmsForm.value.submissions_count = data.submissions_count || 0
      showToast(`Scraped ${data.kind.toUpperCase()}: "${data.title}" from itch.io!`)
    } catch (e) {
      showToast(`Scraper error: ${e.message}`)
    } finally {
      scrapingItch.value = false
    }
  }

  async function publishCMSItem() {
    if (!cmsForm.value.title && !cmsForm.value.url.includes('itch.io')) {
      showToast('Please enter an itch.io URL or a title.')
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
      cmsForm.value.url = ''
      cmsForm.value.secondary_url = ''
      cmsForm.value.image_url = ''
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
    cmsForm,
    autofillFromItch,
    publishCMSItem,
    deleteCMSItem,
    triggerDiscordSync
  }
}
