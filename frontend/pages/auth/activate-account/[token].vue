<script setup lang="ts">
const route = useRoute();

const authStore = useAuthStore();

const { loadingVerifyToken, verifyTokenMessage } = storeToRefs(authStore);

type MyRouteParams = { token?: string };

const params = route.params as MyRouteParams;

onMounted(async () => {
  try {
    if (!params.token) return;

    await authStore.verifyToken(params.token);
  } catch {}
});
</script>

<template>
  <LoadingBar v-if="loadingVerifyToken" />
  <Container v-else>
    <div v-if="verifyTokenMessage">
      <p>Your account has been activated. You can now log in.</p>
      <p>You can close this window now.</p>
    </div>
    <div v-else>
      <p>An unexpected error occurred. Please try again.</p>
      <p>You can close this window now.</p>
    </div>
  </Container>
</template>
