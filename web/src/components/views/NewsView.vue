<script setup>
import { useTriGameDevStore } from '../../composables/useTriGameDevStore'

const { overview, activeSlug, slugify, navigateTo } = useTriGameDevStore()
</script>

<template>
  <div class="flex flex-col gap-5">
    <h1 class="text-2xl font-extrabold text-white uppercase border-b border-zinc-800 pb-3">
      Community News & Announcements
    </h1>
    <div
      v-if="!overview.news || overview.news.length === 0"
      class="p-10 bg-[#121215] border border-zinc-800 text-center text-xs font-mono uppercase text-zinc-500"
    >
      No news or announcements published yet.
    </div>
    <div v-else class="flex flex-col gap-4">
      <article
        v-for="post in overview.news"
        :key="post.id"
        :id="'slug-' + (post.slug || slugify(post.title))"
        class="p-6 bg-[#121215] border flex flex-col gap-2 transition"
        :class="
          (post.slug || slugify(post.title)) === activeSlug
            ? 'border-[#ef4444] ring-1 ring-[#ef4444]'
            : 'border-zinc-800'
        "
      >
        <div class="flex items-center justify-between text-xs font-mono">
          <span class="px-2 py-0.5 bg-zinc-900 border border-zinc-700 text-[#ef4444] uppercase">{{ post.tag }}</span>
          <span class="text-zinc-500">{{ post.published_at }} // BY {{ post.author_name }} ({{ post.author_role }})</span>
        </div>
        <h2
          @click="navigateTo('news', post.slug || slugify(post.title))"
          class="text-lg font-bold text-white hover:text-[#ef4444] uppercase mt-1 cursor-pointer transition"
        >
          {{ post.title }}
        </h2>
        <p class="text-sm text-zinc-300 leading-relaxed whitespace-pre-line">{{ post.body_markdown }}</p>
      </article>
    </div>
  </div>
</template>
