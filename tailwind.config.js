/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/templates/*/*.templ",
    "./web/static/*.js"
  ],
  theme: {
    extend: {
      colors: {
        text: '#ffffff',
        background: '#140004',
        backgroundalt: '#140E20',
        primary: '#016b57',
        subsidary: '#345367',
        secondary: '#002638',
        accent: '#06fe55',
        highlight: '#f9f871',
        warning: '#ea9a27',
        info: '#cef6ff',
      },
      letterSpacing: {
        extreme: '0.2em',
      },
      fontFamily: {
        redhat: "'Red Hat Display', sanf-serif",
        slab: "'Roboto Slab', serif",
      },
      fontSize: {
        '2xs': '.625rem',
        '3xs': '.5rem',
        giant: '20rem',
      },
      borderRadius: {
        "inner": "0.39rem",
      },
    },
  },
  plugins: [],
};
