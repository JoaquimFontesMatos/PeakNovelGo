<script setup lang="ts">
const { tagName } = useRoute().params;
const currentPage = ref(1);
const currentLimit = ref(10);

onMounted(async() => {
  await onPageChange(currentPage.value, currentLimit.value);
});

const {
    fetchingNovel, novelError, paginatedNovelsDataByTag
  } = storeToRefs(
    useNovelStore()
  );

const onPageChange = async(newPage: number, limit: number) => {
  if (
    newPage === currentPage.value &&
    limit === currentLimit.value &&
    paginatedNovelsDataByTag.value.page != 0
    )
    return ;
  await useNovelStore().fetchNovelsByTag(tagName as string, newPage, limit);
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
          path: '/novels/tag/' + tagName,
          name: tagName as string
        },
      ]"
    />

    <VerticalSpacer/>

    <LoadingBar v-show="fetchingNovel"/>

    <PaginatedNovelGallery
      v-show="!fetchingNovel"
      :errorMessage="novelError"
      :paginatedData="paginatedNovelsDataByTag"
      @page-change="onPageChange"
    />
  </Container>
</template>
