<script setup>
import { MapPin } from 'lucide-vue-next'
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-5">
    <h1 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3">
      Upcoming Events & Discord Stages
    </h1>
    <div
      v-if="!overview.events || overview.events.length === 0"
      class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
    >
      No upcoming events scheduled yet.
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div
        v-for="ev in overview.events"
        :key="ev.id"
        class="bg-[#121215] border border-zinc-800 flex flex-col justify-between"
      >
        <img :src="ev.cover_image_url || '/brand/first-meetup.jpg'" class="h-40 w-full object-cover border-b border-zinc-800 bg-black" />
        <div class="p-5 flex-1 flex flex-col justify-between gap-4">
          <div>
            <span class="px-2 py-0.5 text-[11px] font-mono uppercase bg-zinc-900 border border-zinc-700 text-[#ef4444]">
              {{ ev.tag }}
            </span>
            <h3 class="font-bold text-base text-white uppercase mt-2">{{ ev.title }}</h3>
            <p class="text-xs text-zinc-400 mt-1">{{ ev.description }}</p>
          </div>
          <div class="pt-3 border-t border-zinc-800 text-xs font-mono text-zinc-300 flex items-center justify-between">
            <span class="flex items-center gap-1.5">
              <MapPin class="w-3.5 h-3.5 text-[#ef4444]" />
              <span>{{ ev.location }}</span>
            </span>
            <span>{{ ev.rsvp_count }} GOING</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
