<template>
  <div class="home-page">

    <Container>
      <section class="featured-novels py-12">
        <h2 class="mb-6 text-3xl font-bold text-primary-content">Featured Novels</h2>
        <PaginatedNovelGallery v-show="!fetchingNovel" :errorMessage="novelError" :paginatedData="paginatedNovelsData" @page-change="onPageChange" />
      </section>

      <section class="recently-updated py-12">
        <h2 class="mb-6 text-3xl font-bold text-primary-content">Recently Updated</h2>
        <PaginatedNovelGallery v-show="!fetchingNovel" :errorMessage="novelError" :paginatedData="paginatedNovelsData" @page-change="onPageChange" />
      </section>

      <section class="popular-tags py-12">
        <h2 class="mb-6 text-3xl font-bold text-primary-content">Popular Tags</h2>
        <div class="flex flex-wrap gap-2">
          <NuxtLink
            v-for="tag in popularTags"
            :key="tag"
            :to="`/novels/tag/${tag}`"
            class="rounded-full bg-secondary px-4 py-2 text-primary-content transition-colors duration-200 hover:bg-accent-gold hover:text-primary"
          >
            {{ tag }}
          </NuxtLink>
        </div>
      </section>

      <section class="top-rated py-12">
        <h2 class="mb-6 text-3xl font-bold text-primary-content">Top Rated</h2>
        <PaginatedNovelGallery v-show="!fetchingNovel" :errorMessage="novelError" :paginatedData="paginatedNovelsData" @page-change="onPageChange" />
      </section>
    </Container>
  </div>
</template>

<script setup lang="ts">
import { string } from 'yup';

const currentPage = ref(1);
const currentLimit = ref(10);

onMounted(async () => {
  await onPageChange(currentPage.value, currentLimit.value);
  // featuredNovels.value = await fetchFeaturedNovels();

  // Fetch recently updated novels (replace with your actual API call)
  // recentlyUpdated.value = await fetchRecentlyUpdatedNovels();

  // Fetch popular tags (replace with your actual API call)
  popularTags.value = await fetchPopularTags();

  // topRatedNovels.value = await fetchTopRatedNovels();
});

const { fetchingNovel, novelError, paginatedNovelsData } = storeToRefs(useNovelStore());

const onPageChange = async (newPage: number, limit: number) => {
  if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsData.value.page != 0) return;
  await useNovelStore().fetchNovels(newPage, limit);
  currentPage.value = newPage;
  currentLimit.value = limit;
};

//import HeroSection from '~/components/HeroSection.vue';  // Import the hero section

const featuredNovels = ref([]);
const recentlyUpdated = ref([]);
const popularTags: Ref<string[]> = ref([]);
const topRatedNovels = ref([]);

// Placeholder functions - replace these with your actual API calls
async function fetchFeaturedNovels() {
  return [];
}

async function fetchRecentlyUpdatedNovels() {
  return [];
}
async function fetchTopRatedNovels() {
  return [];
}

async function fetchPopularTags(): Promise<string[]> {
  return [
    'Adapted to Manhua',
    'Transformation Ability',
    'Transmigration',
    'Transplanted Memories',
    'Vampires',
    'Weak to Strong',
    'Wealthy Characters',
    'Werebeasts',
    'Zombies',
  ];
}
</script>
