<script setup lang="ts">
    import { storeToRefs } from 'pinia';

    const chapterStore = useChapterStore();
    const { importingChapters, chapterStatuses } = storeToRefs(chapterStore);
    const toastStore = useToastStore();

    const novelUpdatesId = ref('');
    const showPreview = ref(false);
    const errorMessage = ref('');

    // Pagination state
    const currentPage = ref(1);
    const pageSize = ref(10);

    // Computed property for paginated chapter statuses
    const paginatedChapterStatuses = computed(() => {
        const entries = Object.entries(chapterStatuses.value);

        const downloadingChapters = entries.filter(([_, status]) => status === 'downloading');
        const queuedNovels = entries.filter(([_, status]) => status === 'in queue');
        const toDownloadChapters = entries.filter(([_, status]) => status === 'to download');
        const errorChapters = entries.filter(([_, status]) => status === 'error');
        const completedChapters = entries.filter(([_, status]) => status === 'downloaded' || status === 'skipped');

        // Combine the groups in the desired order
        const sortedEntries = [...downloadingChapters, ...queuedNovels, ...toDownloadChapters, ...errorChapters, ...completedChapters];

        const startIndex = (currentPage.value - 1) * pageSize.value;
        const endIndex = startIndex + pageSize.value;
        return new Map(
            sortedEntries.slice(startIndex, endIndex).map(([chapterNoStr, status]) => [Number(chapterNoStr), status]) // Convert keys to numbers
        );
    });

    const completedChaptersCount = computed(() => {
        return Object.keys(chapterStatuses.value).filter(
            chapterNoStr => chapterStatuses.value[Number(chapterNoStr)] === 'downloaded' || chapterStatuses.value[Number(chapterNoStr)] === 'skipped'
        ).length;
    });

    // Computed property for total pages
    const totalPages = computed(() => {
        return Math.ceil(Object.keys(chapterStatuses.value).length / pageSize.value);
    });

    // Handle page change
    const handlePageChange = async (newPage: number, newPageSize: number): Promise<void> => {
        currentPage.value = newPage;
        pageSize.value = newPageSize;
    };

    const onSubmit = async () => {
        if (!novelUpdatesId.value.trim()) {
            errorMessage.value = 'Please enter a Novel Updates ID';
            toastStore.addToast(errorMessage.value, 'error', 'project');
            return;
        }

        errorMessage.value = ''; // Clear error on valid input
        try {
            await chapterStore.importChapters(novelUpdatesId.value);
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
                <label for="novelUpdatesId" class="block text-sm font-medium after:text-error after:content-['*']">Novel Updates ID</label>
                <div class="relative">
                    <input
                        id="novelUpdatesId"
                        v-model="novelUpdatesId"
                        type="text"
                        placeholder="e.g., naruto"
                        class="w-full rounded-md border p-2 focus:ring-2 focus:ring-primary focus:outline-hidden"
                        :class="{ 'border-error': errorMessage }"
                    />
                    <p v-if="errorMessage" class="mt-1 text-sm text-error">
                        {{ errorMessage }}
                    </p>
                </div>
            </div>

            <VerticalSpacer />

            <!-- Import Button -->
            <MainButton :disabled="importingChapters" @click="onSubmit" class="w-full md:w-auto">
                <div v-if="importingChapters" class="flex items-center gap-2">
                    <LoadingSpinner class="h-5 w-5" />
                    <span>Importing Chapters...</span>
                </div>
                <span v-else>Import Chapters</span>
            </MainButton>
        </div>

        <VerticalSpacer />

        <!-- Status Updates Section -->
        <section v-if="importingChapters || showPreview" class="space-y-4">
            <h2 class="text-lg font-semibold">Chapter Import Status</h2>

            <!-- Progress Bar -->
            <div class="relative h-4 w-full overflow-hidden rounded-full bg-secondary">
                <div
                    class="h-full bg-secondary-content transition-all duration-300"
                    :style="{
                        width: `${(completedChaptersCount / Object.keys(chapterStatuses).length) * 100}%`,
                    }"
                ></div>
            </div>
            <p>
                Downloading: {{ completedChaptersCount }} /
                {{ Object.keys(chapterStatuses).length }}
            </p>

            <table class="w-full border-collapse">
                <thead>
                    <tr class="bg-secondary">
                        <th class="p-2 text-left">Chapter No</th>
                        <th class="p-2 text-left">Status</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="[chapterNo, status] in paginatedChapterStatuses" :key="chapterNo" class="border-b">
                        <td class="p-2">Chapter {{ chapterNo }}</td>
                        <td
                            class="p-2"
                            :class="{
                                'text-primary-content': status === 'to download',
                                'text-yellow-500': status === 'in queue',
                                'text-blue-500': status === 'downloading',
                                'text-green-500': status === 'downloaded',
                                'text-red-500': status === 'error',
                                'text-gray-500': status === 'skipped',
                            }"
                        >
                            {{ status }}
                        </td>
                    </tr>
                </tbody>
            </table>

            <!-- Pagination -->
            <NumberedPaginator :totalPages="totalPages" :total="Object.keys(chapterStatuses).length" @page-change="handlePageChange" />
        </section>

        <!-- Preview Section -->
        <section v-if="showPreview && !importingChapters" class="space-y-4">
            <VerticalSpacer />
            <!-- Navigation Button -->
            <MainButton @click="navigateTo('/')" class="w-full md:w-auto">
                <span>Go to Home</span>
            </MainButton>
        </section>
    </Container>
</template>
