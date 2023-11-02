/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.html", "./static/*.js"],
  theme: {
    extend: {
      colors: {
        'text': '#ffffff',
        'background': '#140004',
        'primary': '#016b57',
        'secondary': '#002638',
        'accent': '#06fe55',
      },
    },
  },
  plugins: [],
};
