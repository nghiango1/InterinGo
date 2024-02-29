const plugin = require('tailwindcss/plugin')

module.exports = {
  content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
  theme: {
    extend: {
      typography: {
        DEFAULT: {
          css: {
            h1: {
              "margin-bottom": "0",
              padding: "0.25rem",
              "border-bottom-width": "2px",
            },
            q: {
              display: "block",
              "margin-bottom": "1.25rem",
              "padding-left": "1rem",
              "border-left-width": "4px",
              "color":"#656d76",
            }
          }
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    plugin(function({ addComponents }) {
      addComponents({
        '.scrollbar::-webkit-scrollbar' : {
          width: '10px',
          height: '10px',
        },

        '.scrollbar::-webkit-scrollbar-track' : {
          'border-radius': '100vh',
          background: '#f7f4ed',
        },

        '.scrollbar::-webkit-scrollbar-thumb' : {
          background: '#e0cbcb',
          'border-radius': '100vh',
          border: '2px solid #f6f7ed',
        },

        '.scrollbar::-webkit-scrollbar-thumb:hover' : {
          background: '#c0a0b9',
        },
      })
    })
  ],
}
