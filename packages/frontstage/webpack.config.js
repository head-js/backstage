const path = require('path');
const Config = require('webpack-chain');
const ModuleFederationPlugin = require('webpack/lib/container/ModuleFederationPlugin');
const BannerPlugin = require('webpack/lib/BannerPlugin');


const services = [
  { name: 'www.example.com', version: '0.1.0', match: 'http://127.0.0.1:8000/*' },
];


const config = new Config();

config.mode('production');


services.forEach((svc) => {
  config.entry(svc.name)
  .add(`./src/${svc.name}/frontstage.js`)
  .end();
});

config.output
  .path(path.resolve('.dist'))
  .filename('[name]/frontstage.js')
  .publicPath('chrome://crx/js/');

config.output.set('iife', true);

config.optimization.set('chunkIds', 'named');

// config.optimization.set('splitChunks', {
//   chunks: 'async',
//   // minSize: 0,
//   cacheGroups: {
//     default: false,
//     defaultVendors: false,
//     react: {
//       name: 'vendors-react',
//       chunks: 'initial',
//       minChunks: 1,
//       test: /[\\/]node_modules[\\/](axios|react|react-dom)[\\/]/,
//       priority: -10,
//     },
//   }
// });

config.optimization.set('minimize', false);

config.plugin('module-federation').use(ModuleFederationPlugin, [{
  remotes: {
    'backstage': '__backstagevendors__@chrome-extension://elagegodfhfilllnhpmnbmeokdimoeda/vendors/backstage-vendors.js',
  },
}]);


function makeBanner(name) {
  const svc = services.find((s) => s.name === name);

  return `// ==UserScript==
// @name        frontstage - ${svc.name}
// @namespace   https://lisite.de
// @updateURL   TODO
// @downloadURL TODO
// @version     ${svc.version}
// @description Frontstage
// @author      hbrls
// @match       ${svc.match}
// @icon        data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEKLAAAAAABAAEAAAICTAEA0w==
// @grant       none
// @run-at      document-idle
// @sandbox     raw
// @unwrap
// ==/UserScript==`;
}


config.plugin('banner')
  .use(BannerPlugin, [ { banner: ({ chunk: { name } }) => makeBanner(name), raw: true, entryOnly: true } ]);

const conf = config.toConfig();
// console.log(conf);

module.exports = conf;
