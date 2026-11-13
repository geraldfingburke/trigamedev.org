<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-5">
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
</template>
