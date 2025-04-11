<script setup lang="ts">
    const route = useRoute();
    const router = useRouter();

    const currentPage = ref(Number(route.query.page) || 1);
    const selectedPageSize = ref(Number(route.query.pageSize) || 10);

    onMounted(() => {
        if (!route.query.page) currentPage.value = 1;
        if (!route.query.pageSize) selectedPageSize.value = 10;
        props.onPageChange(currentPage.value, selectedPageSize.value);
    });

    const props = defineProps<{
        totalPages: number;
        total: number;
        onPageChange: (newPage: number, limit: number) => Promise<void>;
    }>();

    function updateQueryParams(page: number, limit: number) {
        router.replace({
            query: { ...route.query, page, pageSize: limit },
        });
    }

    function goToNextPage() {
        if (currentPage.value < props.totalPages) {
            currentPage.value += 1;
            updateQueryParams(currentPage.value, selectedPageSize.value);
            props.onPageChange(currentPage.value, selectedPageSize.value);
        }
    }

    function goToPreviousPage() {
        if (currentPage.value > 1) {
            currentPage.value -= 1;
            updateQueryParams(currentPage.value, selectedPageSize.value);
            props.onPageChange(currentPage.value, selectedPageSize.value);
        }
    }

    function goToSelectedPage(newPage: number, limit: number = selectedPageSize.value) {
        currentPage.value = newPage;
        selectedPageSize.value = limit;
        updateQueryParams(newPage, limit);
        props.onPageChange(newPage, limit);
    }

    /**
     * Calculate the visible page numbers.
     */
    function getVisiblePages(): number[] {
        const maxVisiblePages = 10;
        const pages = [];

        const startPage = Math.max(1, currentPage.value - Math.floor(maxVisiblePages / 2));
        const endPage = Math.min(props.totalPages, startPage + maxVisiblePages - 1);

        for (let i = startPage; i <= endPage; i++) {
            pages.push(i);
        }

        // Ensure the first and last pages are included.
        if (!pages.includes(1)) pages.unshift(1);
        if (!pages.includes(props.totalPages)) pages.push(props.totalPages);

        return pages;
    }

    /**
     * Add a debounce to the page size change to prevent rapid page changes
     */

    // Watch for changes in query params to keep the UI reactive
    watch(
        () => route.query,
        newQuery => {
            if (newQuery.page) currentPage.value = Number(newQuery.page);
            if (newQuery.pageSize) selectedPageSize.value = Number(newQuery.pageSize);
        }
    );
</script>

<template>
    <nav class="flex flex-wrap items-center justify-center gap-2">
        <button
            class="rounded-full border-2 border-transparent p-2 py-1 hover:enabled:border-accent-gold hover:enabled:bg-accent-gold-light hover:enabled:text-primary disabled:cursor-not-allowed disabled:text-secondary-content"
            :disabled="currentPage === 1"
            @click="goToPreviousPage"
        >
            <
        </button>

        <button
            v-for="(page, index) in getVisiblePages()"
            :key="page"
            @click="goToSelectedPage(page)"
            class="flex items-center justify-center gap-2 rounded-full border-2 border-transparent p-2 py-1 hover:enabled:border-accent-gold hover:enabled:bg-accent-gold-light hover:enabled:text-primary disabled:cursor-not-allowed disabled:text-secondary-content"
            :class="currentPage === page ? 'bg-accent-gold' : ''"
            :disabled="currentPage === page"
        >
            <span v-if="page === totalPages && page - 1 != getVisiblePages()[index - 1]">...</span>
            {{ page }}
            <span v-if="page === 1 && page + 1 != getVisiblePages()[index + 1]">...</span>
        </button>

        <button
            class="rounded-full border-2 border-transparent p-2 py-1 hover:enabled:border-accent-gold hover:enabled:bg-accent-gold-light hover:enabled:text-primary disabled:cursor-not-allowed disabled:text-secondary-content"
            :disabled="currentPage === totalPages"
            @click="goToNextPage"
        >
            >
        </button>

        <label for="page-size-select" class="text-secondary-content">Page Size:</label>
        <select
            id="page-size-select"
            class="w-16 rounded-full border-2 border-accent-gold bg-secondary p-2 py-1"
            v-model="selectedPageSize"
            @change="goToSelectedPage(1, selectedPageSize)"
        >
            <option value="10">10</option>
            <option value="25">25</option>
            <option value="50">50</option>
            <option value="100">100</option>
        </select>
    </nav>
</template>
