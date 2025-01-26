// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  alias: {
    "@img": "/assets/img/",
  },
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  compatibilityDate: "2024-11-01",
  devtools: { enabled: true },
  modules: [
    "@nuxtjs/tailwindcss",
    "@pinia/nuxt",
    "@formkit/auto-animate/nuxt",
    [
      "@vee-validate/nuxt",
      {
        // disable or enable auto imports
        autoImports: true,
      },
    ],
    "@nuxt/icon",
  ],
  runtimeConfig: {
    // Keys within public are also exposed client-side
    public: {
      apiUrl: process.env.API_URL,
    },
  },
});