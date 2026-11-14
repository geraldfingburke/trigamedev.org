import { ref } from 'vue'
import { useTriGameDevStore } from './useTriGameDevStore'

const copiedKeyId = ref(null)
const uploadingCover = ref(false)
const donateForm = ref({
  title: '',
  platform: 'Steam',
  store_url: '',
  cover_image_url: '/brand/filler.png',
  description: '',
  raw_keys: ''
})

export function useKeyVault() {
  const {
    currentUser,
    revealedKeyModal,
    showDonateModal,
    showToast,
    apiFetch,
    uploadImageFile,
    loadAllData
  } = useTriGameDevStore()

  async function handleDonateCoverUpload(event) {
    const file = event.target.files?.[0]
    if (!file) return
    uploadingCover.value = true
    try {
      const url = await uploadImageFile(file)
      donateForm.value.cover_image_url = url
      showToast('Featured image uploaded!')
    } catch (e) {
      showToast(e.message)
    } finally {
      uploadingCover.value = false
    }
  }

  async function claimKeyFromPool(pool) {
    if (!currentUser.value) {
      showToast('Please sign in with Discord to claim keys.')
      return
    }
    try {
      const res = await apiFetch('/tgd/keys/claim', {
        method: 'POST',
        body: JSON.stringify({ pool_id: pool.id })
      })
      revealedKeyModal.value = {
        title: pool.title,
        platform: pool.platform,
        key_code: res.key_code
      }
      showToast(`Redeemed key for ${pool.title}! Saved permanently to 'My Keys'.`)
      await loadAllData()
    } catch (e) {
      showToast(e.message)
    }
  }

  async function submitKeyDonation() {
    const keys = donateForm.value.raw_keys
      .split('\n')
      .map(k => k.trim())
      .filter(Boolean)
    if (!donateForm.value.title || keys.length === 0) {
      showToast('Please enter a project title and at least one key code.')
      return
    }
    try {
      const res = await apiFetch('/tgd/keys/pools', {
        method: 'POST',
        body: JSON.stringify({
          title: donateForm.value.title,
          platform: donateForm.value.platform,
          store_url: donateForm.value.store_url,
          cover_image_url: donateForm.value.cover_image_url,
          description: donateForm.value.description,
          keys
        })
      })
      showDonateModal.value = false
      donateForm.value = {
        title: '',
        platform: 'Steam',
        store_url: '',
        cover_image_url: '/brand/filler.png',
        description: '',
        raw_keys: ''
      }
      showToast(`Added ${res.keys_added} key(s) to ${res.pool_id ? 'project' : 'vault'}!`)
      await loadAllData()
    } catch (e) {
      showToast(e.message)
    }
  }

  function copyCode(code, id) {
    navigator.clipboard.writeText(code)
    copiedKeyId.value = id
    showToast('Key code copied to clipboard!')
    setTimeout(() => {
      if (copiedKeyId.value === id) copiedKeyId.value = null
    }, 2000)
  }

  async function deleteKeyProject(project) {
    if (!project?.id) return
    try {
      await apiFetch(`/tgd/keys/pools?id=${project.id}`, {
        method: 'DELETE'
      })
      showToast(`Removed project "${project.title}".`)
      await loadAllData()
    } catch (e) {
      showToast(e.message)
    }
  }

  return {
    copiedKeyId,
    uploadingCover,
    donateForm,
    handleDonateCoverUpload,
    claimKeyFromPool,
    submitKeyDonation,
    deleteKeyProject,
    copyCode
  }
}
