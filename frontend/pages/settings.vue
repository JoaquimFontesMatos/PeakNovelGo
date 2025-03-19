<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

useHead({
  title: "⚙️ Settings"
})
</script>

<template>
  <Container>
    <ReadingPreferencesForm/>
  </Container>
</template>
