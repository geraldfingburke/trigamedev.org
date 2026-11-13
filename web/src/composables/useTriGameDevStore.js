import { ref, computed } from 'vue'

// Singleton reactive state shared across all components & composables
const activeSection = ref('explore') // explore, events, news, jams, showcase, media, keys, library, cms
const keySubTab = ref('pools') // pools, my-keys
const searchQuery = ref('')
const loading = ref(true)
const toastMessage = ref('')

// User / Auth State (Federated Discord OAuth2)
const currentUser = ref(null)
const discordConfigured = ref(false)

// Data Feeds
const overview = ref({
  events: [],
  news: [],
  jams: [],
  showcase: [],
  gallery: [],
  videos: [],
  stats: { available_keys: 0, total_ebooks: 0, total_showcase: 0 }
})
const keyPools = ref([])
const myKeys = ref([])
const ebooks = ref([])

// Global Modals
const revealedKeyModal = ref(null) // { title, platform, key_code }
const showDonateModal = ref(false)
const showUploadEbookModal = ref(false)
const selectedLightbox = ref(null)

function showToast(msg) {
  toastMessage.value = msg
  setTimeout(() => {
    if (toastMessage.value === msg) toastMessage.value = ''
  }, 4200)
}

async function apiFetch(path, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {})
  }
  const res = await fetch(`/api/go${path}`, {
    ...options,
    headers,
    credentials: 'include'
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(data.error || `Request failed (${res.status})`)
  }
  return data
}

async function loadAllData() {
  loading.value = true
  try {
    const [authData, ovData, poolsData, myKeysData, ebooksData] = await Promise.all([
      apiFetch('/auth/me'),
      apiFetch('/tgd/overview'),
      apiFetch('/tgd/keys/pools'),
      apiFetch('/tgd/keys/my-keys'),
      apiFetch('/tgd/library/ebooks')
    ])
    currentUser.value = authData.user
    discordConfigured.value = authData.discord_configured
    overview.value = ovData
    keyPools.value = poolsData.pools || []
    myKeys.value = myKeysData.claimed_keys || []
    ebooks.value = ebooksData.ebooks || []
  } catch (e) {
    showToast(`Notice: ${e.message}`)
  } finally {
    loading.value = false
  }
}

async function logout() {
  try {
    await apiFetch('/auth/logout', { method: 'POST' })
    currentUser.value = null
    myKeys.value = []
    if (activeSection.value === 'cms') activeSection.value = 'explore'
    showToast('Signed out.')
    await loadAllData()
  } catch (e) {
    showToast(e.message)
  }
}

const activeJam = computed(() => overview.value.jams?.[0] || null)
const nextEvent = computed(() => overview.value.events?.[0] || null)

const filteredShowcase = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return overview.value.showcase || []
  return (overview.value.showcase || []).filter(
    g =>
      g.title.toLowerCase().includes(q) ||
      g.developer_name.toLowerCase().includes(q) ||
      g.engine.toLowerCase().includes(q)
  )
})

const filteredKeyPools = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return keyPools.value
  return keyPools.value.filter(
    p =>
      p.title.toLowerCase().includes(q) ||
      p.platform.toLowerCase().includes(q) ||
      p.description.toLowerCase().includes(q)
  )
})

const filteredEbooks = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return ebooks.value
  return ebooks.value.filter(
    b =>
      b.title.toLowerCase().includes(q) ||
      b.author.toLowerCase().includes(q) ||
      b.category.toLowerCase().includes(q)
  )
})

export function useTriGameDevStore() {
  return {
    // State
    activeSection,
    keySubTab,
    searchQuery,
    loading,
    toastMessage,
    currentUser,
    discordConfigured,
    overview,
    keyPools,
    myKeys,
    ebooks,
    // Modals
    revealedKeyModal,
    showDonateModal,
    showUploadEbookModal,
    selectedLightbox,
    // Computed
    activeJam,
    nextEvent,
    filteredShowcase,
    filteredKeyPools,
    filteredEbooks,
    // Actions
    showToast,
    apiFetch,
    loadAllData,
    logout
  }
}
