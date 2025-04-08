<template>
  <div class="home-page">
    <Container>
      <DraggableColumn v-model="items"/>
    </Container>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";

const currentPage = ref(1);
const currentLimit = ref(10);

useHead({
  title: "📖 Home",
});

onMounted(async () => {
  await onPageChange(currentPage.value, currentLimit.value);
  // featuredNovels.value = await fetchFeaturedNovels();

  // Fetch recently updated novels (replace with your actual API call)
  // recentlyUpdated.value = await fetchRecentlyUpdatedNovels();

  // Fetch popular tags (replace with your actual API call)
  popularTags.value = await fetchPopularTags();

  // topRatedNovels.value = await fetchTopRatedNovels();
});

const {fetchingNovel, paginatedNovelsData} = storeToRefs(useNovelStore());

const onPageChange = async (newPage: number, limit: number) => {
  if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsData.value && paginatedNovelsData.value.page != 0) return;

  try {
    await useNovelStore().fetchNovels(newPage, limit);
  } catch {
  }
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

const items = ref([
  {
    id: '1',
    type: 'paginated-novel-gallery',
    props: {
      errorMessage: paginatedNovelsData === null ? 'No Novels Found' : null,
      paginatedData: paginatedNovelsData,
      onPageChange: onPageChange,
    }
  },
  {
    id: '2',
    type: 'paginated-novel-gallery',
    props: {
      errorMessage: paginatedNovelsData === null ? 'No Novels Found' : null,
      paginatedData: paginatedNovelsData,
      onPageChange: onPageChange,
    }
  },
  {
    id: '3',
    type: 'paginated-novel-gallery',
    props: {
      errorMessage: paginatedNovelsData === null ? 'No Novels Found' : null,
      paginatedData: paginatedNovelsData,
      onPageChange: onPageChange,
    }
  }
]);
</script>
