<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  Compass,
  Flame,
  Calendar,
  Newspaper,
  Image as ImageIcon,
  KeyRound,
  BookOpen,
  ShieldAlert,
  Search,
  Bell,
  Sparkles,
  Copy,
  Check,
  Download,
  Plus,
  RefreshCw,
  LogOut,
  Eye,
  Gamepad2,
  ChevronRight,
  Radio,
  Gift,
  Lock
} from 'lucide-vue-next'

// Navigation & UI State
const activeSection = ref('explore') // explore, events, news, jams, showcase, media, keys, library, cms
const keySubTab = ref('pools') // pools, my-keys
const searchQuery = ref('')
const loading = ref(true)
const toastMessage = ref('')

// User / Auth State (Federated Discord OAuth2 Only)
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

// Modals & Interactive States
const revealedKeyModal = ref(null)
const copiedKeyId = ref(null)
const showDonateModal = ref(false)
const showUploadEbookModal = ref(false)
const selectedLightbox = ref(null)

// Forms
const donateForm = ref({
  title: '',
  platform: 'Steam',
  store_url: '',
  cover_image_url: '/inspo/sunset-studio.jpg',
  description: '',
  raw_keys: ''
})

const ebookForm = ref({
  title: '',
  author: '',
  category: 'Game Design',
  description: '',
  file_format: 'PDF',
  file_size_label: '',
  cover_image_url: '/inspo/rooftop-twilight.jpg'
})

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
const scrapingItch = ref(false)

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
      cover_image_url: '/inspo/sunset-studio.jpg',
      description: '',
      raw_keys: ''
    }
    showToast(`Added ${res.keys_added} key(s) to the TriGameDev Key Vault!`)
    await loadAllData()
  } catch (e) {
    showToast(e.message)
  }
}

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

