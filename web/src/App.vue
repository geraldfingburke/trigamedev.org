<script setup>
import { computed, onMounted } from 'vue'
import { useTriGameDevStore } from './composables/useTriGameDevStore'

// Layout Shell Components
import NavRail from './components/layout/NavRail.vue'
import TopHeader from './components/layout/TopHeader.vue'
import RightSidebar from './components/layout/RightSidebar.vue'

// View Components
import ExploreView from './components/views/ExploreView.vue'
import EventsView from './components/views/EventsView.vue'
import NewsView from './components/views/NewsView.vue'
import JamsView from './components/views/JamsView.vue'
import ShowcaseView from './components/views/ShowcaseView.vue'
import MediaView from './components/views/MediaView.vue'
import KeyVaultView from './components/views/KeyVaultView.vue'
import EbookLibraryView from './components/views/EbookLibraryView.vue'
import CmsView from './components/views/CmsView.vue'

// Global Modals & Notifications
import ToastNotice from './components/modals/ToastNotice.vue'
import RevealedKeyModal from './components/modals/RevealedKeyModal.vue'
import DonateKeysModal from './components/modals/DonateKeysModal.vue'
import UploadEbookModal from './components/modals/UploadEbookModal.vue'
import GalleryLightboxModal from './components/modals/GalleryLightboxModal.vue'
import WinnersCircleModal from './components/modals/WinnersCircleModal.vue'

const { activeSection, showToast, loadAllData } = useTriGameDevStore()

const views = {
  explore: ExploreView,
  events: EventsView,
  news: NewsView,
  jams: JamsView,
  showcase: ShowcaseView,
  media: MediaView,
  keys: KeyVaultView,
  library: EbookLibraryView,
  cms: CmsView
}

const currentViewComponent = computed(() => views[activeSection.value] || ExploreView)

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
  <div class="min-h-screen bg-[#09090b] text-[#f4f4f5] flex flex-col">
    <!-- Toast Notification -->
    <ToastNotice />

    <!-- Main Brutalist Grid Frame (Sharp edges, 1px structural borders, no gradients) -->
    <div
      class="flex-1 w-full max-w-[1680px] mx-auto bg-[#0f0f12] border-x border-zinc-800 grid grid-cols-1 lg:grid-cols-[72px_1fr_340px]"
    >
      <!-- 1. Left Navigation Rail -->
      <NavRail />

      <!-- 2. Center Main Stage -->
      <main class="p-6 sm:p-8 flex flex-col gap-8 overflow-y-auto">
        <TopHeader />
        <component :is="currentViewComponent" />
      </main>

      <!-- 3. Right Context Sidebar -->
      <RightSidebar />
    </div>

    <!-- Global Modals -->
    <RevealedKeyModal />
    <DonateKeysModal />
    <UploadEbookModal />
    <GalleryLightboxModal />
    <WinnersCircleModal />
  </div>
</template>
