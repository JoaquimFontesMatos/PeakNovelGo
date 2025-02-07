import type { Config } from 'tailwindcss';

const inlineConfig: Config = {
    darkMode: 'class', // Enables dark mode using a CSS class (e.g., "dark")
    content: [
        './components/**/*.{vue,js,jsx,mjs,ts,tsx}',
        './layouts/**/*.{vue,js,jsx,mjs,ts,tsx}',
        './pages/**/*.{vue,js,jsx,mjs,ts,tsx}',
        './plugins/**/*.{js,ts,mjs}',
        './composables/**/*.{js,ts,mjs}',
        './utils/**/*.{js,ts,mjs}',
        './{A,a}pp.{vue,js,jsx,mjs,ts,tsx}',
        './{E,e}rror.{vue,js,jsx,mjs,ts,tsx}',
        './app.config.{js,ts,mjs}',
    ],
    theme: {
        extend: {
            fontFamily: {
                montserrat: ['Montserrat'],
                noto: ['Noto Sans'],
                raleway: ['Raleway'],
            },
            colors: {
                primary: 'var(--color-primary)',
                secondary: 'var(--color-secondary)',
                'primary-content': 'var(--color-primary-content)',
                'secondary-content': 'var(--color-secondary-content)',
                'highlight-content': 'var(--color-highlight-content)',
                accent: {
                    gold: 'var(--color-accent-gold)',
                    'gold-light': 'var(--color-accent-gold-light)',
                    'gold-dark': 'var(--color-accent-gold-dark)',
                },
                error: 'var(--color-error)',
                'error-darker': 'var(--color-error-darker)',
                success: 'var(--color-success)',
                warning: 'var(--color-warning)',
                info: 'var(--color-info)',
                neutral: {
                    light: 'var(--color-neutral-light)',
                    dark: 'var(--color-neutral-dark)',
                },
                border: 'var(--color-border)',
            },
        },
    },
    plugins: [require('@tailwindcss/forms')],
};

export default inlineConfig satisfies Config;
