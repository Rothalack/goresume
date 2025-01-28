import Vue from 'vue';
import ExampleTest from './components/ExampleTest.vue';

Vue.component('exampleTest', require('./components/ExampleTest.vue').default);

// Initialize Vue instance
new Vue({
  el: '#app', // The element in your HTML to mount Vue
});