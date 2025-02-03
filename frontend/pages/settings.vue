<script setup lang="ts">
definePageMeta({
    middleware: 'auth'
})
</script>

<template>
    <Container>
        <ReadingPreferencesForm/>
    </Container>
</template>
