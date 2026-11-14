<script setup>
import {
  Calendar,
  MapPin,
  ChevronRight,
  Eye,
  Trophy,
  Newspaper,
  ExternalLink
} from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const {
  activeSection,
  overview,
  currentUser,
  staySignedIn,
  activeJam,
  upcomingEvents,
  filteredShowcase,
  selectedWinnersJam,
  loginWithDiscord
} = useTriGameDevStore()

function formatEventDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div class="flex flex-col gap-8">
    <!-- 1. What, Who, and Why is TriGameDev (At the Top) -->
    <section class="grid grid-cols-1 md:grid-cols-3 border border-zinc-800 bg-[#121215] divide-y md:divide-y-0 md:divide-x divide-zinc-800">
      <div class="p-6 flex flex-col gap-2">
        <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">WHAT IS TRIGAMEDEV?</div>
        <h2 class="text-base font-bold text-white uppercase">A Loose Collective</h2>
        <p class="text-xs text-zinc-300 leading-relaxed">
          TriGameDev is a loose collective of game developers that started in the Tri-Cities Region of Northeast Tennessee. While many of our members are still there, there are plenty of us that are not. All are welcome!
        </p>
      </div>

      <div class="p-6 flex flex-col gap-2">
        <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">WHO IS TRIGAMEDEV?</div>
        <h2 class="text-base font-bold text-white uppercase">Pros, Hobbyists &amp; Players</h2>
        <p class="text-xs text-zinc-300 leading-relaxed">
          Pros, Hobbyists, Students, Curious Non-Developers, and Game Playing People.
        </p>
      </div>

      <div class="p-6 flex flex-col gap-2">
        <div class="text-[11px] font-mono uppercase tracking-widest text-[#ef4444]">WHY IS TRIGAMEDEV?</div>
        <h2 class="text-base font-bold text-white uppercase">Let's Be Alone Together</h2>
        <p class="text-xs text-zinc-300 leading-relaxed">
          Solo game dev can be lonely. Let's be alone together.
        </p>
      </div>
    </section>

    <!-- 2. Upcoming Events -->
    <section class="flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <h2 class="text-lg font-bold text-white uppercase tracking-wider flex items-center gap-2">
          <Calendar class="w-4 h-4 text-[#ef4444]" />
          <span>Upcoming Events</span>
        </h2>
        <button
          @click="activeSection = 'events'"
          class="px-3.5 py-1.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider flex items-center gap-1 transition cursor-pointer"
        >
          <span>All Events</span>
          <ChevronRight class="w-3.5 h-3.5" />
        </button>
      </div>

      <div
        v-if="upcomingEvents.length === 0"
        class="p-6 bg-[#121215] border border-zinc-800 flex items-center justify-between text-xs font-mono uppercase text-zinc-400"
      >
        <span>No upcoming events scheduled right now.</span>
        <button @click="activeSection = 'events'" class="text-white hover:underline cursor-pointer">
          View Past Events ->
        </button>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div
          v-for="ev in upcomingEvents.slice(0, 3)"
          :key="ev.id"
          class="bg-[#121215] border border-zinc-800 p-4 flex flex-col justify-between gap-3"
        >
          <div>
            <div class="flex items-center justify-between text-[11px] font-mono">
              <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[#ef4444] uppercase">{{ ev.tag }}</span>
              <span class="text-zinc-400">{{ formatEventDate(ev.start_time) }}</span>
            </div>
            <h3 class="font-bold text-sm text-white uppercase mt-2">{{ ev.title }}</h3>
            <p class="text-xs text-zinc-400 mt-1 line-clamp-2">{{ ev.description }}</p>
          </div>
          <div class="pt-2 border-t border-zinc-800 flex items-center justify-between text-[11px] font-mono text-zinc-300">
            <span class="flex items-center gap-1 truncate">
              <MapPin class="w-3 h-3 text-[#ef4444] shrink-0" />
              <span class="truncate">{{ ev.location }}</span>
            </span>
            <span class="shrink-0">{{ ev.rsvp_count }} GOING</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 3. Latest Jam / Active Jam -->
    <section class="flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <h2 class="text-lg font-bold text-white uppercase tracking-wider">
          {{ activeJam?.status === 'active' ? 'Active Jam' : 'Latest Jam' }}
        </h2>
        <button
          @click="activeSection = 'jams'"
          class="px-3.5 py-1.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider flex items-center gap-1 transition cursor-pointer"
        >
          <span>All Jams</span>
          <ChevronRight class="w-3.5 h-3.5" />
        </button>
      </div>

      <div
        v-if="!activeJam"
        class="p-6 bg-[#121215] border border-zinc-800 text-xs font-mono uppercase text-zinc-500 text-center"
      >
        No jams posted yet.
      </div>

      <div
        v-else
        class="bg-[#121215] border border-zinc-800 flex flex-col md:flex-row items-stretch justify-between"
      >
        <img
          :src="activeJam.cover_image_url || '/brand/filler.png'"
          :alt="activeJam.title"
          class="w-full md:w-72 h-48 md:h-auto object-cover bg-black border-b md:border-b-0 md:border-r border-zinc-800 shrink-0"
        />
        <div class="p-6 flex-1 flex flex-col justify-between gap-4">
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[11px] font-mono uppercase text-[#ef4444]">
                {{ activeJam.month_label }}
              </span>
              <span class="px-2 py-0.5 bg-black border border-zinc-700 text-[11px] font-mono uppercase text-white">
                {{ activeJam.status === 'active' ? 'ACTIVE JAM' : 'LATEST JAM' }}
              </span>
            </div>
            <h3 class="text-xl font-extrabold text-white uppercase">{{ activeJam.title }}</h3>
            <div class="text-xs font-mono text-zinc-300 bg-[#09090b] border border-zinc-800 px-3 py-1.5 w-fit">
              THEME: <strong class="text-white">{{ activeJam.theme }}</strong>
            </div>
          </div>

          <div class="pt-3 border-t border-zinc-800 flex flex-wrap items-center justify-between gap-3 text-xs">
            <span class="font-mono text-zinc-400">
              {{ activeJam.participants_count }} JAMMERS // {{ activeJam.submissions_count }} ENTRIES
            </span>
            <div class="flex items-center gap-2">
              <button
                @click="selectedWinnersJam = activeJam"
                class="px-3.5 py-2 bg-[#09090b] hover:bg-zinc-900 text-white border border-zinc-700 font-bold text-xs uppercase tracking-wider flex items-center gap-1.5 cursor-pointer transition"
              >
                <Trophy class="w-3.5 h-3.5 text-[#ef4444]" />
                <span>Winners Circle</span>
              </button>
              <a
                :href="activeJam.itch_url"
                target="_blank"
                class="px-4 py-2 bg-white hover:bg-zinc-200 text-black font-bold text-xs uppercase tracking-wider flex items-center gap-1 transition"
              >
                <span>Itch.io Page</span>
                <ExternalLink class="w-3.5 h-3.5" />
              </a>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. Latest Games -->
    <section class="flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
        <h2 class="text-lg font-bold text-white uppercase tracking-wider">Latest Games</h2>
        <button
          @click="activeSection = 'showcase'"
          class="px-3.5 py-1.5 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider flex items-center gap-1 transition cursor-pointer"
        >
          <span>Games by Our People</span>
          <ChevronRight class="w-3.5 h-3.5" />
        </button>
      </div>

      <div
        v-if="filteredShowcase.length === 0"
        class="p-8 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No games published yet.
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <article
          v-for="game in filteredShowcase.slice(0, 3)"
          :key="game.id"
          class="group bg-[#121215] border border-zinc-800 hover:border-zinc-500 transition flex flex-col justify-between"
        >
          <div class="relative h-48 overflow-hidden bg-black border-b border-zinc-800">
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

    <!-- 5. News & Announcements (Last, plus Guest-Only Discord Member Perks) -->
    <section class="grid grid-cols-1 gap-5" :class="!currentUser ? 'xl:grid-cols-3' : ''">
      <div
        class="bg-[#121215] p-5 border border-zinc-800 flex flex-col gap-4"
        :class="!currentUser ? 'xl:col-span-2' : ''"
      >
        <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
          <h2 class="font-bold text-base text-white uppercase tracking-wider flex items-center gap-2">
            <Newspaper class="w-4 h-4 text-[#ef4444]" />
            <span>News &amp; Announcements</span>
          </h2>
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

      <!-- Guest-Only Discord Member Perks -->
      <div
        v-if="!currentUser"
        class="bg-[#121215] p-6 border border-zinc-800 flex flex-col justify-between gap-4"
      >
        <div class="flex flex-col gap-2">
          <span class="px-2 py-0.5 w-fit text-[10px] font-mono uppercase tracking-widest bg-zinc-900 border border-zinc-700 text-[#ef4444]">
            DISCORD MEMBER PERKS
          </span>
          <h3 class="text-xl font-extrabold text-white uppercase mt-2">
            Access to Free Games, Books, and More...
          </h3>
        </div>

        <div class="flex flex-col gap-3 pt-3 border-t border-zinc-800">
          <button
            @click="loginWithDiscord"
            class="w-full py-2.5 px-4 bg-white hover:bg-zinc-200 text-black text-xs font-bold uppercase tracking-wider cursor-pointer transition"
          >
            Sign In with Discord
          </button>
          <label class="flex items-center gap-2 text-xs font-mono text-zinc-400 cursor-pointer select-none">
            <input
              v-model="staySignedIn"
              type="checkbox"
              class="accent-[#ef4444] w-3.5 h-3.5 cursor-pointer"
            />
            <span>Stay signed in</span>
          </label>
        </div>
      </div>
    </section>
  </div>
</template>
