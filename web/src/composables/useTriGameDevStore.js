import { ref, computed, watch } from 'vue'

export function slugify(raw) {
  return (
    String(raw || '')
      .toLowerCase()
      .trim()
      .replace(/['\u2019]/g, '')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '') || 'post'
  )
}

const sectionFromPathSegment = {
  '': 'explore',
  explore: 'explore',
  jams: 'jams',
  events: 'events',
  news: 'news',
  announcements: 'news',
  games: 'showcase',
  showcase: 'showcase',
  media: 'media',
  gallery: 'media',
  videos: 'media',
  keys: 'keys',
  library: 'library',
  ebooks: 'library',
  cms: 'cms'
}

const pathSegmentFromSection = {
  explore: '',
  jams: 'jams',
  events: 'events',
  news: 'news',
  showcase: 'games',
  media: 'media',
  keys: 'keys',
  library: 'library',
  cms: 'cms'
}

// Singleton reactive state shared across all components & composables
const activeSection = ref('explore') // explore, events, news, jams, showcase, media, keys, library, cms
const activeSlug = ref('')
let suppressUrlPush = false

function parseCurrentPath() {
  if (typeof window === 'undefined') return { section: 'explore', slug: '' }
  const parts = window.location.pathname
    .split('/')
    .map(p => p.trim().toLowerCase())
    .filter(Boolean)
  const seg1 = parts[0] || ''
  const sec = sectionFromPathSegment[seg1] || 'explore'
  const slug = parts[1] ? decodeURIComponent(parts[1]) : ''
  return { section: sec, slug }
}

function buildPathForRoute(section, slug = '') {
  const base = pathSegmentFromSection[section] ?? ''
  if (!base) return '/'
  if (slug) return `/${base}/${encodeURIComponent(slug)}`
  return `/${base}`
}

function scrollToActiveSlug() {
  if (typeof window === 'undefined' || !activeSlug.value) return
  setTimeout(() => {
    const el = document.getElementById(`slug-${activeSlug.value}`)
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
  }, 120)
}

function syncRouteFromLocation() {
  const { section, slug } = parseCurrentPath()
  suppressUrlPush = true
  activeSection.value = section
  activeSlug.value = slug
  suppressUrlPush = false
  if (slug) {
    scrollToActiveSlug()
  }
}

function navigateTo(section, slug = '') {
  suppressUrlPush = true
  activeSection.value = section
  activeSlug.value = slug ? slugify(slug) : ''
  suppressUrlPush = false
  if (typeof window !== 'undefined') {
    const targetPath = buildPathForRoute(activeSection.value, activeSlug.value)
    if (window.location.pathname !== targetPath) {
      window.history.pushState({}, '', targetPath)
    }
  }
  if (activeSlug.value) {
    scrollToActiveSlug()
  }
}

// Capture token from OAuth redirect query params if present, and initialize route from URL path
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
  syncRouteFromLocation()
  window.addEventListener('popstate', syncRouteFromLocation)
}
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

watch(activeSection, (newSec, oldSec) => {
  if (!suppressUrlPush && newSec !== oldSec) {
    activeSlug.value = ''
    if (typeof window !== 'undefined') {
      const targetPath = buildPathForRoute(newSec, '')
      if (window.location.pathname !== targetPath) {
        window.history.pushState({}, '', targetPath)
      }
    }
  }
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
    if (!currentUser.value && ['keys', 'library', 'cms'].includes(activeSection.value) && !activeSlug.value) {
      activeSection.value = 'explore'
    }
    overview.value = {
      ...overview.value,
      ...ovData,
      unread: ovData.unread || {}
    }
    keyPools.value = poolsData.pools || []
    myKeys.value = myKeysData.claimed_keys || []
    ebooks.value = ebooksData.ebooks || []
    if (activeSlug.value) {
      scrollToActiveSlug()
    }
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
    if (['keys', 'library', 'cms'].includes(activeSection.value)) {
      activeSection.value = 'explore'
    }
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
  const base = !q
    ? [...(overview.value.showcase || [])]
    : (overview.value.showcase || []).filter(
        g =>
          g.title.toLowerCase().includes(q) ||
          g.developer_name.toLowerCase().includes(q) ||
          g.engine.toLowerCase().includes(q)
      )
  if (activeSlug.value) {
    base.sort((a, b) => {
      const aMatch = (a.slug || slugify(a.title)) === activeSlug.value ? -1 : 0
      const bMatch = (b.slug || slugify(b.title)) === activeSlug.value ? -1 : 0
      return aMatch - bMatch
    })
  }
  return base
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
    activeSlug,
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
    slugify,
    navigateTo,
    showToast,
    loginWithDiscord,
    apiFetch,
    uploadImageFile,
    markSectionAsRead,
    loadAllData,
    logout
  }
}
