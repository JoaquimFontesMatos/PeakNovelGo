// tailwind.config.ts
import type { Config } from 'tailwindcss';

const inlineConfig: Config = {
    content: [
        './components/**/*.{vue,js,jsx,mjs,ts,tsx}',
        './components/global/**/*.{vue,js,jsx,mjs,ts,tsx}',
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
                primary: '#1a202c', // equivalent to bg-gray-900
                secondary: '#2d3748', // equivalent to bg-gray-800,
                'primary-content': '#e2e8f0', // equivalent to text-gray-200
                'secondary-content': '#a0aec0', // equivalent to text-gray-400
                'highlight-content': '#fbbf24', // equivalent to text-yellow-400
                accent: {
                    gold: '#fbbf24', // main gold
                    'gold-light': '#facc15', // light gold
                    'gold-dark': '#d97706', // dark gold
                },
                error: '#e53e3e', // red-600
                'error-darker': '#942828',
                success: '#38a169', // green-600
                warning: '#dd6b20', // orange-600
                info: '#3182ce', // blue-600
                neutral: {
                    light: '#f7fafc', // gray-100
                    dark: '#2d3748', // gray-800
                },
                border: '#4a5568', // gray-600
            },
        },
    },
    plugins: [require('@tailwindcss/forms')],
};

export default inlineConfig satisfies Config;