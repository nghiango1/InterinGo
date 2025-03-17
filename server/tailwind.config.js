const plugin = require('tailwindcss/plugin')

module.exports = {
  content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
  theme: {
    extend: {
      gridTemplateColumns: {
        'docs': '19rem 1fr',
      },
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
    plugin(function({ addUtilities, addComponents, e, prefix, config }) {
      const newUtilities = {
        '.horizontal-tb': {
          writingMode: 'horizontal-tb',
        },
        '.vertical-rl': {
          writingMode: 'vertical-rl',
          '-webkit-writing-mode': 'vertical-rl',
          '-ms-writing-mode': 'vertical-rl',
        },
        '.vertical-lr': {
          writingMode: 'vertical-lr',
          '-webkit-writing-mode': 'vertical-lr',
          '-ms-writing-mode': 'vertical-lr',
        }
      }
      addUtilities(newUtilities)
    }),
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
