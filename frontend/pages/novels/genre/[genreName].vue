<script setup lang="ts">
const { genreName } = useRoute().params as { genreName: string };
const currentPage = ref(1);
const currentLimit = ref(10);

onMounted(async () => {
  try {
    await onPageChange(currentPage.value, currentLimit.value);
  } catch {}
});

const { fetchingNovel, paginatedNovelsDataByGenre } = storeToRefs(useNovelStore());

const onPageChange = async (newPage: number, limit: number) => {
  if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsDataByGenre.value && paginatedNovelsDataByGenre.value.page != 0) return;
  try {
    await useNovelStore().fetchNovelsByGenre(genreName as string, newPage, limit);
  } catch {}
  currentPage.value = newPage;
  currentLimit.value = limit;
};
</script>

<template>
  <Container>
    <RouteTree
      :routes="[
        { path: '/', name: 'Home' },
        { path: '/novels', name: 'Novels' },
        {
          path: '/novels/genre/' + genreName,
          name: genreName as string
        },
      ]"
    />

    <VerticalSpacer />

    <LoadingBar v-show="fetchingNovel" />

    <PaginatedNovelGallery
      v-show="!fetchingNovel"
      :errorMessage="paginatedNovelsDataByGenre === null ? 'No Novels Found' : null"
      :paginatedData="paginatedNovelsDataByGenre"
      @page-change="onPageChange"
    />
  </Container>
</template>
