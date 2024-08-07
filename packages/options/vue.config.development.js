module.exports = {
  publicPath: '/',

  devServer: {
    host: '0.0.0.0',
    port: 3001,
  },

  chainWebpack: (config) => {
    config.module
      .rule('vue').use('vue-loader').tap((options) => {
        options.prettify = false;
        return options;
      });

    config.plugin('html-umi').tap((args) => {
      const options = args[0];
      options.minify = false;
      options.inject = false;
      return args;
    });
  },

  lintOnSave: false,
}
