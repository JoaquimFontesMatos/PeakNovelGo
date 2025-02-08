<script setup lang="ts">
const currentPage = ref(1);
const currentLimit = ref(10);

onMounted(async () => {
  try {
    await onPageChange(currentPage.value, currentLimit.value);
  } catch {}
});

const { fetchingNovel, paginatedNovelsData } = storeToRefs(useNovelStore());

const onPageChange = async (newPage: number, limit: number) => {
  if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsData.value && paginatedNovelsData.value.page != 0) return;
  try {
    await useNovelStore().fetchNovels(newPage, limit);
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
      ]"
    />

    <VerticalSpacer />

    <LoadingBar v-show="fetchingNovel" />

    <PaginatedNovelGallery
      v-show="!fetchingNovel"
      :errorMessage="paginatedNovelsData === null ? 'No Novels Found' : null"
      :paginatedData="paginatedNovelsData"
      @page-change="onPageChange"
    />
  </Container>
</template>
