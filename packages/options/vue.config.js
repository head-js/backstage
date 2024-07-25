module.exports = () => {
  const commonConfig = {
    pages: {
      umi: {
        entry: 'src/index.js',
        template: 'src/document.ejs',
        filename: 'options.html',
      },
    },

    configureWebpack: {
      output: {
        libraryTarget: 'system',
      },

      module: {
        rules: [
          { parser: { system: false } },
        ],
      },
    },
  };

  let config = {};
  if (process.env.NODE_ENV === 'production') {
    config = require('./vue.config.production');
  } else {
    config = require('./vue.config.development');
  }

  return { ...commonConfig, ...config };
}
