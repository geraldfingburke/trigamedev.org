<script setup>
import { ref, computed } from 'vue'
import { RefreshCw, Terminal, Upload, Trophy, ChevronDown, ChevronRight } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'
import { useCms } from '../../composables/useCms'

const { currentUser, overview } = useTriGameDevStore()
const {
  scrapingItch,
  uploadingCmsImage,
  cmsForm,
  autofillFromItch,
  handleCmsImageUpload,
  handleVideoUrlInput,
  publishCMSItem,
  deleteCMSItem,
  triggerDiscordSync
} = useCms()

const onlyMyPosts = ref(false)
const openPostCategories = ref({})

function isAuthoredByUser(item) {
  const u = currentUser.value
  if (!u) return false
  if (item.created_by && item.created_by === u.id) return true
  const disp = (u.display_name || '').toLowerCase()
  const uname = (u.username || '').toLowerCase()
  const fields = [item.author_name, item.uploaded_by, item.developer_name, item.speaker_name]
  return fields.some(val => {
    if (!val) return false
    const low = val.toLowerCase()
    return (disp && low.includes(disp)) || (uname && low.includes(uname))
  })
}

const postCategories = computed(() => {
  const raw = [
    { key: 'jams', label: 'Jams', items: overview.value.jams || [] },
    { key: 'showcase', label: 'Games', items: overview.value.showcase || [] },
    { key: 'events', label: 'Events', items: overview.value.events || [] },
    { key: 'news', label: 'Announcements', items: overview.value.news || [] },
    { key: 'gallery', label: 'Photos', items: overview.value.gallery || [] },
    { key: 'videos', label: 'Videos', items: overview.value.videos || [] }
  ]
  return raw.map(cat => ({
    ...cat,
    items: onlyMyPosts.value ? cat.items.filter(isAuthoredByUser) : cat.items
  }))
})

function isCategoryOpen(key) {
  return !!openPostCategories.value[key]
}

