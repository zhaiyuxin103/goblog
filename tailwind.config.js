/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './resources/views/**/*.tmpl',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}
