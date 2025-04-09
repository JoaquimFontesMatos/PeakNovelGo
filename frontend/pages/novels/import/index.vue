<script setup lang="ts">
const novelStore = useNovelStore();
const {importingNovel, novel} = storeToRefs(novelStore);
const toastStore = useToastStore();

const novelUpdatesId = ref('');
const showPreview = ref(false); // Toggle for preview
const errorMessage = ref(''); // Error message for validation

const onSubmit = async () => {
  if (!novelUpdatesId.value.trim()) {
    errorMessage.value = 'Please enter a Novel Updates ID';
    toastStore.addToast(errorMessage.value, 'error', 'project');
    return;
  }

  errorMessage.value = ''; // Clear error on valid input
  try {
    await novelStore.importByNovelUpdatesId(novelUpdatesId.value);
    showPreview.value = true; // Show preview after successful import
  } catch (error) {
    errorMessage.value = 'Failed to import novel. Please check the ID and try again.';
    toastStore.addToast(errorMessage.value, 'error', 'project');
  }
};
</script>

<template>
  <Container>
    <!-- Input Section -->
    <div>
      <div class="w-full space-y-2 md:w-2/3">
        <label for="novelUpdatesId" class="block text-sm font-medium after:text-error after:content-['*']"> Novel
          Updates ID </label>
        <div class="relative">
          <input
              id="novelUpdatesId"
              v-model="novelUpdatesId"
              type="text"
              placeholder="e.g., reverend-insanity"
              class="w-full rounded-md border p-2 focus:outline-hidden focus:ring-2 focus:ring-primary"
              :class="{ 'border-error': errorMessage }"
          />
          <p v-if="errorMessage" class="mt-1 text-sm text-error">{{ errorMessage }}</p>
        </div>
      </div>

      <VerticalSpacer/>

      <!-- Import Button -->
      <Button :disabled="importingNovel" @click="onSubmit" class="w-full md:w-auto">
        <div v-if="importingNovel" class="flex items-center gap-2">
          <LoadingSpinner class="h-5 w-5"/>
          <span>Importing Novel...</span>
        </div>
        <span v-else>Import Novel</span>
      </Button>
    </div>

    <VerticalSpacer/>

    <!-- Preview Section -->
    <section v-if="novel" class="space-y-4">
      <div class="rounded-lg border-border bg-secondary p-4">
        <p class="text-lg font-semibold">Novel imported successfully!</p>
        <p class="text-secondary-content">You can close this window now.</p>
        <SmallVerticalSpacer/>
        <button @click="showPreview = !showPreview" class="w-full hover:text-accent-gold hover:underline md:w-auto">
          <span>{{ showPreview ? 'Hide Details' : 'Show Details' }}</span>
        </button>
      </div>

      <!-- Collapsible Novel Info -->
      <div class="space-y-2">
        <div v-if="showPreview" class="space-y-2 rounded-lg border-border bg-secondary p-4">
          <p><span class="font-medium">Novel ID:</span> {{ novel.ID }}</p>
          <p><span class="font-medium">Title:</span> {{ novel.title }}</p>
          <p><span class="font-medium">Synopsis:</span> {{ novel.synopsis }}</p>
          <p><span class="font-medium">Cover URL:</span> {{ novel.coverUrl }}</p>
          <p><span class="font-medium">Language:</span> {{ novel.language }}</p>
          <p><span class="font-medium">Status:</span> {{ novel.status }}</p>
          <p><span class="font-medium">Novel Updates URL:</span> {{ novel.novelUpdatesUrl }}</p>
          <span class="font-medium">Tag(s):</span>
          <ul class="flex max-h-28 flex-wrap overflow-scroll lg:max-h-64 xl:h-auto">
            <li v-for="({ name, id }, index) in novel.tags" :key="id" class="text-accent-gold">
              <NuxtLink :to="'/novels/tag/' + name" class="hover:underline">
                {{ name }}
              </NuxtLink>
              <span v-if="index !== novel.tags.length - 1" class="mr-2">,</span>
            </li>
          </ul>

          <span class="font-medium">Genre(s):</span>
          <ul class="flex max-h-28 flex-wrap overflow-scroll lg:max-h-64 xl:h-auto">
            <li v-for="({ name, id }, index) in novel.genres" :key="id" class="text-accent-gold">
              <NuxtLink :to="'/novels/genre/' + name" class="hover:underline">
                {{ name }}
              </NuxtLink>
              <span v-if="index !== novel.genres.length - 1" class="mr-2">,</span>
            </li>
          </ul>

          <span class="font-medium">Author(s):</span>
          <ul class="flex max-h-28 flex-wrap overflow-scroll lg:max-h-64 xl:h-auto">
            <li v-for="({ name, id }, index) in novel.authors" :key="id" class="text-accent-gold">
              <NuxtLink :to="'/novels/author/' + name" class="hover:underline">
                {{ name }}
              </NuxtLink>
              <span v-if="index !== novel.authors.length - 1" class="mr-2">,</span>
            </li>
          </ul>
          <p><span class="font-medium">Year:</span> {{ novel.year }}</p>
          <p><span class="font-medium">Release Frequency:</span> {{ novel.releaseFrequency }}</p>
          <p><span class="font-medium">Novel Updates ID:</span> {{ novel.novelUpdatesId }}</p>
          <p><span class="font-medium">Latest Chapter:</span> {{ novel.latestChapter }}</p>
        </div>
      </div>

      <VerticalSpacer/>

      <!-- Navigation Button -->
      <Button @click="navigateTo('/')" class="w-full md:w-auto">
        <span>Go to Home</span>
      </Button>
    </section>
  </Container>
</template>
