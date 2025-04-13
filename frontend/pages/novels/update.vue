<script setup lang="ts">
    import { storeToRefs } from 'pinia';

    const novelStore = useNovelStore();
    const { updatingNovels, novelStatuses } = storeToRefs(novelStore);
    const toastStore = useToastStore();

    const showPreview = ref(false);
    const errorMessage = ref('');

    // Pagination state
    const currentPage = ref(1);
    const pageSize = ref(10);

    // Computed property for paginated chapter statuses
    const paginatedStatuses = computed(() => {
        const entries = Object.entries(novelStatuses.value);

        const toUpdateNovels = entries.filter(([_, status]) => status === 'to update');
        const queuedNovels = entries.filter(([_, status]) => status === 'in queue');
        const updatingNovels = entries.filter(([_, status]) => status === 'updating');
        const errorNovels = entries.filter(([_, status]) => status === 'error');
        const completedNovels = entries.filter(([_, status]) => status === 'updated' || status === 'skipped');

        // Combine the groups in the desired order
        const sortedEntries = [...updatingNovels, ...queuedNovels, ...toUpdateNovels, ...errorNovels, ...completedNovels];

        const startIndex = (currentPage.value - 1) * pageSize.value;
        const endIndex = startIndex + pageSize.value;
        return new Map(
            sortedEntries.slice(startIndex, endIndex).map(([novelId, status]) => [novelId, status]) // Convert keys to numbers
        );
    });

    const completedNovelCount = computed(() => {
        return Object.keys(novelStatuses.value).filter(novelId => novelStatuses.value[novelId] === 'updated' || novelStatuses.value[novelId] === 'skipped')
            .length;
    });

    // Computed property for total pages
    const totalPages = computed(() => {
        return Math.ceil(Object.keys(novelStatuses.value).length / pageSize.value);
    });

    // Handle page change
    const handlePageChange = async (newPage: number, newPageSize: number): Promise<void> => {
        currentPage.value = newPage;
        pageSize.value = newPageSize;
    };

    const onSubmit = async () => {
        errorMessage.value = ''; // Clear error on valid input
        try {
            await novelStore.batchUpdateNovels();
            showPreview.value = true; // Show preview after successful import
        } catch (error) {
            errorMessage.value = 'Failed to update novels: ' + error;
            toastStore.addToast(errorMessage.value, 'error', 'project');
        }
    };

    const retry = async (novelUpdatesID: string) => {
        try {
            await novelStore.importByNovelUpdatesId(novelUpdatesID);
            novelStatuses.value[novelUpdatesID] = 'updated';
        } catch (error) {
            toastStore.addToast(`Retry didn't work for novel: ${novelUpdatesID}`, 'error', 'novel');
        }
    };
</script>

<template>
    <Container>
        <!-- Input Section -->
        <div>
            <!-- Import Button -->
            <MainButton :disabled="updatingNovels" @click="onSubmit" class="w-full md:w-auto">
                <div v-if="updatingNovels" class="flex items-center gap-2">
                    <LoadingSpinner class="h-5 w-5" />
                    <span>Updating Novels...</span>
                </div>
                <span v-else>Update Novels</span>
            </MainButton>
        </div>

        <VerticalSpacer />

        <!-- Status Updates Section -->
        <section v-if="updatingNovels || showPreview" class="space-y-4">
            <h2 class="text-lg font-semibold">Novel Update Status</h2>

            <!-- Progress Bar -->
            <div class="relative h-4 w-full overflow-hidden rounded-full bg-secondary">
                <div
                    class="h-full bg-secondary-content transition-all duration-300"
                    :style="{
                        width: `${(completedNovelCount / Object.keys(novelStatuses).length) * 100}%`,
                    }"
                ></div>
            </div>
            <p>
                Downloading: {{ completedNovelCount }} /
                {{ Object.keys(novelStatuses).length }}
            </p>

            <table class="w-full border-collapse">
                <thead>
                    <tr class="bg-secondary">
                        <th class="p-2 text-left">NovelUpdatesID</th>
                        <th class="p-2 text-left">Status</th>
                        <th class="p-2 text-left">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="[novelUpdatesID, status] in paginatedStatuses" :key="novelUpdatesID" class="border-b">
                        <td class="p-2">{{ novelUpdatesID }}</td>
                        <td
                            class="p-2"
                            :class="{
                                'text-primary-content': status === 'to download' || status === 'to update',
                                'text-yellow-500': status === 'in queue',
                                'text-blue-500': status === 'updating',
                                'text-green-500': status === 'updated',
                                'text-red-500': status === 'error',
                                'text-gray-500': status === 'skipped',
                            }"
                        >
                            {{ status }}
                        </td>
                        <td class="p-2">
                            <button v-if="novelStatuses[novelUpdatesID] === 'error'" @click="retry(novelUpdatesID)">Retry</button>
                        </td>
                    </tr>
                </tbody>
            </table>

            <!-- Pagination -->
            <NumberedPaginator :totalPages="totalPages" :total="Object.keys(novelStatuses).length" @page-change="handlePageChange" />
        </section>

        <!-- Preview Section -->
        <section v-if="showPreview && !novelStatuses" class="space-y-4">
            <VerticalSpacer />
            <!-- Navigation Button -->
            <MainButton @click="navigateTo('/')" class="w-full md:w-auto">
                <span>Go to Home</span>
            </MainButton>
        </section>
    </Container>
</template>