function toggleCategory(key) {
  openPostCategories.value = {
    ...openPostCategories.value,
    [key]: !openPostCategories.value[key]
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
      <div>
        <h1 class="text-2xl font-extrabold text-white uppercase">Content Management System</h1>
      </div>
      <button
        @click="triggerDiscordSync"
        class="px-4 py-2.5 bg-[#121215] hover:bg-zinc-800 border border-zinc-700 text-xs font-bold uppercase tracking-wider text-white flex items-center gap-2 cursor-pointer transition"
      >
        <RefreshCw class="w-4 h-4 text-[#ef4444]" />
        <span>Sync Discord Events</span>
      </button>
    </div>

    <div
      v-if="!currentUser?.is_og_admin"
      class="p-6 bg-[#121215] border border-[#ef4444] text-sm text-zinc-200"
    >
      You must be signed in via Discord and hold the <strong class="text-white">OG</strong> role in the TriGameDev server to access the CMS.
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Dynamic Form based on Post Type -->
      <div class="bg-[#121215] p-6 border border-zinc-800 flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
          <h2 class="text-base font-bold text-white uppercase tracking-wider">Create Post</h2>
          <div class="flex items-center gap-2">
            <label class="text-xs font-mono uppercase text-zinc-400">Post Type:</label>
            <select
              v-model="cmsForm.section"
              class="px-3 py-1.5 bg-[#09090b] border border-zinc-700 text-xs font-mono uppercase text-white focus:outline-none focus:border-white"
            >
              <option value="jams">Game Jam</option>
              <option value="showcase">Game (By Our People)</option>
              <option value="news">News / Announcement</option>
              <option value="events">Event</option>
              <option value="gallery">Photo (Gallery)</option>
              <option value="videos">Video Link</option>
            </select>
          </div>
        </div>

        <!-- Game Importer (Jams and Games) -->
        <div
          v-if="cmsForm.section === 'jams' || cmsForm.section === 'showcase'"
          class="p-4 bg-[#09090b] border border-zinc-700 flex flex-col gap-2.5"
        >
          <label class="text-xs font-mono uppercase text-white flex items-center gap-1.5">
            <Terminal class="w-3.5 h-3.5 text-[#ef4444]" />
            <span>Game Importer</span>
          </label>
          <div class="flex flex-col sm:flex-row gap-2">
            <input
              v-model="cmsForm.url"
              type="url"
              placeholder="https://..."
              class="flex-1 px-3.5 py-2 bg-[#121215] border border-zinc-700 text-sm text-white placeholder-zinc-600 focus:outline-none focus:border-white"
            />
            <button
              type="button"
              @click="autofillFromItch"
              :disabled="scrapingItch"
              class="px-4 py-2 bg-white hover:bg-zinc-200 disabled:opacity-50 text-black text-xs font-bold uppercase tracking-wider whitespace-nowrap cursor-pointer transition"
            >
              {{ scrapingItch ? 'Importing...' : 'Import' }}
            </button>
          </div>
        </div>

        <!-- DYNAMIC FIELDS: JAMS -->
        <template v-if="cmsForm.section === 'jams'">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Jam Category</label>
              <select
                v-model="cmsForm.jam_type"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              >
                <option value="monthly">Monthly Game Jam</option>
                <option value="other">Other Jams</option>
              </select>
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Month / Label</label>
              <input
                v-model="cmsForm.tag"
                type="text"
                placeholder="e.g. April 2025"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Jam Title</label>
            <input
              v-model="cmsForm.title"
              type="text"
              placeholder="TriGameDev Monthly Game Jam #1"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Jam Theme</label>
            <input
              v-model="cmsForm.subtitle"
              type="text"
              placeholder="Next line after 'The Theme'"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
          </div>

          <div v-if="cmsForm.winners && cmsForm.winners.length > 0" class="p-3 bg-[#09090b] border border-zinc-800 flex flex-col gap-1.5">
            <div class="text-[11px] font-mono uppercase text-[#ef4444] flex items-center gap-1.5">
              <Trophy class="w-3.5 h-3.5" />
              <span>Scraped Winners Circle ({{ cmsForm.winners.length }})</span>
            </div>
            <div v-for="w in cmsForm.winners" :key="w.rank" class="text-xs font-mono text-zinc-300">
              {{ w.medal.toUpperCase() }} (#{{ w.rank }}): <strong class="text-white">{{ w.title }}</strong> by {{ w.authors }}
            </div>
          </div>
        </template>

        <!-- DYNAMIC FIELDS: SHOWCASE (GAMES BY OUR PEOPLE) -->
        <template v-else-if="cmsForm.section === 'showcase'">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Game Title</label>
              <input
                v-model="cmsForm.title"
                type="text"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Engine / Tools</label>
              <input
                v-model="cmsForm.tag"
                type="text"
                placeholder="Godot, Unity, Unreal..."
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Developer(s)</label>
              <input
                v-model="cmsForm.subtitle"
                type="text"
                class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Button Label</label>
              <select
                v-model="cmsForm.action_label"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              >
                <option value="Buy on Steam">Buy on Steam</option>
                <option value="Wishlist on Steam">Wishlist on Steam</option>
                <option value="Play Now on Itch.io">Play Now on Itch.io</option>
                <option value="Download from Itch.io">Download from Itch.io</option>
              </select>
            </div>
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Short Description</label>
            <textarea
              v-model="cmsForm.description"
              rows="2"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            ></textarea>
          </div>
        </template>

        <!-- DYNAMIC FIELDS: NEWS / ANNOUNCEMENT -->
        <template v-else-if="cmsForm.section === 'news'">
          <div class="grid grid-cols-3 gap-3">
            <div class="col-span-2">
              <label class="text-xs font-mono uppercase text-zinc-400">Headline</label>
              <input
                v-model="cmsForm.title"
                type="text"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Tag</label>
              <input
                v-model="cmsForm.tag"
                type="text"
                placeholder="Announcement"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Short Summary</label>
            <input
              v-model="cmsForm.subtitle"
              type="text"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Announcement Body</label>
            <textarea
              v-model="cmsForm.description"
              rows="4"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            ></textarea>
          </div>
        </template>

        <!-- DYNAMIC FIELDS: EVENTS -->
        <template v-else-if="cmsForm.section === 'events'">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Event Title</label>
              <input
                v-model="cmsForm.title"
                type="text"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Event Date &amp; Time</label>
              <input
                v-model="cmsForm.start_date"
                type="datetime-local"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Location</label>
              <input
                v-model="cmsForm.subtitle"
                type="text"
                placeholder="TriGameDev Discord / Willow Tree"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Event Type</label>
              <input
                v-model="cmsForm.tag"
                type="text"
                placeholder="Meetup"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Event Details</label>
            <textarea
              v-model="cmsForm.description"
              rows="3"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            ></textarea>
          </div>
        </template>

        <!-- DYNAMIC FIELDS: GALLERY PHOTO -->
        <template v-else-if="cmsForm.section === 'gallery'">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Photo Title</label>
              <input
                v-model="cmsForm.title"
                type="text"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Event / Album Label</label>
              <input
                v-model="cmsForm.tag"
                type="text"
                placeholder="Community Meetup"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Caption</label>
            <input
              v-model="cmsForm.description"
              type="text"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
          </div>
        </template>

        <!-- DYNAMIC FIELDS: VIDEO LINK -->
        <template v-else-if="cmsForm.section === 'videos'">
          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Video Link (YouTube, Twitch, Vimeo...)</label>
            <input
              v-model="cmsForm.url"
              @input="handleVideoUrlInput"
              type="url"
              placeholder="https://www.youtube.com/watch?v=..."
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
          </div>

          <div class="grid grid-cols-3 gap-3">
            <div class="col-span-2">
              <label class="text-xs font-mono uppercase text-zinc-400">Video Title</label>
              <input
                v-model="cmsForm.title"
                type="text"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
            <div>
              <label class="text-xs font-mono uppercase text-zinc-400">Duration</label>
              <input
                v-model="cmsForm.tag"
                type="text"
                placeholder="18:40"
                class="w-full mt-1 px-3 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
              />
            </div>
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Speaker / Creator</label>
            <input
              v-model="cmsForm.subtitle"
              type="text"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            />
          </div>

          <div>
            <label class="text-xs font-mono uppercase text-zinc-400">Description</label>
            <textarea
              v-model="cmsForm.description"
              rows="2"
              class="w-full mt-1 px-3.5 py-2 bg-[#09090b] border border-zinc-700 text-sm text-white focus:outline-none focus:border-white"
            ></textarea>
          </div>
        </template>

        <!-- Cover / Photo Image Upload or URL (For all types except News) -->
        <div v-if="cmsForm.section !== 'news'" class="p-3.5 bg-[#09090b] border border-zinc-800 flex flex-col gap-2.5">
          <div class="flex items-center justify-between">
            <label class="text-xs font-mono uppercase text-zinc-300">
              {{ cmsForm.section === 'gallery' ? 'Photo File or URL' : 'Cover Image (Upload or URL)' }}
            </label>
            <label class="px-3 py-1 bg-white hover:bg-zinc-200 text-black text-[11px] font-bold uppercase tracking-wider flex items-center gap-1.5 cursor-pointer">
              <Upload class="w-3 h-3" />
              <span>{{ uploadingCmsImage ? 'Uploading...' : 'Upload Image' }}</span>
              <input type="file" accept="image/*" @change="handleCmsImageUpload" class="hidden" />
            </label>
          </div>

          <div class="flex items-center gap-3">
            <img
              v-if="cmsForm.image_url"
              :src="cmsForm.image_url"
              class="w-16 h-12 object-cover border border-zinc-700 bg-black shrink-0"
            />
            <input
              v-model="cmsForm.image_url"
              type="text"
              placeholder="Or paste an image URL..."
              class="flex-1 px-3 py-2 bg-[#121215] border border-zinc-700 text-xs text-white focus:outline-none focus:border-white"
            />
          </div>
        </div>

        <button
          @click="publishCMSItem"
          class="py-3 bg-[#ef4444] hover:bg-red-500 text-white font-extrabold text-xs uppercase tracking-wider cursor-pointer transition"
        >
          Publish {{ cmsForm.section.toUpperCase() }}
        </button>
      </div>

      <!-- Posts (Category Accordions, Default Collapsed) -->
      <div class="bg-[#121215] p-6 border border-zinc-800 flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
          <h2 class="text-base font-bold text-white uppercase tracking-wider">
            Posts
          </h2>
          <label class="flex items-center gap-2 text-xs font-mono uppercase text-zinc-400 hover:text-white cursor-pointer select-none">
            <input
              v-model="onlyMyPosts"
              type="checkbox"
              class="accent-[#ef4444] w-3.5 h-3.5 cursor-pointer"
            />
            <span>Authored by me</span>
          </label>
        </div>

        <div class="flex flex-col gap-2.5">
          <div
            v-for="cat in postCategories"
            :key="cat.key"
            class="bg-[#09090b] border border-zinc-800"
          >
            <button
              type="button"
              @click="toggleCategory(cat.key)"
              class="w-full px-4 py-3 flex items-center justify-between gap-3 text-left hover:bg-zinc-900/70 transition cursor-pointer"
              :class="isCategoryOpen(cat.key) ? 'border-b border-zinc-800 bg-zinc-900/40' : ''"
            >
              <div class="flex items-center gap-2.5">
                <ChevronDown
                  v-if="isCategoryOpen(cat.key)"
                  class="w-4 h-4 text-[#ef4444] shrink-0"
                />
                <ChevronRight
                  v-else
                  class="w-4 h-4 text-zinc-400 shrink-0"
                />
                <span class="font-bold text-sm text-white uppercase tracking-wider">
                  {{ cat.label }}
                </span>
              </div>
              <span class="px-2 py-0.5 bg-[#121215] border border-zinc-700 text-[11px] font-mono text-zinc-300">
                {{ cat.items.length }}
              </span>
            </button>

            <div
              v-if="isCategoryOpen(cat.key)"
              class="p-3 flex flex-col gap-2 max-h-72 overflow-y-auto"
            >
              <div
                v-if="cat.items.length === 0"
                class="py-3 text-center text-xs font-mono uppercase text-zinc-500"
              >
                No posts in this category.
              </div>
              <div
                v-for="item in cat.items"
                :key="cat.key + '-' + item.id"
                class="p-2.5 bg-[#121215] border border-zinc-800 flex items-center justify-between gap-2 text-xs"
              >
                <span class="text-white font-semibold truncate">{{ item.title }}</span>
                <button
                  @click="deleteCMSItem(cat.key, item.id)"
                  class="text-[#ef4444] font-mono uppercase hover:underline shrink-0 cursor-pointer"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
