module.exports = {
  plugins: [
    '@head.js/babel-plugin-console',
  ],
  presets: [
    ['@babel/preset-env', {
      useBuiltIns: 'usage',
      corejs: '3',
    }]
  ],
}
