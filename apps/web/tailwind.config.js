/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        forge: {
          bg: "#0B0F17",
          card: "#111827",
          border: "#1F2937",
          accent: "#F85D7F",
          gold: "#F8D866",
          cyan: "#38BDF8",
        },
      },
    },
  },
  plugins: [],
};
