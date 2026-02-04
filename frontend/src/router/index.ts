import { createRouter, createWebHashHistory } from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'
import ScriptsPage from '../pages/ScriptsPage.vue'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [{ path: '', component: ScriptsPage }],
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
