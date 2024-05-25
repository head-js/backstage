import eslint from '@rollup/plugin-eslint';
import commonjs from '@rollup/plugin-commonjs';
import resolve from '@rollup/plugin-node-resolve';
import babel from '@rollup/plugin-babel';
import json from '@rollup/plugin-json';


export default [
  {
    input: 'src/frontstage.js',

    external: [
      'core-js/modules/es.promise.js',
    ],

    plugins: [
      eslint(),

      json(),

      commonjs({
        sourceMap: false,
      }),

      resolve({
        browser: true,
      }),

      babel({
        exclude: 'node_modules/**',
        babelHelpers: 'bundled',
      }),
    ],

    output: [
      { file: 'crx/js/frontstage.js', format: 'iife', exports: 'none' },
    ],
  },
  {
    input: 'src/backstage.js',

    external: [
      'core-js/modules/es.promise.js',
    ],

    plugins: [
      eslint(),

      json(),

      commonjs({
        sourceMap: false,
      }),

      resolve({
        browser: true,
      }),

      babel({
        exclude: 'node_modules/**',
        babelHelpers: 'bundled',
      }),
    ],

    output: [
      { file: 'crx/js/backstage.js', format: 'iife', exports: 'none' },
    ],
  },
  {
    input: 'src/background.js',

    external: [
      'core-js/modules/es.promise.js',
    ],

    plugins: [
      eslint(),

      json(),

      commonjs({
        sourceMap: false,
      }),

      resolve({
        browser: true,
      }),

      babel({
        exclude: 'node_modules/**',
        babelHelpers: 'bundled',
      }),
    ],

    output: [
      { file: 'crx/js/background.js', format: 'iife', exports: 'none' },
    ],
  },
];
