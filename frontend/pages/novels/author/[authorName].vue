<script setup lang="ts">
const { authorName } = useRoute().params;
const currentPage = ref(1);
const currentLimit = ref(10);

onMounted(async () => {
  try {
    await onPageChange(currentPage.value, currentLimit.value);
  } catch {}
});

const { fetchingNovel, paginatedNovelsDataByAuthor } = storeToRefs(useNovelStore());

const onPageChange = async (newPage: number, limit: number) => {
  if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsDataByAuthor.value && paginatedNovelsDataByAuthor.value.page != 0) return;
  try {
    await useNovelStore().fetchNovelsByAuthor(authorName as string, newPage, limit);
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
          path: '/novels/author/' + authorName,
          name: authorName as string
        },
      ]"
    />

    <VerticalSpacer />

    <LoadingBar v-show="fetchingNovel" />

    <PaginatedNovelGallery
      v-show="!fetchingNovel"
      :errorMessage="paginatedNovelsDataByAuthor === null ? 'No Novels Found' : null"
      :paginatedData="paginatedNovelsDataByAuthor"
      @page-change="onPageChange"
    />
  </Container>
</template>
