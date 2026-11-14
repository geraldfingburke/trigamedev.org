import { ref, computed, watch } from 'vue'

// Capture token from OAuth redirect query params if present
if (typeof window !== 'undefined') {
  const params = new URLSearchParams(window.location.search)
  const tok = params.get('auth_token')
  const sessionOnly = params.get('session_only')
  if (tok) {
    localStorage.setItem('tgd_auth_token', tok)
    params.delete('auth_token')
    const clean = params.toString()
    window.history.replaceState({}, '', window.location.pathname + (clean ? `?${clean}` : ''))
  } else if (sessionOnly) {
    localStorage.removeItem('tgd_auth_token')
    params.delete('session_only')
    const clean = params.toString()
    window.history.replaceState({}, '', window.location.pathname + (clean ? `?${clean}` : ''))
  }
}

// Singleton reactive state shared across all components & composables
const activeSection = ref('explore') // explore, events, news, jams, showcase, media, keys, library, cms
const keySubTab = ref('pools') // pools, my-keys
const searchQuery = ref('')
const loading = ref(true)
const toastMessage = ref('')

// User / Auth State (Federated Discord OAuth2 + Local Token)
const currentUser = ref(null)
const discordConfigured = ref(false)
const staySignedIn = ref(
  typeof window !== 'undefined' ? localStorage.getItem('tgd_stay_signed_in') !== 'false' : true
)

// Data Feeds
const overview = ref({
  events: [],
  news: [],
  jams: [],
  showcase: [],
  gallery: [],
  videos: [],
  unread: {},
  stats: {
    total_members: 162,
    total_jams: 0,
    released_games: 0,
    available_keys: 0,
    total_projects: 0,
    total_ebooks: 0,
    total_showcase: 0
  }
})
const keyPools = ref([])
const myKeys = ref([])
const ebooks = ref([])

// Global Modals
const revealedKeyModal = ref(null) // { title, platform, key_code }
const showDonateModal = ref(false)
const showUploadEbookModal = ref(false)
const selectedLightbox = ref(null)
const selectedWinnersJam = ref(null)

function showToast(msg) {
  toastMessage.value = msg
  setTimeout(() => {
    if (toastMessage.value === msg) toastMessage.value = ''
  }, 4200)
}

function loginWithDiscord() {
  if (typeof window !== 'undefined') {
    localStorage.setItem('tgd_stay_signed_in', staySignedIn.value ? '1' : '0')
    window.location.href = `/api/go/auth/discord/login?remember=${staySignedIn.value ? '1' : '0'}`
  }
}

async function apiFetch(path, options = {}) {
  const isFormData = typeof FormData !== 'undefined' && options.body instanceof FormData
  const token = typeof window !== 'undefined' ? localStorage.getItem('tgd_auth_token') : null
  const headers = {
    ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
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

async function uploadImageFile(file) {
  const fd = new FormData()
  fd.append('file', file)
  const res = await apiFetch('/tgd/upload/image', {
    method: 'POST',
    body: fd
  })
  return res.url
}

async function markSectionAsRead(sec) {
  if (!currentUser.value || !sec || sec === 'explore' || sec === 'cms') return
  if (overview.value.unread && overview.value.unread[sec]) {
    overview.value.unread[sec] = 0
  }
  try {
    await apiFetch('/tgd/notifications/read', {
      method: 'POST',
      body: JSON.stringify({ section: sec })
    })
  } catch {
    // ignore background read-receipt errors
  }
}

watch(activeSection, newSec => {
  markSectionAsRead(newSec)
})

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
    overview.value = {
      ...overview.value,
      ...ovData,
      unread: ovData.unread || {}
    }
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
    if (typeof window !== 'undefined') {
      localStorage.removeItem('tgd_auth_token')
    }
    await apiFetch('/auth/logout', { method: 'POST' })
    currentUser.value = null
    myKeys.value = []
    if (overview.value.unread) {
      overview.value.unread = {}
    }
    if (activeSection.value === 'cms') activeSection.value = 'explore'
    showToast('Signed out.')
    await loadAllData()
  } catch (e) {
    showToast(e.message)
  }
}

const activeJam = computed(() => overview.value.jams?.[0] || null)

const upcomingEvents = computed(() => {
  const now = Date.now()
  return (overview.value.events || [])
    .filter(ev => new Date(ev.start_time).getTime() >= now)
    .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime())
})

const pastEvents = computed(() => {
  const now = Date.now()
  return (overview.value.events || [])
    .filter(ev => new Date(ev.start_time).getTime() < now)
    .sort((a, b) => new Date(b.start_time).getTime() - new Date(a.start_time).getTime())
})

const nextEvent = computed(() => upcomingEvents.value[0] || overview.value.events?.[0] || null)

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
    staySignedIn,
    overview,
    keyPools,
    myKeys,
    ebooks,
    // Modals
    revealedKeyModal,
    showDonateModal,
    showUploadEbookModal,
    selectedLightbox,
    selectedWinnersJam,
    // Computed
    activeJam,
    upcomingEvents,
    pastEvents,
    nextEvent,
    filteredShowcase,
    filteredKeyPools,
    filteredEbooks,
    // Actions
    showToast,
    loginWithDiscord,
    apiFetch,
    uploadImageFile,
    markSectionAsRead,
    loadAllData,
    logout
  }
}
