import Vue from 'vue';
import VueRouter, { Route, RouteConfig } from 'vue-router';
import Login from '@/views/Login.vue';
import ELayout from '@/layouts/layout.enum';

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Home',
    component: () => import(/* webpackChunkName: "about" */ '../views/Home.vue'),
    meta: {
    },
  },
  {
    path: '/Login',
    name: 'Login',
    component: Login,
    meta: {
      layout: ELayout.Fullscreen,
      unauthorized: true,
    },
  },
  {
    path: '/joinMember',
    name: 'JoinMember',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/JoinMember.vue'),
    meta: {
      layout: ELayout.Fullscreen,
      unauthorized: true,
    },
  },
];

const router = new VueRouter({
  mode: 'history',
  routes,
});
//
// router.beforeEach((to: Route, from: Route, next: (vm?: string) => void) => {
//   if (!to.matched.some((record) => record.meta.unauthorized)
//     || !Vue.$cookies.get('access-token')
//     || !Vue.$cookies.get('refresh-token')) {
//     console.log(!Vue.$cookies.get('access-token'));
//
//     window.alert('로그인 해주세요.');
//     return next('/Login');
//   }
//   return next();
// });

export default router;
