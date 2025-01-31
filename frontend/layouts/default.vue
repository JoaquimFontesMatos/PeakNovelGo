<script setup lang="ts">
onMounted(async() => {
  await useAuthStore().initSession();
});
</script>
<template>
  <div>
    <AppHeader class="fixed top-0 left-0 w-full"/>
    <VerticalSpacer/>
    <div class="min-h-svh"><slot/></div>
    <AppFooter/>
    <Toast/>
  </div>
</template>
