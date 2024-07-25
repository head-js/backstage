import { registerApplication, start } from 'single-spa';
import { constructRoutes, constructApplications, constructLayoutEngine } from 'single-spa-layout';
import Vue from 'vue';
import router from './router';
import AsideMenu from './components/AsideMenu';
import 'element-ui/lib/theme-chalk/index.css';


const routes = constructRoutes(`
<single-spa-router mode="hash" containerEl="#root">
  <route path="/account"><application name="app-account"></application></route>
  <route path="/options"><application name="app-options"></application></route>
  <router default><h1 class="text-center">404 Not Found</h1></router>
</single-spa-router>
`);


const applications = constructApplications({
  routes,
  loadApp({ name }) { return System.import(name); }, // eslint-disable-line no-undef
});


const layoutEngine = constructLayoutEngine({ routes, applications });


applications.forEach(registerApplication);


function $ready(fn) {
  fn();
}


$ready(() => {
  Vue.config.productionTip = false;

  new Vue({
    router,
    render: (h) => h(AsideMenu),
  }).$mount('#aside');

  start();
});
