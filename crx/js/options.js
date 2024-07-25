// System.addImportMap({
//   imports: {
//   },
// });
// app-*
System.addImportMap({
  imports: {
    'app-account': '/app-account/rsrc/dist/umi.js',
    'app-options': '/app-options/rsrc/dist/umi.js',
  },
});
// entry
System.import('/rsrc/dist/options-vendors-vue.js');
System.import('/rsrc/dist/options-umi.js');
