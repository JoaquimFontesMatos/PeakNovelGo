// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  experimental: {
    typedPages: true,
  },
  alias: {
    '@img': '/assets/img/',
  },
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@formkit/auto-animate/nuxt',
    [
      '@vee-validate/nuxt',
      {
        // disable or enable auto imports
        autoImports: true,
      },
    ],
    '@nuxt/icon',
    '@nuxt/fonts',
  ],
  fonts: {
    defaults: {
      weights: [400],
      styles: ['normal', 'italic'],
    },
    families: [
      { name: 'Montserratt' },
      { name: 'Noto Sans' },
      { name: 'Raleway' },
    ],
  },
  runtimeConfig: {
    // Keys within public are also exposed client-side
    public: {
      apiUrl: process.env.API_URL,
    },
  },
});