function copyCode(code, id) {
  navigator.clipboard.writeText(code)
  copiedKeyId.value = id
  showToast('Key code copied to clipboard!')
  setTimeout(() => {
    if (copiedKeyId.value === id) copiedKeyId.value = null
  }, 2000)
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

onMounted(() => {
  const params = new URLSearchParams(window.location.search)
  if (params.get('error') === 'discord_oauth_not_configured') {
    showToast('Discord OAuth credentials are not set in /opt/infra/.env yet.')
    window.history.replaceState({}, '', '/')
  }
  loadAllData()
})
</script>

<template>
  <div class="min-h-screen bg-[#12121b] text-[#f4f4f9] flex flex-col p-2 sm:p-4 lg:p-6">
    <!-- Toast Notification -->
    <transition name="fade">
      <div
        v-if="toastMessage"
        class="fixed bottom-6 right-6 z-50 flex items-center gap-3 px-5 py-3.5 rounded-2xl bg-[#252338]/95 border border-[#ff4d8d]/50 text-sm font-medium shadow-2xl shadow-[#ff4d8d]/20 backdrop-blur-xl"
      >
        <Sparkles class="w-4 h-4 text-[#ff4d8d] shrink-0" />
        <span>{{ toastMessage }}</span>
      </div>
    </transition>

    <!-- Main Dashboard App Frame -->
    <div
      class="flex-1 w-full max-w-[1640px] mx-auto bg-[#1b1b26] rounded-[28px] border border-white/[0.07] shadow-[0_30px_90px_rgba(0,0,0,0.65)] overflow-hidden grid grid-cols-1 lg:grid-cols-[76px_1fr_340px]"
    >
      <!-- 1. LEFT VERTICAL ICON DOCK / NAVIGATION RAIL -->
      <aside
        class="bg-[#232234] border-b lg:border-b-0 lg:border-r border-white/[0.06] flex lg:flex-col items-center justify-between px-4 py-3 lg:py-6 gap-2"
      >
        <!-- Brand Emblem -->
        <button
          @click="activeSection = 'explore'"
          class="w-11 h-11 rounded-2xl bg-gradient-to-br from-[#6c5ce7] via-[#8b5cf6] to-[#ff4d8d] flex items-center justify-center shadow-lg shadow-[#6c5ce7]/30 hover:scale-105 transition cursor-pointer"
          title="TriGameDev.org Home"
        >
          <Gamepad2 class="w-6 h-6 text-white" />
        </button>

        <!-- Primary Icon Rail -->
        <nav class="flex lg:flex-col items-center gap-2 sm:gap-3 overflow-x-auto">
          <button
            v-for="item in [
              { id: 'explore', label: 'Explore Hub', icon: Compass },
              { id: 'jams', label: 'Monthly Game Jam', icon: Flame },
              { id: 'events', label: 'Events & Meetups', icon: Calendar },
              { id: 'news', label: 'News & Announcements', icon: Newspaper },
              { id: 'showcase', label: 'Member Showcase', icon: Sparkles },
              { id: 'media', label: 'Gallery & Videos', icon: ImageIcon },
              { id: 'keys', label: 'Member Key Vault', icon: KeyRound, badge: overview.stats?.available_keys },
              { id: 'library', label: 'Ebook Library', icon: BookOpen }
            ]"
            :key="item.id"
            @click="activeSection = item.id"
            class="relative group w-11 h-11 rounded-2xl flex items-center justify-center transition-all cursor-pointer"
            :class="
              activeSection === item.id
                ? 'bg-[#6c5ce7]/20 text-white shadow-inner border border-[#8b5cf6]/40'
                : 'text-[#8e8ea9] hover:text-white hover:bg-white/[0.05]'
            "
            :title="item.label"
          >
            <span
              v-if="activeSection === item.id"
              class="hidden lg:block absolute -left-3 w-1.5 h-6 rounded-r-full bg-gradient-to-b from-[#8b5cf6] to-[#ff4d8d]"
            ></span>
            <component :is="item.icon" class="w-5 h-5" />
            <span
              v-if="item.badge"
              class="absolute -top-1 -right-1 px-1.5 py-0.2 rounded-full text-[10px] font-bold bg-[#ff4d8d] text-white"
            >
              {{ item.badge }}
            </span>
          </button>
        </nav>

        <!-- OG Admin CMS Button at Bottom of Rail (Only shown if logged in as OG Admin) -->
        <div class="flex items-center">
          <button
            v-if="currentUser?.is_og_admin"
            @click="activeSection = 'cms'"
            class="relative w-11 h-11 rounded-2xl flex items-center justify-center transition cursor-pointer"
            :class="
              activeSection === 'cms'
                ? 'bg-[#ff9f43]/25 text-[#ff9f43] border border-[#ff9f43]/50'
                : 'text-[#8e8ea9] hover:text-[#ff9f43] hover:bg-white/[0.05]'
            "
            title="OG Role Admin CMS"
          >
            <ShieldAlert class="w-5 h-5" />
          </button>
        </div>
      </aside>

      <!-- 2. CENTER MAIN STAGE -->
      <main class="p-5 sm:p-8 flex flex-col gap-7 overflow-y-auto">
        <!-- Top Search Bar -->
        <header class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4">
          <div class="relative flex-1 max-w-xl">
            <Search class="w-4 h-4 text-[#8e8ea9] absolute left-4 top-1/2 -translate-y-1/2" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search TriGameDev games, keys, ebooks, or events..."
              class="w-full pl-11 pr-4 py-3 rounded-2xl bg-[#252438] border border-white/[0.06] text-sm text-white placeholder-[#7c7c9a] focus:outline-none focus:border-[#8b5cf6] transition"
            />
          </div>
        </header>

        <!-- Horizontal Category Underline Tabs -->
        <div class="flex items-center justify-between border-b border-white/[0.07] pb-1 overflow-x-auto gap-6">
          <div class="flex items-center gap-7 text-sm font-medium shrink-0">
            <button
              v-for="tab in [
                { id: 'explore', label: 'Explore' },
                { id: 'events', label: 'Events' },
                { id: 'news', label: 'Announcements' },
                { id: 'jams', label: 'Monthly Jam' },
                { id: 'showcase', label: 'Showcase' },
                { id: 'media', label: 'Gallery & Videos' },
                { id: 'keys', label: 'Key Vault' },
                { id: 'library', label: 'Ebook Library' },
                ...(currentUser?.is_og_admin ? [{ id: 'cms', label: '⚡ OG CMS' }] : [])
              ]"
              :key="tab.id"
              @click="activeSection = tab.id"
              class="pb-3 relative transition cursor-pointer whitespace-nowrap"
              :class="
                activeSection === tab.id
                  ? 'text-white font-semibold'
                  : 'text-[#8e8ea9] hover:text-white'
              "
            >
              {{ tab.label }}
              <span
                v-if="activeSection === tab.id"
                class="absolute bottom-0 left-0 right-0 h-[3px] rounded-full bg-gradient-to-r from-[#6c5ce7] via-[#8b5cf6] to-[#ff4d8d]"
              ></span>
            </button>
          </div>
        </div>

        <!-- =====================================================================
             VIEW 1: EXPLORE HUB
             ===================================================================== -->
        <div v-if="activeSection === 'explore'" class="flex flex-col gap-8">
          <!-- Full-Width Lofi Sunset Hero Banner -->
          <section
            class="relative rounded-3xl overflow-hidden border border-white/[0.1] min-h-[230px] flex items-center p-6 sm:p-10 shadow-2xl"
          >
            <img
              src="/inspo/sunset-studio.jpg"
              alt="TriGameDev Lofi Sunset Banner"
              class="absolute inset-0 w-full h-full object-cover object-center scale-105 filter brightness-90"
            />
            <div
              class="absolute inset-0 bg-gradient-to-r from-[#1b1838]/95 via-[#433089]/80 to-[#ff4d8d]/35"
            ></div>

            <div class="relative z-10 max-w-xl flex flex-col gap-4">
              <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-black/40 backdrop-blur-md border border-white/15 w-fit text-xs font-semibold text-[#ff7eb3]">
                <Radio class="w-3.5 h-3.5 text-[#ff4d8d]" />
                <span>{{ activeJam ? `${activeJam.month_label} • LIVE JAM` : 'TRIGAMEDEV COLLECTIVE' }}</span>
              </div>

              <h1 class="text-2xl sm:text-4xl font-extrabold tracking-tight leading-tight text-white">
                {{ activeJam ? activeJam.title : 'TriGameDev Community Hub' }}
              </h1>

              <p class="text-sm sm:text-base text-white/85 leading-relaxed">
                <template v-if="activeJam">
                  Theme: <span class="font-mono font-semibold text-[#ffd166]">{{ activeJam.theme }}</span> —
                  {{ activeJam.description }}
                </template>
                <template v-else>
                  Connect with local & indie game developers on Discord. Share builds, join monthly jams, swap keys, and explore community resources.
                </template>
              </p>

              <div class="flex flex-wrap items-center gap-3 pt-1">
                <button
                  @click="activeSection = 'jams'"
                  class="px-6 py-2.5 rounded-xl bg-white text-[#18162b] font-bold text-sm shadow-lg hover:bg-[#f0ecff] transition cursor-pointer"
                >
                  Monthly Jams
                </button>
                <button
                  @click="activeSection = 'keys'"
                  class="px-5 py-2.5 rounded-xl bg-black/45 hover:bg-black/65 text-white border border-white/20 backdrop-blur-md font-semibold text-sm flex items-center gap-2 transition cursor-pointer"
                >
                  <KeyRound class="w-4 h-4 text-[#00f5d4]" />
                  <span>Key Vault ({{ overview.stats?.available_keys || 0 }} Available)</span>
                </button>
              </div>
            </div>
          </section>

          <!-- 3-Column Featured Showcase Cards -->
          <section class="flex flex-col gap-4">
            <div class="flex items-center justify-between">
              <h2 class="text-xl font-bold text-white flex items-center gap-2">
                <span>Member Showcase</span>
                <span class="text-xs px-2.5 py-0.5 rounded-full bg-[#6c5ce7]/25 text-[#a78bfa] font-semibold">
                  Community Builds
                </span>
              </h2>
              <button
                @click="activeSection = 'showcase'"
                class="px-4 py-2 rounded-xl bg-[#6c5ce7] hover:bg-[#5b4bc4] text-xs font-semibold text-white flex items-center gap-1.5 transition cursor-pointer shadow-md shadow-[#6c5ce7]/25"
              >
                <span>View All</span>
                <ChevronRight class="w-4 h-4" />
              </button>
            </div>

            <div
              v-if="filteredShowcase.length === 0"
              class="p-8 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
            >
              No showcase games published yet. OG Admins can add projects from the CMS.
            </div>

            <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
              <article
                v-for="game in filteredShowcase.slice(0, 3)"
                :key="game.id"
                class="group bg-[#242336] rounded-3xl p-3.5 border border-white/[0.06] hover:border-[#8b5cf6]/50 transition-all flex flex-col justify-between shadow-xl"
              >
                <div class="relative h-56 rounded-2xl overflow-hidden bg-[#181726]">
                  <img
                    :src="game.thumbnail_url || '/inspo/sunset-studio.jpg'"
                    :alt="game.title"
                    class="w-full h-full object-cover group-hover:scale-105 transition duration-500"
                  />
                  <div
                    class="absolute top-3 right-3 px-2.5 py-1 rounded-full bg-white/90 text-[#14141d] text-[11px] font-bold flex items-center gap-1 shadow"
                  >
                    <Eye class="w-3.5 h-3.5" />
                    <span>{{ game.views_count }}</span>
                  </div>
                  <div
                    class="absolute bottom-3 left-3 right-3 p-3 rounded-2xl bg-[#1c1b29]/80 backdrop-blur-md border border-white/15 flex items-center justify-between gap-2"
                  >
                    <div class="min-w-0">
                      <h3 class="font-bold text-sm text-white truncate">{{ game.title }}</h3>
                      <p class="text-[11px] text-[#b4b4cc] truncate">by {{ game.developer_name }}</p>
                    </div>
                    <a
                      :href="game.play_url || '#'"
                      target="_blank"
                      class="w-9 h-9 rounded-full bg-gradient-to-tr from-[#ff2e7e] to-[#ff6b9d] flex items-center justify-center text-white shrink-0 shadow-lg shadow-[#ff2e7e]/40 hover:scale-105 transition"
                    >
                      <Radio class="w-4 h-4" />
                    </a>
                  </div>
                </div>

                <div class="pt-4 px-1 flex flex-col gap-3">
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-[#8e8ea9]">Engine:</span>
                    <span class="font-semibold text-[#00f5d4] flex items-center gap-1">
                      <Sparkles class="w-3.5 h-3.5 text-[#8b5cf6]" />
                      {{ game.engine }}
                    </span>
                  </div>
                  <div class="grid grid-cols-2 gap-2.5">
                    <a
                      :href="game.play_url || '#'"
                      target="_blank"
                      class="py-2 px-3 rounded-xl bg-gradient-to-r from-[#6c5ce7] to-[#8b5cf6] hover:opacity-95 text-white font-semibold text-xs text-center shadow-md shadow-[#6c5ce7]/25 transition"
                    >
                      Play Build
                    </a>
                    <a
                      :href="game.source_url || '#'"
                      target="_blank"
                      class="py-2 px-3 rounded-xl bg-transparent border border-white/20 hover:border-white/40 text-white/90 font-semibold text-xs text-center transition"
                    >
                      Source / Info
                    </a>
                  </div>
                </div>
              </article>
            </div>
          </section>

          <!-- 7-Column Compact Gallery Strip -->
          <section v-if="overview.gallery && overview.gallery.length > 0" class="flex flex-col gap-3">
            <div class="flex items-center justify-between">
              <h2 class="text-base font-bold text-white">Recent Gallery Snaps</h2>
              <button
                @click="activeSection = 'media'"
                class="text-xs font-semibold text-[#a78bfa] hover:underline cursor-pointer"
              >
                Open Gallery & Videos →
              </button>
            </div>

            <div class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-7 gap-3">
              <button
                v-for="snap in overview.gallery.slice(0, 7)"
                :key="snap.id"
                @click="selectedLightbox = snap"
                class="group relative aspect-square rounded-2xl overflow-hidden bg-[#242336] border border-white/[0.07] hover:border-[#ff4d8d] transition cursor-pointer"
              >
                <img
                  :src="snap.image_url"
                  :alt="snap.title"
                  class="w-full h-full object-cover group-hover:scale-110 transition duration-300"
                />
              </button>
            </div>
          </section>

          <!-- Two-Column Quick Feed: Latest News & Member Portals -->
          <section class="grid grid-cols-1 xl:grid-cols-2 gap-5">
            <!-- Latest News & Discord Announcements -->
            <div class="bg-[#242336] rounded-3xl p-5 border border-white/[0.06] flex flex-col gap-4">
              <div class="flex items-center justify-between">
                <h3 class="font-bold text-base text-white flex items-center gap-2">
                  <Newspaper class="w-4 h-4 text-[#ff9f43]" />
                  <span>News & Announcements</span>
                </h3>
                <button
                  @click="activeSection = 'news'"
                  class="text-xs text-[#a78bfa] hover:underline cursor-pointer"
                >
                  View All
                </button>
              </div>

              <div v-if="!overview.news || overview.news.length === 0" class="text-xs text-[#8e8ea9] py-6 text-center">
                No announcements posted yet.
              </div>

              <div v-else class="flex flex-col gap-3">
                <div
                  v-for="post in overview.news.slice(0, 3)"
                  :key="post.id"
                  class="p-3.5 rounded-2xl bg-[#1b1a29] border border-white/[0.04] flex flex-col gap-1"
                >
                  <div class="flex items-center justify-between text-[11px]">
                    <span class="px-2 py-0.5 rounded-md bg-[#6c5ce7]/20 text-[#a78bfa] font-semibold">
                      {{ post.tag }}
                    </span>
                    <span class="text-[#7c7c9a]">{{ post.published_at }} • {{ post.author_name }}</span>
                  </div>
                  <h4 class="font-bold text-sm text-white mt-1">{{ post.title }}</h4>
                  <p class="text-xs text-[#9e9eb8] line-clamp-2">{{ post.summary }}</p>
                </div>
              </div>
            </div>

            <!-- Member Vault Quick Access -->
            <div class="bg-[#242336] rounded-3xl p-5 border border-white/[0.06] flex flex-col justify-between gap-4">
              <div class="flex flex-col gap-2">
                <div class="flex items-center justify-between">
                  <span class="px-2.5 py-1 rounded-full text-[11px] font-bold bg-[#00f5d4]/15 text-[#00f5d4]">
                    DISCORD MEMBER PERKS
                  </span>
                </div>
                <h3 class="text-lg font-bold text-white mt-1">
                  Key Sharing Portal & Ebook Library
                </h3>
                <p class="text-xs text-[#9e9eb8] leading-relaxed">
                  Sign in with your TriGameDev Discord account to donate or claim game keys into your permanent
                  <strong class="text-white">My Keys</strong> vault and download community ebooks.
                </p>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <button
                  @click="activeSection = 'keys'"
                  class="p-4 rounded-2xl bg-gradient-to-br from-[#2d294e] to-[#1f1d36] border border-[#8b5cf6]/30 hover:border-[#8b5cf6] text-left transition cursor-pointer group"
                >
                  <KeyRound class="w-5 h-5 text-[#ff4d8d] mb-2 group-hover:scale-110 transition" />
                  <div class="font-bold text-sm text-white">Key Sharing Vault</div>
                  <div class="text-xs text-[#00f5d4] mt-0.5">
                    {{ overview.stats?.available_keys || 0 }} keys available →
                  </div>
                </button>

                <button
                  @click="activeSection = 'library'"
                  class="p-4 rounded-2xl bg-gradient-to-br from-[#2b2540] to-[#1d1b2f] border border-[#ff4d8d]/30 hover:border-[#ff4d8d] text-left transition cursor-pointer group"
                >
                  <BookOpen class="w-5 h-5 text-[#ff9f43] mb-2 group-hover:scale-110 transition" />
                  <div class="font-bold text-sm text-white">Ebook Library</div>
                  <div class="text-xs text-[#ff9f43] mt-0.5">
                    {{ overview.stats?.total_ebooks || 0 }} direct-download books →
                  </div>
                </button>
              </div>
            </div>
          </section>

          <!-- Bottom Discord Connect Bar -->
          <section
            v-if="!currentUser"
            class="p-4 sm:px-6 rounded-2xl bg-[#242336] border border-white/[0.07] flex flex-col sm:flex-row items-center justify-between gap-4"
          >
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 rounded-full bg-[#ff4d8d]"></div>
              <p class="text-xs sm:text-sm text-[#c4c4dc]">
                Member of the <strong class="text-white">TriGameDev Discord</strong>? Sign in to unlock the Key Vault, Ebook Library, and OG CMS.
              </p>
            </div>
            <a
              href="/api/go/auth/discord/login"
              class="px-5 py-2.5 rounded-xl bg-[#5865F2] hover:bg-[#4752c4] text-white text-xs font-bold whitespace-nowrap shadow-md transition"
            >
              Sign In with Discord
            </a>
          </section>
        </div>

        <!-- =====================================================================
             VIEW 2: KEY SHARING PORTAL
             ===================================================================== -->
        <div v-else-if="activeSection === 'keys'" class="flex flex-col gap-6">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 class="text-2xl font-extrabold text-white flex items-center gap-2.5">
                <KeyRound class="w-6 h-6 text-[#ff4d8d]" />
                <span>Member Key Sharing Portal</span>
              </h1>
              <p class="text-xs sm:text-sm text-[#9e9eb8] mt-1">
                Single redemption per pool — claiming a key permanently ties it to your account and moves it to
                <strong class="text-white">My Keys</strong>.
              </p>
            </div>

            <div class="flex items-center gap-2.5">
              <div class="bg-[#232234] p-1 rounded-xl border border-white/10 flex items-center">
                <button
                  @click="keySubTab = 'pools'"
                  class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer"
                  :class="keySubTab === 'pools' ? 'bg-[#6c5ce7] text-white' : 'text-[#8e8ea9] hover:text-white'"
                >
                  Available Pools ({{ keyPools.length }})
                </button>
                <button
                  @click="keySubTab = 'my-keys'"
                  class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer flex items-center gap-1.5"
                  :class="keySubTab === 'my-keys' ? 'bg-[#ff4d8d] text-white' : 'text-[#8e8ea9] hover:text-white'"
                >
                  <span>My Keys</span>
                  <span class="px-1.5 py-0.2 rounded-full bg-black/30 text-[10px]">{{ myKeys.length }}</span>
                </button>
              </div>

              <button
                v-if="currentUser"
                @click="showDonateModal = true"
                class="px-4 py-2 rounded-xl bg-gradient-to-r from-[#ff4d8d] to-[#ff7eb3] text-white text-xs font-bold flex items-center gap-1.5 shadow-lg shadow-[#ff4d8d]/25 cursor-pointer"
              >
                <Plus class="w-4 h-4" />
                <span>Donate Keys</span>
              </button>
            </div>
          </div>

          <div
            v-if="!currentUser"
            class="p-6 rounded-3xl bg-[#242336] border border-[#ff4d8d]/40 flex flex-col sm:flex-row items-center justify-between gap-4"
          >
            <div class="flex items-center gap-4">
              <Lock class="w-8 h-8 text-[#ff4d8d] shrink-0" />
              <div>
                <h3 class="font-bold text-white">Discord Member Login Required</h3>
                <p class="text-xs text-[#9e9eb8]">
                  Sign in with your TriGameDev Discord account to donate keys, redeem keys, or view your claimed keys.
                </p>
              </div>
            </div>
            <a
              href="/api/go/auth/discord/login"
              class="px-4 py-2 rounded-xl bg-[#5865F2] text-white text-xs font-bold shrink-0"
            >
              Sign In with Discord
            </a>
          </div>

          <div v-if="keySubTab === 'pools'">
            <div
              v-if="filteredKeyPools.length === 0"
              class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
            >
              No key pools have been added yet.
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div
                v-for="pool in filteredKeyPools"
                :key="pool.id"
                class="bg-[#242336] rounded-3xl p-4 border border-white/[0.07] flex flex-col sm:flex-row gap-4 justify-between"
              >
                <img
                  :src="pool.cover_image_url || '/inspo/sunset-studio.jpg'"
                  :alt="pool.title"
                  class="w-full sm:w-36 h-36 rounded-2xl object-cover shrink-0"
                />
                <div class="flex-1 flex flex-col justify-between gap-3">
                  <div>
                    <div class="flex items-center justify-between gap-2">
                      <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-[#6c5ce7]/25 text-[#a78bfa]">
                        {{ pool.platform }}
                      </span>
                      <span
                        class="text-xs font-mono font-semibold"
                        :class="pool.available_count > 0 ? 'text-[#00f5d4]' : 'text-[#ff4d8d]'"
                      >
                        {{ pool.available_count }} / {{ pool.total_count }} Keys Left
                      </span>
                    </div>
                    <h3 class="font-bold text-base text-white mt-1.5">{{ pool.title }}</h3>
                    <p class="text-xs text-[#9e9eb8] line-clamp-2 mt-1">{{ pool.description }}</p>
                  </div>

                  <div class="flex items-center justify-between pt-2 border-t border-white/[0.06]">
                    <span class="text-[11px] text-[#7c7c9a]">Donated by {{ pool.donated_by_name }}</span>

                    <button
                      v-if="pool.claimed_by_user"
                      @click="keySubTab = 'my-keys'"
                      class="px-3.5 py-1.5 rounded-xl bg-[#00f5d4]/15 text-[#00f5d4] border border-[#00f5d4]/30 text-xs font-bold cursor-pointer"
                    >
                      ✓ Claimed (View in My Keys)
                    </button>
                    <button
                      v-else-if="pool.available_count > 0"
                      @click="claimKeyFromPool(pool)"
                      class="px-4 py-1.5 rounded-xl bg-gradient-to-r from-[#6c5ce7] to-[#8b5cf6] hover:opacity-95 text-white text-xs font-bold shadow-md shadow-[#6c5ce7]/30 cursor-pointer"
                    >
                      Claim 1 Key
                    </button>
                    <span v-else class="text-xs font-semibold text-[#7c7c9a]">Pool Depleted</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="flex flex-col gap-4">
            <div
              v-if="myKeys.length === 0"
              class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center flex flex-col items-center gap-3"
            >
              <Gift class="w-10 h-10 text-[#ff4d8d]" />
              <h3 class="font-bold text-base text-white">You haven't claimed any keys yet</h3>
              <p class="text-xs text-[#9e9eb8] max-w-md">
                Keys you redeem from community pools are permanently bound to your account and displayed here.
              </p>
            </div>

            <div
              v-for="item in myKeys"
              :key="item.id"
              class="p-4 rounded-3xl bg-[#242336] border border-[#00f5d4]/30 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
            >
              <div class="flex items-center gap-4">
                <img
                  :src="item.cover_image_url || '/inspo/sunset-studio.jpg'"
                  class="w-16 h-16 rounded-2xl object-cover"
                />
                <div>
                  <div class="flex items-center gap-2">
                    <span class="px-2 py-0.5 rounded-md bg-[#6c5ce7]/25 text-[#a78bfa] text-[10px] font-bold">
                      {{ item.platform }}
                    </span>
                    <span class="text-[11px] text-[#7c7c9a]">Claimed {{ item.claimed_at }}</span>
                  </div>
                  <h4 class="font-bold text-base text-white mt-0.5">{{ item.title }}</h4>
                  <p class="text-xs text-[#8e8ea9]">Donated by {{ item.donated_by_name }}</p>
                </div>
              </div>

              <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
                <code
                  class="px-3.5 py-2 rounded-xl bg-[#14141d] border border-white/10 text-xs sm:text-sm font-mono text-[#00f5d4] tracking-wider"
                >
                  {{ item.key_code }}
                </code>
                <button
                  @click="copyCode(item.key_code, item.id)"
                  class="p-2.5 rounded-xl bg-[#6c5ce7] hover:bg-[#5a4ad1] text-white transition cursor-pointer"
                  title="Copy Key"
                >
                  <Check v-if="copiedKeyId === item.id" class="w-4 h-4" />
                  <Copy v-else class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- =====================================================================
             VIEW 3: EBOOK LIBRARY
             ===================================================================== -->
        <div v-else-if="activeSection === 'library'" class="flex flex-col gap-6">
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

        <!-- =====================================================================
             VIEW 4: EVENTS, NEWS, JAMS, SHOWCASE, MEDIA
             ===================================================================== -->
        <div v-else-if="activeSection === 'events'" class="flex flex-col gap-5">
          <h1 class="text-2xl font-extrabold text-white">Upcoming Events & Discord Stages</h1>
          <div
            v-if="!overview.events || overview.events.length === 0"
            class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
          >
            No upcoming events scheduled yet.
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            <div
              v-for="ev in overview.events"
              :key="ev.id"
              class="bg-[#242336] rounded-3xl overflow-hidden border border-white/[0.07] flex flex-col justify-between"
            >
              <img :src="ev.cover_image_url || '/inspo/sunset-studio.jpg'" class="h-40 w-full object-cover" />
              <div class="p-5 flex-1 flex flex-col justify-between gap-4">
                <div>
                  <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-[#ff4d8d]/20 text-[#ff7eb3]">
                    {{ ev.tag }}
                  </span>
                  <h3 class="font-bold text-base text-white mt-2">{{ ev.title }}</h3>
                  <p class="text-xs text-[#9e9eb8] mt-1">{{ ev.description }}</p>
                </div>
                <div class="pt-3 border-t border-white/[0.06] text-xs text-[#00f5d4] flex items-center justify-between">
                  <span>📍 {{ ev.location }}</span>
                  <span>{{ ev.rsvp_count }} Going</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="activeSection === 'news'" class="flex flex-col gap-5">
          <h1 class="text-2xl font-extrabold text-white">Community News & Announcements</h1>
          <div
            v-if="!overview.news || overview.news.length === 0"
            class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
          >
            No news or announcements published yet.
          </div>
          <div v-else class="flex flex-col gap-4">
            <article
              v-for="post in overview.news"
              :key="post.id"
              class="p-6 rounded-3xl bg-[#242336] border border-white/[0.07] flex flex-col gap-2"
            >
              <div class="flex items-center justify-between text-xs">
                <span class="px-2.5 py-0.5 rounded-full bg-[#6c5ce7]/25 text-[#a78bfa] font-bold">{{ post.tag }}</span>
                <span class="text-[#8e8ea9]">{{ post.published_at }} • by {{ post.author_name }} ({{ post.author_role }})</span>
              </div>
              <h2 class="text-lg font-bold text-white mt-1">{{ post.title }}</h2>
              <p class="text-sm text-[#c4c4dc] leading-relaxed">{{ post.body_markdown }}</p>
            </article>
          </div>
        </div>

        <div v-else-if="activeSection === 'jams'" class="flex flex-col gap-5">
          <h1 class="text-2xl font-extrabold text-white">Monthly Game Jams</h1>
          <div
            v-if="!overview.jams || overview.jams.length === 0"
            class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
          >
            No game jams posted yet.
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div
              v-for="jam in overview.jams"
              :key="jam.id"
              class="bg-[#242336] rounded-3xl overflow-hidden border border-white/[0.08] flex flex-col"
            >
              <img :src="jam.cover_image_url || '/inspo/sunset-studio.jpg'" class="h-48 w-full object-cover" />
              <div class="p-6 flex flex-col gap-3">
                <div class="flex items-center justify-between">
                  <span class="text-xs font-bold uppercase tracking-wider text-[#ff4d8d]">{{ jam.month_label }}</span>
                  <span class="px-2.5 py-0.5 rounded-full bg-[#00f5d4]/15 text-[#00f5d4] text-xs font-bold uppercase">
                    {{ jam.status }}
                  </span>
                </div>
                <h2 class="text-xl font-extrabold text-white">{{ jam.title }}</h2>
                <div class="text-xs font-mono text-[#ffd166]">THEME: {{ jam.theme }}</div>
                <p class="text-xs text-[#9e9eb8] leading-relaxed">{{ jam.description }}</p>
                <div class="flex items-center justify-between pt-3 border-t border-white/[0.06] text-xs">
                  <span class="text-[#8e8ea9]">{{ jam.participants_count }} Jammers • {{ jam.submissions_count }} Entries</span>
                  <a
                    :href="jam.itch_url"
                    target="_blank"
                    class="px-4 py-2 rounded-xl bg-[#6c5ce7] text-white font-bold"
                  >
                    Itch.io Jam Page →
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="activeSection === 'showcase'" class="flex flex-col gap-5">
          <h1 class="text-2xl font-extrabold text-white">Member Game Showcase</h1>
          <div
            v-if="filteredShowcase.length === 0"
            class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
          >
            No showcase builds published yet.
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
            <article
              v-for="game in filteredShowcase"
              :key="game.id"
              class="bg-[#242336] rounded-3xl p-4 border border-white/[0.07] flex flex-col justify-between gap-4"
            >
              <div class="relative h-48 rounded-2xl overflow-hidden">
                <img :src="game.thumbnail_url || '/inspo/sunset-studio.jpg'" class="w-full h-full object-cover" />
              </div>
              <div>
                <span class="text-[11px] font-bold text-[#00f5d4]">{{ game.engine }}</span>
                <h3 class="font-bold text-base text-white mt-0.5">{{ game.title }}</h3>
                <p class="text-xs text-[#a78bfa]">by {{ game.developer_name }}</p>
                <p class="text-xs text-[#9e9eb8] mt-2">{{ game.description }}</p>
              </div>
              <div class="grid grid-cols-2 gap-2">
                <a :href="game.play_url" target="_blank" class="py-2 rounded-xl bg-[#6c5ce7] text-white text-xs font-bold text-center">
                  Play Game
                </a>
                <a :href="game.source_url" target="_blank" class="py-2 rounded-xl border border-white/20 text-white text-xs font-bold text-center">
                  Source
                </a>
              </div>
            </article>
          </div>
        </div>

        <div v-else-if="activeSection === 'media'" class="flex flex-col gap-7">
          <div>
            <h1 class="text-2xl font-extrabold text-white mb-4">Community Photo Gallery</h1>
            <div
              v-if="!overview.gallery || overview.gallery.length === 0"
              class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
            >
              No gallery photos uploaded yet.
            </div>
            <div v-else class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div
                v-for="img in overview.gallery"
                :key="img.id"
                @click="selectedLightbox = img"
                class="group rounded-3xl overflow-hidden bg-[#242336] border border-white/[0.07] cursor-pointer"
              >
                <img :src="img.image_url" class="h-44 w-full object-cover group-hover:scale-105 transition" />
                <div class="p-3">
                  <div class="text-[10px] font-bold text-[#ff4d8d]">{{ img.event_label }}</div>
                  <div class="text-xs font-bold text-white truncate">{{ img.title }}</div>
                </div>
              </div>
            </div>
          </div>

          <div>
            <h2 class="text-2xl font-extrabold text-white mb-4">Talks & Workshop Videos</h2>
            <div
              v-if="!overview.videos || overview.videos.length === 0"
              class="p-10 rounded-3xl bg-[#242336] border border-white/[0.06] text-center text-sm text-[#8e8ea9]"
            >
              No videos archived yet.
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
              <a
                v-for="vid in overview.videos"
                :key="vid.id"
                :href="vid.video_url"
                target="_blank"
                class="bg-[#242336] rounded-3xl overflow-hidden border border-white/[0.07] hover:border-[#8b5cf6] transition flex flex-col"
              >
                <div class="relative h-40">
                  <img :src="vid.thumbnail_url || '/inspo/sunset-studio.jpg'" class="w-full h-full object-cover" />
                  <span v-if="vid.duration" class="absolute bottom-2 right-2 px-2 py-0.5 rounded bg-black/80 text-[11px] font-mono text-white">
                    {{ vid.duration }}
                  </span>
                </div>
                <div class="p-4">
                  <h3 class="font-bold text-sm text-white">{{ vid.title }}</h3>
                  <p class="text-xs text-[#a78bfa] mt-1">{{ vid.speaker_name }}</p>
                  <p class="text-xs text-[#9e9eb8] mt-1 line-clamp-2">{{ vid.description }}</p>
                </div>
              </a>
            </div>
          </div>
        </div>

        <!-- =====================================================================
             VIEW 5: OG ROLE ADMIN CMS
             ===================================================================== -->
        <div v-else-if="activeSection === 'cms'" class="flex flex-col gap-6">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <span class="px-2.5 py-1 rounded-full bg-[#ff9f43]/20 text-[#ff9f43] text-[11px] font-bold">
                OG DISCORD ROLE VERIFIED
              </span>
              <h1 class="text-2xl font-extrabold text-white mt-1">TriGameDev OG Content Management System</h1>
            </div>
            <button
              @click="triggerDiscordSync"
              class="px-4 py-2.5 rounded-xl bg-[#252438] hover:bg-[#312f4a] border border-[#8b5cf6]/40 text-xs font-bold text-white flex items-center gap-2 cursor-pointer"
            >
              <RefreshCw class="w-4 h-4 text-[#00f5d4]" />
              <span>Sync Discord Events & Announcements</span>
            </button>
          </div>

          <div
            v-if="!currentUser?.is_og_admin"
            class="p-6 rounded-3xl bg-[#242336] border border-[#ff9f43]/40 text-sm text-[#ff9f43]"
          >
            You must be signed in via Discord and hold the <strong>OG Role</strong> in the TriGameDev server to access the CMS.
          </div>

          <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div class="bg-[#242336] rounded-3xl p-6 border border-white/[0.07] flex flex-col gap-4">
              <h2 class="text-lg font-bold text-white">Publish New Community Entry</h2>

              <!-- One-Click Itch.io Scraper Box -->
              <div class="p-4 rounded-2xl bg-[#181726] border border-[#8b5cf6]/40 flex flex-col gap-2.5">
                <div class="flex items-center justify-between">
                  <label class="text-xs font-bold text-[#a78bfa] flex items-center gap-1.5">
                    <Sparkles class="w-3.5 h-3.5 text-[#ff4d8d]" />
                    <span>One-Click Itch.io Jam / Game Importer</span>
                  </label>
                  <span class="text-[11px] text-[#8e8ea9]">Paste a link & publish directly or auto-fill</span>
                </div>
                <div class="flex flex-col sm:flex-row gap-2">
                  <input
                    v-model="cmsForm.url"
                    type="url"
                    placeholder="https://itch.io/jam/... or https://creator.itch.io/game"
                    class="flex-1 px-3.5 py-2.5 rounded-xl bg-[#12121b] border border-white/10 text-sm text-white placeholder-[#6c6c8a]"
                  />
                  <button
                    type="button"
                    @click="autofillFromItch"
                    :disabled="scrapingItch"
                    class="px-4 py-2.5 rounded-xl bg-[#6c5ce7] hover:bg-[#5a4ad1] disabled:opacity-50 text-white text-xs font-bold whitespace-nowrap cursor-pointer transition"
                  >
                    {{ scrapingItch ? 'Scraping Itch.io...' : 'Auto-Fill from Itch.io' }}
                  </button>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="text-xs text-[#8e8ea9]">Target Section</label>
                  <select
                    v-model="cmsForm.section"
                    class="w-full mt-1 px-3 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
                  >
                    <option value="jams">Monthly Game Jam</option>
                    <option value="showcase">Game Showcase</option>
                    <option value="news">News / Announcement</option>
                    <option value="events">Event / Meetup</option>
                  </select>
                </div>
                <div>
                  <label class="text-xs text-[#8e8ea9]">Tag / Engine / Month</label>
                  <input
                    v-model="cmsForm.tag"
                    type="text"
                    placeholder="Auto-filled for Itch links"
                    class="w-full mt-1 px-3 py-2 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
                  />
                </div>
              </div>

              <div>
                <label class="text-xs text-[#8e8ea9]">Title (Optional if Itch.io URL is provided above)</label>
                <input
                  v-model="cmsForm.title"
                  type="text"
                  placeholder="Leave blank to auto-scrape on publish, or click Auto-Fill to preview"
                  class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
                />
              </div>

              <div>
                <label class="text-xs text-[#8e8ea9]">Subtitle / Location / Theme / Developer</label>
                <input
                  v-model="cmsForm.subtitle"
                  type="text"
                  class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
                />
              </div>

              <div>
                <label class="text-xs text-[#8e8ea9]">Cover Image URL</label>
                <input
                  v-model="cmsForm.image_url"
                  type="text"
                  class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
                />
              </div>

              <div>
                <label class="text-xs text-[#8e8ea9]">Body / Description</label>
                <textarea
                  v-model="cmsForm.description"
                  rows="3"
                  class="w-full mt-1 px-3.5 py-2.5 rounded-xl bg-[#181726] border border-white/10 text-sm text-white"
                ></textarea>
              </div>

              <button
                @click="publishCMSItem"
                class="py-3 rounded-xl bg-gradient-to-r from-[#ff9f43] to-[#ff4d8d] text-[#14141d] font-extrabold text-sm cursor-pointer shadow-lg"
              >
                Publish Entry
              </button>
            </div>

            <div class="bg-[#242336] rounded-3xl p-6 border border-white/[0.07] flex flex-col gap-4">
              <h2 class="text-lg font-bold text-white">Manage Live Entries</h2>
              <div class="flex flex-col gap-2.5 max-h-[420px] overflow-y-auto pr-1">
                <div
                  v-for="item in overview.news"
                  :key="'n-' + item.id"
                  class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
                >
                  <div class="truncate">
                    <span class="text-[#ff9f43] font-bold mr-2">[NEWS]</span>
                    <span class="text-white font-semibold">{{ item.title }}</span>
                  </div>
                  <button @click="deleteCMSItem('news', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
                    Delete
                  </button>
                </div>

                <div
                  v-for="item in overview.events"
                  :key="'e-' + item.id"
                  class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
                >
                  <div class="truncate">
                    <span class="text-[#00f5d4] font-bold mr-2">[EVENT]</span>
                    <span class="text-white font-semibold">{{ item.title }}</span>
                  </div>
                  <button @click="deleteCMSItem('events', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
                    Delete
                  </button>
                </div>

                <div
                  v-for="item in overview.jams"
                  :key="'j-' + item.id"
                  class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
                >
                  <div class="truncate">
                    <span class="text-[#ff4d8d] font-bold mr-2">[JAM]</span>
                    <span class="text-white font-semibold">{{ item.title }}</span>
                  </div>
                  <button @click="deleteCMSItem('jams', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
                    Delete
                  </button>
                </div>

                <div
                  v-for="item in overview.showcase"
                  :key="'s-' + item.id"
                  class="p-3 rounded-2xl bg-[#181726] flex items-center justify-between gap-2 text-xs"
                >
                  <div class="truncate">
                    <span class="text-[#a78bfa] font-bold mr-2">[SHOWCASE]</span>
                    <span class="text-white font-semibold">{{ item.title }}</span>
                  </div>
                  <button @click="deleteCMSItem('showcase', item.id)" class="text-[#ff4d8d] hover:underline shrink-0 cursor-pointer">
                    Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>

      <!-- 3. RIGHT CONTEXT SIDEBAR -->
      <aside
        class="bg-gradient-to-b from-[#292643] via-[#211f35] to-[#1d1c2d] border-t lg:border-t-0 lg:border-l border-white/[0.07] p-6 flex flex-col justify-between gap-6"
      >
        <div class="flex flex-col gap-6">
          <!-- Top User Profile & Discord Status Header -->
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3 min-w-0">
              <div v-if="currentUser" class="min-w-0">
                <div class="text-xs font-bold text-white truncate">{{ currentUser.display_name }}</div>
                <div class="text-[11px] font-semibold flex items-center gap-1">
                  <span
                    class="w-2 h-2 rounded-full"
                    :class="currentUser.is_og_admin ? 'bg-[#ff9f43]' : 'bg-[#00f5d4]'"
                  ></span>
                  <span :class="currentUser.is_og_admin ? 'text-[#ff9f43]' : 'text-[#00f5d4]'">
                    {{ currentUser.is_og_admin ? 'OG Role Admin' : 'TriGameDev Member' }}
                  </span>
                </div>
              </div>
              <div v-else>
                <div class="text-xs font-bold text-white">Guest Visitor</div>
                <div class="text-[11px] text-[#8e8ea9]">Sign in via Discord</div>
              </div>
            </div>

            <div v-if="currentUser" class="flex items-center gap-2">
              <img
                :src="currentUser.avatar_url"
                alt="Avatar"
                class="w-10 h-10 rounded-2xl border-2 border-[#8b5cf6] bg-[#14141d] object-cover"
              />
              <button
                @click="logout"
                class="p-2 rounded-xl bg-white/[0.06] hover:bg-white/10 text-[#9e9eb8] hover:text-white cursor-pointer"
                title="Sign Out"
              >
                <LogOut class="w-4 h-4" />
              </button>
            </div>
            <a
              v-else
              href="/api/go/auth/discord/login"
              class="px-3.5 py-2 rounded-xl bg-[#5865F2] hover:bg-[#4752c4] text-white text-xs font-bold transition"
            >
              Discord Login
            </a>
          </div>

          <!-- Featured Event / Jam Spotlight Card -->
          <div class="flex flex-col gap-3">
            <div class="flex items-center justify-between">
              <h3 class="font-bold text-sm text-white">Featured Spotlight</h3>
              <button
                @click="activeSection = 'events'"
                class="text-xs text-[#c4c4dc] hover:text-white underline cursor-pointer"
              >
                Events
              </button>
            </div>

            <div class="bg-[#1b1a29] rounded-3xl p-4 border border-white/[0.08] flex flex-col gap-4 shadow-xl">
              <div class="relative h-44 rounded-2xl overflow-hidden">
                <img
                  :src="nextEvent?.cover_image_url || '/inspo/rooftop-twilight.jpg'"
                  alt="Featured Event"
                  class="w-full h-full object-cover"
                />
                <span
                  class="absolute bottom-2.5 right-2.5 px-2.5 py-1 rounded-lg bg-white text-[#14141d] text-[10px] font-extrabold shadow"
                >
                  TriGameDev ✓
                </span>
              </div>

              <div
                class="p-3 rounded-2xl bg-[#27263b] border border-white/10 flex items-center justify-between gap-2"
              >
                <div class="min-w-0">
                  <div class="font-bold text-xs text-white truncate">
                    {{ nextEvent?.title || 'TriGameDev Community Hub' }}
                  </div>
                  <div class="text-[10px] text-[#9e9eb8] truncate">
                    {{ nextEvent?.location || 'Join us on Discord' }}
                  </div>
                </div>
                <span
                  class="w-8 h-8 rounded-full bg-[#ff2e7e] flex items-center justify-center text-white shrink-0 shadow-md shadow-[#ff2e7e]/40"
                >
                  <Radio class="w-4 h-4" />
                </span>
              </div>

              <div class="flex flex-col gap-2 text-xs px-1">
                <div class="flex items-center justify-between">
                  <span class="text-[#8e8ea9]">Available Vault Keys:</span>
                  <span class="font-bold text-[#00f5d4]">{{ overview.stats?.available_keys || 0 }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-[#8e8ea9]">Library Ebooks:</span>
                  <span class="font-bold text-[#00f5d4]">{{ overview.stats?.total_ebooks || 0 }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-[#8e8ea9]">Showcase Games:</span>
                  <span class="font-bold text-[#00f5d4]">{{ overview.stats?.total_showcase || 0 }}</span>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-2.5 pt-1">
                <button
                  @click="activeSection = 'keys'"
                  class="py-2 rounded-xl bg-gradient-to-r from-[#6c5ce7] to-[#8b5cf6] text-white text-xs font-bold cursor-pointer shadow-md shadow-[#6c5ce7]/25"
                >
                  Key Vault
                </button>
                <button
                  @click="activeSection = 'library'"
                  class="py-2 rounded-xl border border-white/20 hover:border-white/40 text-white text-xs font-bold cursor-pointer"
                >
                  Ebook Library
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Bottom Pulse Decorative Graphic -->
        <div class="flex flex-col gap-2 pt-2 border-t border-white/[0.06]">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-white">TriGameDev Signal</span>
            <span class="text-xs font-bold text-[#00f5d4]">ONLINE</span>
          </div>
          <div class="relative pt-2 pb-1">
            <svg viewBox="0 0 260 70" class="w-full h-14 overflow-visible">
              <defs>
                <linearGradient id="pulseGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#8b5cf6" stop-opacity="0.45" />
                  <stop offset="100%" stop-color="#8b5cf6" stop-opacity="0.0" />
                </linearGradient>
              </defs>
              <path
                d="M0,55 C35,50 55,20 90,38 C125,55 145,15 180,22 C215,28 235,8 260,12 L260,70 L0,70 Z"
                fill="url(#pulseGrad)"
              />
              <path
                d="M0,55 C35,50 55,20 90,38 C125,55 145,15 180,22 C215,28 235,8 260,12"
                fill="none"
                stroke="#8b5cf6"
                stroke-width="2.5"
              />
              <circle cx="90" cy="38" r="4" fill="#00f5d4" />
              <circle cx="180" cy="22" r="4.5" fill="#00f5d4" stroke="#ffffff" stroke-width="1.5" />
              <circle cx="235" cy="11" r="4" fill="#00f5d4" />
            </svg>
          </div>
        </div>
      </aside>
    </div>

    <!-- MODAL: Revealed Key Code upon Redemption -->
    <div
      v-if="revealedKeyModal"
      class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
    >
      <div class="bg-[#232234] border border-[#00f5d4]/40 rounded-3xl max-w-md w-full p-6 flex flex-col gap-4 shadow-2xl">
        <div class="flex items-center justify-between">
          <span class="px-2.5 py-1 rounded-full bg-[#00f5d4]/15 text-[#00f5d4] text-xs font-bold">
            KEY CLAIMED & BOUND TO YOUR ACCOUNT
          </span>
          <button @click="revealedKeyModal = null" class="text-[#8e8ea9] hover:text-white text-sm cursor-pointer">✕</button>
        </div>
        <h3 class="text-xl font-extrabold text-white">{{ revealedKeyModal.title }}</h3>
        <p class="text-xs text-[#9e9eb8]">
          This key has been removed from the community pool and permanently saved under your
          <strong class="text-white">My Keys</strong> tab.
        </p>
        <div class="p-4 rounded-2xl bg-[#14141d] border border-white/10 flex items-center justify-between gap-2">
          <code class="text-sm font-mono text-[#00f5d4] font-bold">{{ revealedKeyModal.key_code }}</code>
          <button
            @click="copyCode(revealedKeyModal.key_code, 'modal')"
            class="px-3 py-1.5 rounded-xl bg-[#6c5ce7] text-white text-xs font-bold cursor-pointer"
          >
            Copy
          </button>
        </div>
        <button
          @click="revealedKeyModal = null; keySubTab = 'my-keys'"
          class="w-full py-2.5 rounded-xl bg-white text-[#14141d] font-bold text-xs cursor-pointer"
        >
          Go to My Keys Vault
        </button>
      </div>
    </div>

    <!-- MODAL: Donate Pool of Keys -->
    <div
      v-if="showDonateModal"
      class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
    >
      <div class="bg-[#232234] border border-white/15 rounded-3xl max-w-lg w-full p-6 flex flex-col gap-4 shadow-2xl">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-extrabold text-white">Donate Game Keys to the Pool</h3>
          <button @click="showDonateModal = false" class="text-[#8e8ea9] hover:text-white cursor-pointer">✕</button>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div class="sm:col-span-2">
            <label class="text-xs text-[#8e8ea9]">Game / Tool Title</label>
            <input
              v-model="donateForm.title"
              type="text"
              class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
            />
          </div>
          <div>
            <label class="text-xs text-[#8e8ea9]">Platform</label>
            <select
              v-model="donateForm.platform"
              class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
            >
              <option>Steam</option>
              <option>Itch.io</option>
              <option>GOG</option>
              <option>Epic</option>
            </select>
          </div>
        </div>
        <div>
          <label class="text-xs text-[#8e8ea9]">Short Description</label>
          <input
            v-model="donateForm.description"
            type="text"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          />
        </div>
        <div>
          <label class="text-xs text-[#8e8ea9]">Paste Key Codes (One key per line)</label>
          <textarea
            v-model="donateForm.raw_keys"
            rows="4"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-xs font-mono text-[#00f5d4]"
          ></textarea>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button @click="showDonateModal = false" class="px-4 py-2 rounded-xl text-xs text-[#8e8ea9] cursor-pointer">
            Cancel
          </button>
          <button
            @click="submitKeyDonation"
            class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-[#ff4d8d] to-[#ff7eb3] text-white text-xs font-bold cursor-pointer"
          >
            Add Keys to Pool
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Upload Ebook -->
    <div
      v-if="showUploadEbookModal"
      class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
    >
      <div class="bg-[#232234] border border-white/15 rounded-3xl max-w-lg w-full p-6 flex flex-col gap-4 shadow-2xl">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-extrabold text-white">Add Ebook to TriGameDev Library</h3>
          <button @click="showUploadEbookModal = false" class="text-[#8e8ea9] hover:text-white cursor-pointer">✕</button>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="text-xs text-[#8e8ea9]">Book Title</label>
            <input
              v-model="ebookForm.title"
              type="text"
              class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
            />
          </div>
          <div>
            <label class="text-xs text-[#8e8ea9]">Author</label>
            <input
              v-model="ebookForm.author"
              type="text"
              class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
            />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="text-xs text-[#8e8ea9]">Category</label>
            <input
              v-model="ebookForm.category"
              type="text"
              class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
            />
          </div>
          <div>
            <label class="text-xs text-[#8e8ea9]">Format</label>
            <select
              v-model="ebookForm.file_format"
              class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
            >
              <option>PDF</option>
              <option>EPUB</option>
              <option>ZIP</option>
            </select>
          </div>
        </div>
        <div>
          <label class="text-xs text-[#8e8ea9]">Synopsis / Notes</label>
          <textarea
            v-model="ebookForm.description"
            rows="3"
            class="w-full mt-1 px-3 py-2 rounded-xl bg-[#14141d] border border-white/10 text-sm text-white"
          ></textarea>
        </div>
        <div class="flex justify-end gap-2">
          <button @click="showUploadEbookModal = false" class="px-4 py-2 rounded-xl text-xs text-[#8e8ea9] cursor-pointer">
            Cancel
          </button>
          <button
            @click="submitEbookUpload"
            class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-[#ff9f43] to-[#ff4d8d] text-[#14141d] text-xs font-extrabold cursor-pointer"
          >
            Publish to Library
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Gallery Lightbox -->
    <div
      v-if="selectedLightbox"
      @click.self="selectedLightbox = null"
      class="fixed inset-0 z-50 bg-black/85 backdrop-blur-md flex items-center justify-center p-4"
    >
      <div class="bg-[#232234] border border-white/15 rounded-3xl max-w-2xl w-full overflow-hidden shadow-2xl">
        <img :src="selectedLightbox.image_url" class="w-full max-h-[65vh] object-cover" />
        <div class="p-5 flex items-center justify-between">
          <div>
            <span class="text-xs font-bold text-[#ff4d8d]">{{ selectedLightbox.event_label }}</span>
            <h3 class="text-lg font-bold text-white">{{ selectedLightbox.title }}</h3>
            <p class="text-xs text-[#9e9eb8] mt-0.5">{{ selectedLightbox.caption }}</p>
          </div>
          <button
            @click="selectedLightbox = null"
            class="px-4 py-2 rounded-xl bg-white/10 text-white text-xs font-bold cursor-pointer"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
