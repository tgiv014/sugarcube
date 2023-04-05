/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,svelte}'
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          'InterVariable',
          'sans-serif',
        ],
        mono: ['IBM Plex Mono',
          'monospace'
        ],
      }
    },
  },
  plugins: [],
}

