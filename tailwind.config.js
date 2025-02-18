/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: "jit",
  content: ["./interface/html/**/*.tmpl.html"],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        'accent': '#22201e',
        'accent-dark': '#fafafa',
        'neutral-915': '#151515',
        'neutral-925': '#101010',
      },
      boxShadow: {
        'b': '0px 10px 15px -15px rgba(0, 0, 0, 0.3)',
        'inner-lg': 'inset 0px 0px 6px 0px rgba(0, 0, 0, 0.1)',
      },
      fontFamily: {
        'source-sans': ["'Source Sans 3'", 'sans-serif'],
        'open-sans': ["'Open Sans'", 'sans-serif'],
        'nunito': ["'Nunito Sans'", 'sans-serif'],
        'noto-sans': ["'Noto Sans'", 'sans-serif'],
        'hubba': ["'Hubba'", 'sans-serif'],
        'zing': ["'Zing Rust Demo'", 'sans-serif'],
      },
    },
  },
  plugins: [],
}

