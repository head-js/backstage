module.exports = {
  root: true,

  env: {
    browser: true,
  },

  parserOptions: {
    parser: 'babel-eslint',
  },

  globals: {
    head: true,
  },

  extends: [
    'plugin:eslint-plugin-vue/essential',
    '@vue/eslint-config-airbnb',
  ],

  rules: {
    'import/extensions': [ 'error', { 'js': 'never', 'vue': 'never', 'json': 'always' } ],
    'no-console': 'off',
    'no-multiple-empty-lines': [ 'error', { 'max': 2 } ],
    'no-unused-vars': 'off',
    'object-curly-newline': 'off',
  },
};
