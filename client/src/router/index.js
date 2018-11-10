import Vue from 'vue'
import Router from 'vue-router'
import Page1 from '@/components/Page1'
import Page2 from '@/components/Page2'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Page1',
      component: Page1
    },
    {
      path: '/alter',
      name: 'Page2',
      component: Page2
    }
  ]
})
