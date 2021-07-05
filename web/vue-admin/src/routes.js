import Vue from 'vue'
import Router from 'vue-router'
Vue.use(Router)
export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/login',
      component: () => import('./views/Login.vue'),
      name: '',
      hidden: true
    },
    {
      path: '/404',
      component: () => import('./views/404.vue'),
      name: '',
      hidden: true
    },
    // {
    //   path: '/',
    //   component: () => import('./views/Home.vue'),
    //   name: '导航一',
    //   iconCls: 'el-icon-message', // 图标样式class
    //   children: [
    //     { path: '/main', component: () => import('./views/Main.vue'), name: '主页', hidden: true },
    //     { path: '/table', component: () => import('./views/nav1/Table.vue'), name: 'Table' },
    //     { path: '/form', component: () => import('./views/nav1/Form.vue'), name: 'Form' },
    //     { path: '/user', component: () => import('./views/nav1/user.vue'), name: '列表' }
    //   ]
    // },
    // {
    //   path: '/',
    //   component: () => import('./views/Home.vue'),
    //   name: '导航二',
    //   iconCls: 'fa fa-id-card-o',
    //   children: [
    //     { path: '/page4', component: () => import('./views/nav2/Page4.vue'), name: '页面4' },
    //     { path: '/page5', component: () => import('./views/nav2/Page5.vue'), name: '页面5' }
    //   ]
    // },
    {
      path: '/',
      component: () => import('./views/Home.vue'),
      name: '用户信息',
      iconCls: 'fa fa-id-card-o',
      leaf: true,
      children: [
        { path: '/userinfo', component: () => import('./views/UserInfo.vue'), name: '用户信息' }
      ]
    },
    {
      path: '/',
      component: () => import('./views/Home.vue'),
      name: '接收信息',
      iconCls: 'fa fa-address-card',
      leaf: true,
      children: [
        { path: '/receiverinfo', component: () => import('./views/ReceiverInfo.vue'), name: '接收信息' }
      ]
    },
    {
      path: '/',
      component: () => import('./views/Home.vue'),
      name: '值班信息',
      iconCls: 'fa fa-address-card',
      leaf: true,
      children: [
        { path: '/userrota', component: () => import('./views/UserRota.vue'), name: '值班信息' }
      ]
    },
    {
      path: '*',
      hidden: true,
      redirect: { path: '/404' }
    }
  ]
})
