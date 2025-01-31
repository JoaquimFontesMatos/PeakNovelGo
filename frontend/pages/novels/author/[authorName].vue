<script setup lang="ts">
const { authorName } = useRoute().params;
const currentPage = ref(1);
const currentLimit = ref(10);

onMounted(async() => {
  await onPageChange(currentPage.value, currentLimit.value);
});

const {
    fetchingNovel, novelError, paginatedNovelsDataByAuthor
  } = storeToRefs(
    useNovelStore()
  );

const onPageChange = async(newPage: number, limit: number) => {
  if (
    newPage === currentPage.value &&
    limit === currentLimit.value &&
    paginatedNovelsDataByAuthor.value.page != 0
    )
    return ;
  await useNovelStore().fetchNovelsByAuthor(
    authorName as string,
    newPage,
    limit
  );
  currentPage.value = newPage;
  currentLimit.value = limit;
};
</script>

<template>
  <main class="px-20 py-10">
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

    <VerticalSpacer/>

    <LoadingBar v-show="fetchingNovel"/>

    <PaginatedNovelGallery
      v-show="!fetchingNovel"
      :errorMessage="novelError"
      :paginatedData="paginatedNovelsDataByAuthor"
      @page-change="onPageChange"
    />
  </main>
</template>
