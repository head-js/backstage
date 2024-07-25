const ElementImporter = require('unplugin-vue-components/webpack');
const { ElementUiResolver } = require('unplugin-vue-components/resolvers');


module.exports = {
  outputDir: '.dist/rsrc/dist',
  publicPath: '/rsrc/dist',

  chainWebpack: (config) => {
    config.output
      .filename('options-[name].js')
      .chunkFilename('options-[name].js');

    config.optimization.splitChunks({
      cacheGroups: {
        vue: {
          name: 'vendors-vue',
          test: /[\\/]node_modules[\\/](vue|vue-router|vuex)[\\/]/,
          priority: -10,
          chunks: 'initial',
        },
        umi: {
          name: 'vendors-umi',
          test: /[\\/]node_modules[\\/](axios)[\\/]/,
          priority: -11,
          chunks: 'initial',
        },
        default: {},
        vendors: {},
      },
    });

    config.plugin('element-importer').use(ElementImporter({
      resolvers: [
        ElementUiResolver(),
      ],
    }));

    config.plugin('html-umi').tap((args) => {
      const options = args[0];
      options.minify = false;
      options.inject = false;
      return args;
    });

    config.plugin('extract-css').tap((args) => {
      const options = args[0];
      options.filename = 'options-[name].css';
      options.chunkFilename = 'options-[name].css';
      return args;
    });
  },


  productionSourceMap: false,
};
