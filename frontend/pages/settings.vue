<script setup lang="ts">
definePageMeta({
    middleware: 'auth'
})
</script>

<template>
    <main class="px-20 py-10">
        <LazyReadingPreferencesForm/>
    </main>
</template>
