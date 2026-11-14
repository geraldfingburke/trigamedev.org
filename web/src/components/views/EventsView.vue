<script setup>
import { MapPin, Calendar } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { upcomingEvents, pastEvents } = useTriGameDevStore()

function formatEventDateTime(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit'
  })
}
</script>

<template>
  <div class="flex flex-col gap-8">
    <h1 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3">
      Events
    </h1>

    <!-- Upcoming Events -->
    <section class="flex flex-col gap-4">
      <h2 class="text-sm font-mono uppercase tracking-widest text-[#ef4444]">
        // UPCOMING EVENTS ({{ upcomingEvents.length }})
      </h2>

      <div
        v-if="upcomingEvents.length === 0"
        class="p-8 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No upcoming events scheduled on Discord right now.
      </div>

      <div v-else class="flex flex-col gap-3">
        <div
          v-for="ev in upcomingEvents"
          :key="'up-' + ev.id"
          class="bg-[#121215] border border-zinc-800 p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
        >
          <div class="flex items-start gap-4">
            <img
              :src="ev.cover_image_url || '/brand/first-meetup.jpg'"
              class="w-24 h-20 object-cover border border-zinc-800 bg-black shrink-0"
            />
            <div>
              <div class="flex items-center gap-2 text-[11px] font-mono">
                <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[#ef4444] uppercase">{{ ev.tag }}</span>
                <span class="text-zinc-300 flex items-center gap-1">
                  <Calendar class="w-3 h-3 text-zinc-400" />
                  {{ formatEventDateTime(ev.start_time) }}
                </span>
              </div>
              <h3 class="font-bold text-base text-white uppercase mt-1">{{ ev.title }}</h3>
              <p class="text-xs text-zinc-400 mt-0.5">{{ ev.description }}</p>
            </div>
          </div>

          <div class="flex sm:flex-col items-end justify-between w-full sm:w-auto gap-1 shrink-0 font-mono text-xs">
            <span class="text-zinc-300 flex items-center gap-1">
              <MapPin class="w-3.5 h-3.5 text-[#ef4444]" />
              <span>{{ ev.location }}</span>
            </span>
            <span class="text-zinc-500">{{ ev.rsvp_count }} GOING</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Past Events -->
    <section class="flex flex-col gap-4">
      <h2 class="text-sm font-mono uppercase tracking-widest text-zinc-400">
        // PAST EVENTS ({{ pastEvents.length }})
      </h2>

      <div
        v-if="pastEvents.length === 0"
        class="p-8 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
      >
        No past events archived yet.
      </div>

      <div v-else class="flex flex-col gap-3">
        <div
          v-for="ev in pastEvents"
          :key="'past-' + ev.id"
          class="bg-[#121215] border border-zinc-800 p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 opacity-85"
        >
          <div class="flex items-start gap-4">
            <img
              :src="ev.cover_image_url || '/brand/first-meetup.jpg'"
              class="w-24 h-20 object-cover border border-zinc-800 bg-black shrink-0 grayscale"
            />
            <div>
              <div class="flex items-center gap-2 text-[11px] font-mono">
                <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-zinc-400 uppercase">{{ ev.tag }}</span>
                <span class="text-zinc-400">{{ formatEventDateTime(ev.start_time) }}</span>
              </div>
              <h3 class="font-bold text-base text-white uppercase mt-1">{{ ev.title }}</h3>
              <p class="text-xs text-zinc-400 mt-0.5">{{ ev.description }}</p>
            </div>
          </div>

          <div class="flex sm:flex-col items-end justify-between w-full sm:w-auto gap-1 shrink-0 font-mono text-xs">
            <span class="text-zinc-400 flex items-center gap-1">
              <MapPin class="w-3.5 h-3.5 text-zinc-500" />
              <span>{{ ev.location }}</span>
            </span>
            <span class="text-zinc-500">{{ ev.rsvp_count }} ATTENDED</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
