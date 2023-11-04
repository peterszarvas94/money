/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.html", "./static/*.js"],
  theme: {
    extend: {
      colors: {
        text: '#ffffff',
        background: '#140004',
        primary: '#016b57',
        secondary: '#002638',
        accent: '#06fe55',
        highlight: '#f9f871',
        warning: '#ea9a27',
        info: '#cef6ff',
        subsidary: '#345367',
      },
      letterSpacing: {
        extreme: '0.2em',
      },
      fontFamily: {
        redhat: "'Red Hat Display', sanf-serif",
        slab: "'Roboto Slab', serif",
      },
      fontSize: {
        giant: '20rem',
      },
    },
  },
  plugins: [],
};
