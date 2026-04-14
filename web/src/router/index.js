import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue'),
        meta: { guest: true },
    },
    {
        path: '/',
        name: 'Main',
        component: () => import('../views/Main.vue'),
        meta: { auth: true },
    },
    {
        path: '/history',
        name: 'History',
        component: () => import('../views/History.vue'),
        meta: { auth: true },
    },
    {
        path: '/admin',
        name: 'Admin',
        component: () => import('../views/Admin.vue'),
        meta: { auth: true, admin: true },
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach(async (to, from, next) => {
    const auth = useAuthStore()

    if (!auth.initialized) {
        await auth.ensureAuthLoaded()
    }

    if (to.meta.auth && !auth.isLoggedIn) {
        next('/login')
    } else if (to.meta.guest && auth.isLoggedIn) {
        next('/')
    } else if (to.meta.admin && auth.user?.role !== 'admin') {
        next('/')
    } else {
        next()
    }
})

export default router
