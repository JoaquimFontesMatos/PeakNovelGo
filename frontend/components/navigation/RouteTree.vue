<script setup lang="ts">
    defineProps<{
        routes: { path: string; name: string }[];
    }>();
</script>

<template>
    <div class="flex items-center">
        <div v-for="({ path, name }, index) in routes" :key="path">
            <div>
                <div v-if="index !== routes.length - 1" class="flex flex-row text-lg text-secondary hover:text-accent-gold md:text-xl">
                    <NuxtLink :to="path" class="line-clamp-1 hover:underline">
                        {{ name }}
                    </NuxtLink>
                    <span class="mx-2">/</span>
                </div>
                <div v-else class="h-full">
                    <span class="md:text-md line-clamp-1 text-accent-gold md:text-lg">
                        {{ name }}
                    </span>
                </div>
            </div>
        </div>
    </div>
</template>
