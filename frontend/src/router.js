import { createRouter, createWebHashHistory } from 'vue-router'
import Home from './views/Home.vue'
import { useStore } from './store'

const routes = [
	{
		path: '/',
		name: 'Home',
		component: Home,
	},
	{
		path: '/platform',
		name: 'platform',
		component: () => import('./views/Platform.vue'),
	},
	{
		path: '/environments/:id',
		name: 'EnvironmentsWithId',
		component: () => import('./views/Environments.vue'),
	},
	{
		path: '/environments',
		name: 'Environments',
		component: () => import('./views/Environments.vue'),
	},
	{
		path: '/environment-setup',
		name: 'EnvironmentSetup',
		component: () => import('./views/EnvironmentSetup.vue'),
	},
	{
		path: '/variables',
		name: 'Variables',
		component: () => import('./views/Variables.vue')
	},
	{
		path: '/install',
		name: 'Install',
		component: () => import('./views/Install.vue')
	},
	{
		path: '/about',
		name: 'About',
		component: () => import('./views/About.vue')
	},
	{
		path: '/populate/:id',
		name: 'Populate',
		component: () => import('./views/Populate.vue')
	},
	{
		path: '/populate-install/:id',
		name: 'PopulateInstall',
		component: () => import('./views/PopulateInstall.vue')
	}
]

const router = createRouter({
	history: createWebHashHistory(),
	routes,
})

router.beforeEach((to, from, next) => {
	// when going to the home reset the insallation state
	if (to.path === '/' || to.path === '/environments') {
		useStore().commit('resetInstallationState');
	}
	next();
})

export default router