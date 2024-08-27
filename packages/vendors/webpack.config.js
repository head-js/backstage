const path = require('path');
const Config = require('webpack-chain');
const ModuleFederationPlugin = require('webpack/lib/container/ModuleFederationPlugin');


const config = new Config();

config.mode('production');

config.entry('umi')
  .add('./src/index.js')
  .end();

config.output
  .path(path.resolve('../../crx/vendors'))
  .filename('[name].js')
  .publicPath('chrome-extension://elagegodfhfilllnhpmnbmeokdimoeda/vendors/');

config.optimization.set('chunkIds', 'named');

config.optimization.set('minimize', false);

config.plugin('module-federation').use(ModuleFederationPlugin, [{
  name: '__backstagevendors__',
  filename: 'backstage-vendors.js',
  remotes: {},
  exposes: {
    './ajv': './src/ajv',
    './jsonata': './src/jsonata',
    './rxjs': './src/rxjs',
  },
}]);

const conf = config.toConfig();
// console.log(conf);

module.exports = conf;
