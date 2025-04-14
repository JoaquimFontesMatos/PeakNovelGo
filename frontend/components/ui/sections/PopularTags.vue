<script setup lang="ts">
    const popularTags: Ref<string[]> = ref([]);

    onMounted(async () => {
        popularTags.value = await fetchPopularTags();
    });

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
</script>

<template>
    <section class="">
        <h2 class="mb-6 text-3xl font-bold text-primary-content">Popular Tags</h2>
        <div class="flex flex-wrap gap-2">
            <NuxtLink
                v-for="tag in popularTags"
                :key="tag"
                :to="`/novels/tag/${tag}`"
                class="rounded-full bg-secondary px-4 py-2 text-primary-content transition-colors duration-200 hover:bg-accent-gold hover:text-primary"
            >
                {{ tag }}
            </NuxtLink>
        </div>
    </section>
</template>
