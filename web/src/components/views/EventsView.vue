<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-5">
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
</template>
