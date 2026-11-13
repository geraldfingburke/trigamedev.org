<script setup>
import {
  Radio,
  KeyRound,
  ChevronRight,
  Eye,
  Newspaper,
  BookOpen,
  ExternalLink
} from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const {
  activeSection,
  overview,
  currentUser,
  activeJam,
  filteredShowcase,
  selectedLightbox
} = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-8">
    <!-- Brutalist Brandkit Hero Banner -->
    <section
      class="relative overflow-hidden border border-zinc-800 bg-black min-h-[260px] flex flex-col lg:flex-row items-stretch justify-between"
    >
      <img
        src="/brand/filler.png"
        alt="TriGameDev Space Invaders"
        class="absolute inset-0 w-full h-full object-cover object-center opacity-25 pointer-events-none"
      />

      <div class="relative z-10 p-6 sm:p-10 flex flex-col justify-between gap-6 max-w-2xl">
        <div class="flex flex-col gap-3">
          <div class="inline-flex items-center gap-2 px-3 py-1 bg-[#09090b] border border-zinc-700 w-fit text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">
            <Radio class="w-3.5 h-3.5 text-[#ef4444]" />
            <span>{{ activeJam ? `${activeJam.month_label} // LIVE JAM` : 'NORTHEAST TENNESSEE & BEYOND' }}</span>
          </div>

          <img
            src="/brand/logo-large.png"
            alt="TRIGAMEDEV"
            class="h-14 sm:h-20 w-auto object-contain object-left"
          />

          <p class="text-sm sm:text-base text-zinc-300 leading-relaxed max-w-xl">
            <template v-if="activeJam">
              Active Jam: <strong class="text-white">{{ activeJam.title }}</strong> — Theme:
              <span class="font-mono font-semibold text-[#ef4444]">{{ activeJam.theme }}</span>.
              {{ activeJam.description }}
            </template>
            <template v-else>
              A loose collective of game developers that started in the Tri-Cities Region of Northeast Tennessee.
              Solo game dev can be lonely — let's be alone together.
            </template>
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <button
            @click="activeSection = 'jams'"
            class="px-6 py-2.5 bg-white text-black font-bold text-xs uppercase tracking-wider hover:bg-zinc-200 transition cursor-pointer"
          >
            Monthly Jams
          </button>
          <button
            @click="activeSection = 'keys'"
            class="px-5 py-2.5 bg-[#121215] hover:bg-zinc-800 text-white border border-zinc-700 font-bold text-xs uppercase tracking-wider flex items-center gap-2 transition cursor-pointer"
          >
            <KeyRound class="w-3.5 h-3.5 text-[#ef4444]" />
            <span>Key Vault ({{ overview.stats?.available_keys || 0 }} Available)</span>
          </button>
          <a
            href="https://discord.gg/6jdhMV6sqt"
            target="_blank"
            class="px-5 py-2.5 bg-transparent hover:bg-zinc-900 text-zinc-300 border border-zinc-800 font-bold text-xs uppercase tracking-wider flex items-center gap-1.5 transition"
          >
            <span>Join Discord</span>
            <ExternalLink class="w-3.5 h-3.5" />
          </a>
        </div>
      </div>

      <!-- Right Hero Brand Badge Panel -->
      <div class="relative z-10 hidden lg:flex flex-col items-center justify-center p-8 border-l border-zinc-800 bg-[#0c0c0f]/90 w-80 shrink-0 gap-4">
        <img
          src="/brand/circle-logo-large.png"
          alt="TriGameDev Badge"
          class="w-36 h-36 object-contain border border-zinc-800"
        />
        <img
          src="/brand/tres.png"
          alt="Una Substantia, Tres Personae"
          class="h-6 object-contain"
        />
      </div>
    </section>

    <!-- What / Who / Why Manifesto Strip (From tgd-brandkit) -->
    <section class="grid grid-cols-1 md:grid-cols-3 border border-zinc-800 bg-[#121215] divide-y md:divide-y-0 md:divide-x divide-zinc-800">
      <div class="p-6 flex flex-col gap-2">
        <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">01 // WHAT IS TRIGAMEDEV?</div>
        <h3 class="text-base font-bold text-white uppercase">A Loose Collective</h3>
        <p class="text-xs text-zinc-400 leading-relaxed">
          TriGameDev started in the Tri-Cities Region of Northeast Tennessee. Over the years, people have moved away and new members have found us online. All are welcome.
        </p>
      </div>

      <div class="p-6 flex flex-col gap-2">
        <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">02 // WHO IS TRIGAMEDEV?</div>
        <h3 class="text-base font-bold text-white uppercase">Pros, Hobbyists & Players</h3>
        <p class="text-xs text-zinc-400 leading-relaxed">
          Industry professionals, hobbyists, students, curious non-developers, and game playing people. Every discipline and experience level has a seat at the table.
        </p>
      </div>

      <div class="p-6 flex flex-col gap-2">
        <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">03 // WHY IS TRIGAMEDEV?</div>
        <h3 class="text-base font-bold text-white uppercase">Let's Be Alone Together</h3>
        <p class="text-xs text-zinc-400 leading-relaxed">
          Game dev is hard. Solo game dev can be lonely. We host monthly 48-hour game jams, share resources, swap keys, and celebrate shipping builds together.
        </p>
      </div>
    </section>

    <!-- 3-Column Featured Showcase Cards -->
    <section class="flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <div class="flex items-center gap-3">
          <h2 class="text-lg font-bold text-white uppercase tracking-wider">Member Showcase</h2>
          <span class="text-[11px] font-mono px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-zinc-300 uppercase">
            Community Builds
          </span>
        </div>
        <button
          @click="activeSection = 'showcase'"
          class="px-3.5 py-1.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider flex items-center gap-1 transition cursor-pointer"
        >
          <span>View All</span>
          <ChevronRight class="w-3.5 h-3.5" />
        </button>
      </div>

      <div
        v-if="filteredShowcase.length === 0"
        class="p-8 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No showcase games published yet. OG Admins can add projects from the CMS.
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <article
          v-for="game in filteredShowcase.slice(0, 3)"
          :key="game.id"
          class="group bg-[#121215] border border-zinc-800 hover:border-zinc-500 transition flex flex-col justify-between"
        >
          <div class="relative h-52 overflow-hidden bg-black border-b border-zinc-800">
            <img
              :src="game.thumbnail_url || '/brand/filler.png'"
              :alt="game.title"
              class="w-full h-full object-cover"
            />
            <div
              class="absolute top-3 right-3 px-2 py-0.5 bg-black border border-zinc-700 text-white text-[11px] font-mono flex items-center gap-1"
            >
              <Eye class="w-3 h-3 text-zinc-400" />
              <span>{{ game.views_count }}</span>
            </div>
            <div
              class="absolute bottom-0 inset-x-0 p-3 bg-black/90 border-t border-zinc-800 flex items-center justify-between gap-2"
            >
              <div class="min-w-0">
                <h3 class="font-bold text-sm text-white truncate uppercase">{{ game.title }}</h3>
                <p class="text-[11px] font-mono text-zinc-400 truncate">BY {{ game.developer_name }}</p>
              </div>
              <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[10px] font-mono text-[#ef4444] shrink-0 uppercase">
                {{ game.engine }}
              </span>
            </div>
          </div>

          <div class="p-4 flex flex-col gap-3">
            <p class="text-xs text-zinc-400 line-clamp-2">{{ game.description }}</p>
            <div class="grid grid-cols-2 gap-2">
              <a
                :href="game.play_url || '#'"
                target="_blank"
                class="py-2 px-3 bg-white hover:bg-zinc-200 text-black font-bold text-xs uppercase tracking-wider text-center transition"
              >
                Play Build
              </a>
              <a
                :href="game.source_url || '#'"
                target="_blank"
                class="py-2 px-3 bg-transparent border border-zinc-700 hover:border-white text-zinc-300 hover:text-white font-bold text-xs uppercase tracking-wider text-center transition"
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
      <div class="flex items-center justify-between border-b border-zinc-800 pb-2">
        <h2 class="text-sm font-bold text-white uppercase tracking-wider">Recent Gallery Snaps</h2>
        <button
          @click="activeSection = 'media'"
          class="text-xs font-mono uppercase text-zinc-400 hover:text-white cursor-pointer"
        >
          Open Gallery & Videos ->
        </button>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-7 gap-3">
        <button
          v-for="snap in overview.gallery.slice(0, 7)"
          :key="snap.id"
          @click="selectedLightbox = snap"
          class="group relative aspect-square overflow-hidden bg-[#121215] border border-zinc-800 hover:border-white transition cursor-pointer"
        >
          <img
            :src="snap.image_url"
            :alt="snap.title"
            class="w-full h-full object-cover"
          />
        </button>
      </div>
    </section>

    <!-- Two-Column Quick Feed: Latest News & Member Portals -->
    <section class="grid grid-cols-1 xl:grid-cols-2 gap-5">
      <!-- Latest News & Discord Announcements -->
      <div class="bg-[#121215] p-5 border border-zinc-800 flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
          <h3 class="font-bold text-sm text-white uppercase tracking-wider flex items-center gap-2">
            <Newspaper class="w-4 h-4 text-[#ef4444]" />
            <span>News & Announcements</span>
          </h3>
          <button
            @click="activeSection = 'news'"
            class="text-xs font-mono uppercase text-zinc-400 hover:text-white cursor-pointer"
          >
            View All ->
          </button>
        </div>

        <div v-if="!overview.news || overview.news.length === 0" class="text-xs font-mono uppercase text-zinc-500 py-6 text-center">
          No announcements posted yet.
        </div>

        <div v-else class="flex flex-col gap-3">
          <div
            v-for="post in overview.news.slice(0, 3)"
            :key="post.id"
            class="p-3.5 bg-[#09090b] border border-zinc-800 flex flex-col gap-1"
          >
            <div class="flex items-center justify-between text-[11px] font-mono">
              <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[#ef4444] uppercase">
                {{ post.tag }}
              </span>
              <span class="text-zinc-500">{{ post.published_at }} // {{ post.author_name }}</span>
            </div>
            <h4 class="font-bold text-sm text-white mt-1">{{ post.title }}</h4>
            <p class="text-xs text-zinc-400 line-clamp-2">{{ post.summary }}</p>
          </div>
        </div>
      </div>

      <!-- Member Vault Quick Access -->
      <div class="bg-[#121215] p-5 border border-zinc-800 flex flex-col justify-between gap-4">
        <div class="flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span class="px-2 py-0.5 text-[10px] font-mono uppercase tracking-widest bg-zinc-900 border border-zinc-700 text-zinc-300">
              DISCORD MEMBER PERKS
            </span>
          </div>
          <h3 class="text-lg font-bold text-white uppercase mt-1">
            Key Sharing Portal & Ebook Library
          </h3>
          <p class="text-xs text-zinc-400 leading-relaxed">
            Sign in with your TriGameDev Discord account to donate or claim game keys into your permanent
            <strong class="text-white">My Keys</strong> vault and download community ebooks.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <button
            @click="activeSection = 'keys'"
            class="p-4 bg-[#09090b] border border-zinc-800 hover:border-white text-left transition cursor-pointer"
          >
            <KeyRound class="w-5 h-5 text-[#ef4444] mb-2" />
            <div class="font-bold text-sm text-white uppercase">Key Sharing Vault</div>
            <div class="text-xs font-mono text-zinc-400 mt-1">
              {{ overview.stats?.available_keys || 0 }} KEYS AVAILABLE ->
            </div>
          </button>

          <button
            @click="activeSection = 'library'"
            class="p-4 bg-[#09090b] border border-zinc-800 hover:border-white text-left transition cursor-pointer"
          >
            <BookOpen class="w-5 h-5 text-white mb-2" />
            <div class="font-bold text-sm text-white uppercase">Ebook Library</div>
            <div class="text-xs font-mono text-zinc-400 mt-1">
              {{ overview.stats?.total_ebooks || 0 }} DIRECT-DOWNLOAD BOOKS ->
            </div>
          </button>
        </div>
      </div>
    </section>

    <!-- Bottom Discord Connect Bar -->
    <section
      v-if="!currentUser"
      class="p-4 sm:px-6 bg-[#121215] border border-zinc-800 flex flex-col sm:flex-row items-center justify-between gap-4"
    >
      <div class="flex items-center gap-3">
        <div class="w-2.5 h-2.5 bg-[#ef4444]"></div>
        <p class="text-xs sm:text-sm text-zinc-300">
          Member of the <strong class="text-white">TriGameDev Discord</strong>? Sign in to unlock the Key Vault, Ebook Library, and OG CMS.
        </p>
      </div>
      <a
        href="/api/go/auth/discord/login"
        class="px-5 py-2.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider whitespace-nowrap transition"
      >
        Sign In with Discord
      </a>
    </section>
  </div>
</template>
