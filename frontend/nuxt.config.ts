// https://nuxt.com/docs/api/configuration/nuxt-config

import tailwindcss from '@tailwindcss/vite';

export default defineNuxtConfig({
    components: [
        {
            path: '~/components',
            pathPrefix: false,
        },
    ],
    experimental: {
        typedPages: true,
    },
    alias: {
        '@img': '/assets/img/',
    },
    css: ['~/assets/css/tailwind.css'],
    compatibilityDate: '2025-02-27',
    devtools: { enabled: true },
    modules: ['@pinia/nuxt', '@formkit/auto-animate/nuxt', '@nuxt/icon', '@nuxt/fonts', '@nuxtjs/color-mode', '@vite-pwa/nuxt'],
    fonts: {
        defaults: {
            weights: [400],
            styles: ['normal', 'italic'],
        },
        families: [{ name: 'Montserratt' }, { name: 'Noto Sans' }, { name: 'Raleway' }],
    },
    runtimeConfig: {
        // Keys within public are also exposed client-side
        public: {
            apiUrl: process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8081',
            runMode: process.env.NUXT_PUBLIC_RUN_MODE || 'development',
        },
    },
    colorMode: {
        classSuffix: '',
        preference: 'dark',
        fallback: 'dark',
    },
    vite: {
        plugins: [tailwindcss()],
    },
    pwa: {
        registerType: 'autoUpdate',
        manifest: {
            name: 'PeakNovelGo',
            short_name: 'PeakNovelGo',
            theme_color: '#ffffff',
            start_url: '/',
            display: 'standalone',
            icons: [
                {
                    src: '/android-chrome-192x192.png',
                    sizes: '192x192',
                    type: 'image/png',
                },
                {
                    src: '/android-chrome-512x512.png',
                    sizes: '512x512',
                    type: 'image/png',
                },
            ],
        },
    },
});
