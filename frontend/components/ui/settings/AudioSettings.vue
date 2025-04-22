<script setup lang="ts">
const userStore = useUserStore();
const { user } = storeToRefs(userStore);
const colorMode = useColorMode();

const handleChangeReadingPreferences = async () => {
  if (!user.value) return;

  const fields = {
    readingPreferences: JSON.stringify(user.value.readingPreferences),
  };

  if (user.value.readingPreferences.theme !== undefined && user.value.readingPreferences.theme) {
    colorMode.preference = user.value.readingPreferences.theme;
  }

  try {
    await userStore.updateUserFields(fields);
  } catch {}

  userStore.saveUserLocalStorage();
};
</script>

<template>
  <div v-if="user" class="menu-container form-container">
    <!-- Text-to-Speech Section -->
    <SectionHeader :title="'Text-to-Speech'" :is-main-header="true">
      <div class="mt-2 space-y-4">
        <!-- Autoplay -->
        <div class="form-group flex items-center space-x-2">
          <input
              id="autoplay"
              name="autoplay"
              type="checkbox"
              v-model="user.readingPreferences.tts.autoplay"
              @change="handleChangeReadingPreferences()"
          />
          <label for="autoplay" class="text-sm font-medium text-secondary-content">Autoplay</label>
        </div>

        <!-- Voice Selector -->
        <div class="form-group">
          <label for="voice" class="block text-sm font-medium text-secondary-content">Voice</label>
          <select id="voice" name="voice" v-model="user.readingPreferences.tts.voice" @change="handleChangeReadingPreferences()">
            <option value="en-US-AriaNeural">en-US-AriaNeural</option>
            <option value="en-US-MichelleNeural">en-US-MichelleNeural</option>
            <option value="en-US-ChristopherNeural">en-US-ChristopherNeural</option>
            <option value="en-US-AnaNeural">en-US-AnaNeural</option>
            <option value="en-GB-SoniaNeural">en-GB-SoniaNeural</option>
            <option value="en-US-AvaNeural">en-US-AvaNeural</option>
            <option value="en-GB-LibbyNeural">en-GB-LibbyNeural</option>
            <option value="es-ES-ElviraNeural">es-ES-ElviraNeural</option>
            <option value="en-US-SteffanNeural">en-US-SteffanNeural</option>
            <option value="en-US-JennyNeural">en-US-JennyNeural</option>
            <option value="en-US-GuyNeural">en-US-GuyNeural</option>
            <option value="en-US-RogerNeural">en-US-RogerNeural</option>
            <option value="en-US-AvaMultilingualNeural">en-US-AvaMultilingualNeural</option>
          </select>
        </div>

        <!-- Rate Input -->
        <div class="form-group">
          <label for="rate" class="block text-sm font-medium text-secondary-content">Rate</label>
          <input
              id="rate"
              name="rate"
              type="range"
              step="5"
              min="-100"
              max="100"
              v-model.number="user.readingPreferences.tts.rate"
              @change="handleChangeReadingPreferences()"
          />
          <span class="block text-sm font-medium text-secondary-content">{{ user.readingPreferences.tts.rate }}%</span>
        </div>
      </div>
    </SectionHeader>
  </div>
</template>
