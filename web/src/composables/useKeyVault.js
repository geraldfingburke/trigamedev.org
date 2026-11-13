import { ref } from 'vue'
import { useTriGameDevStore } from './useTriGameDevStore'

export function useKeyVault() {
  const {
    currentUser,
    revealedKeyModal,
    showDonateModal,
    showToast,
    apiFetch,
    loadAllData
  } = useTriGameDevStore()

  const copiedKeyId = ref(null)
  const donateForm = ref({
    title: '',
    platform: 'Steam',
    store_url: '',
    cover_image_url: '/brand/filler.png',
    description: '',
    raw_keys: ''
  })

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
      showToast('Please enter a game title and at least one key code.')
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
      showToast(`Added ${res.keys_added} key(s) to the TriGameDev Key Vault!`)
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

  return {
    copiedKeyId,
    donateForm,
    claimKeyFromPool,
    submitKeyDonation,
    copyCode
  }
}
