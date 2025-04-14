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
        <SectionHeader :title="'Reading Preferences'" :is-main-header="true">
            <!-- Atomic Reading -->
            <div class="form-group flex items-center space-x-2">
                <input
                    id="atomicReading"
                    name="atomicReading"
                    type="checkbox"
                    v-model="user.readingPreferences.atomicReading"
                    @change="handleChangeReadingPreferences()"
                />
                <label for="atomicReading" class="text-sm font-medium text-secondary-content">Enable Atomic Reading</label>
            </div>

            <!-- Font Selector -->
            <div class="form-group mt-2">
                <label for="fontPreview" class="block text-sm font-medium text-secondary-content">Font</label>
                <select id="fontPreview" name="font" v-model="user.readingPreferences.font" @change="handleChangeReadingPreferences()">
                    <option value="font-serif">Serif</option>
                    <option value="font-sans">Sans Serif</option>
                    <option value="font-mono">Monospace</option>
                    <option value="font-montserrat">Montserrat</option>
                    <option value="font-noto">Noto Sans</option>
                    <option value="font-raleway">Raleway</option>
                </select>
                <p class="mt-2" :class="user.readingPreferences.font">Preview: The quick brown fox jumps over the lazy dog.</p>
            </div>
        </SectionHeader>

        <!-- Theme Selector -->
        <div class="form-group">
            <label for="theme" class="block text-sm font-medium text-secondary-content">Theme</label>
            <select id="theme" name="theme" v-model="user.readingPreferences.theme" @change="handleChangeReadingPreferences()">
                <option value="dark">dark</option>
                <option value="light">light</option>
                <option value="cyberpunk">cyberpunk</option>
                <option value="forest">forest</option>
                <option value="heaven">heaven</option>
                <option value="deep-blue">deep-blue</option>
                <option value="purple-dusk">purple-dusk</option>
                <option value="crimson-night">crimson-night</option>
                <option value="cyber-green">cyber-green</option>
                <option value="warm-amber">Warm Amber</option>
                <option value="midnight-indigo">midnight-indigo</option>
                <option value="neon-pink">neon-pink</option>
                <option value="emerald-twilight">emerald-twilight</option>
                <option value="smoky-quartz">smoky-quartz</option>
                <option value="obsidian-flame">obsidian-flame</option>
                <option value="velvet-noir">velvet-noir</option>
            </select>
        </div>
    </div>
</template>
