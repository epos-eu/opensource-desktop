import { createStore } from 'vuex'

const initialState = {
	id: null,	// just a string {environmentSetup.name}-{environmentSetup.version} only used in the frontend
	platform: null,
	environmentSetup: {
		name: null,
		version: '0.0',
	},
	// variables is an array of Section objects
	// class Section {
	//     name: string;
	//     variables: {[key: string]: string};
	variables: null,
	skipImagesAutoupdate: false,
};

// Create and export the Vuex store
export const store = createStore({
	state() {
		return {
			installationState: JSON.parse(JSON.stringify(initialState)),
			editingEnvironment: null,	// the original environment being edited
			checkForUpdatesDone: false,	// flag to indicate if the check for updates has been done
			homeBannerFullscreen: true,	// true when the app starts, false when the user dismisses the banner
			populateState: {
				name: null,
				version: null,
				id: null,
				path: ""
			},	// the state of the populate process
		}
	},
	mutations: {
		setPlatform(state, platform) {
			state.installationState.platform = platform;
			// Invalidate the variables if the platform changes
			state.installationState.variables = null;
		},
		setEnvironment(state, environment) {
			state.installationState.environmentSetup = environment;
		},
		setEnvironmentName(state, name) {
			state.installationState.environmentSetup.name = name;
		},
		setEnvironmentVersion(state, version) {
			state.installationState.environmentSetup.version = version;
		},
		setEnvironmentContext(state, context) {
			state.installationState.environmentSetup.context = context;
		},
		resetInstallationState(state) {
			state.installationState = JSON.parse(JSON.stringify(initialState));
		},
		setVariables(state, variables) {
			state.installationState.variables = variables
		},
		editEnvironmentInit(state, environment) {
			state.installationState = environment;
			// copy the environment being edited
			state.editingEnvironment = JSON.parse(JSON.stringify(environment));
		},
		resetEditingEnvironment(state) {
			state.editingEnvironment = null;
		},
		setSkipImagesAutoupdate(state, skip) {
			state.installationState.skipImagesAutoupdate = skip;
		},
		setCheckForUpdatesDone(state, done) {
			state.checkForUpdatesDone = done;
		},
		setHomeBannerFullscreen(state, fullscreen) {
			state.homeBannerFullscreen = fullscreen;
		},
		setPopulateState(state, populateState) {
			state.populateState = populateState;
		},
		resetPopulateState(state) {
			state.populateState = {
				name: null,
				version: null,
				id: null,
				path: ""
			};
		}
	}
})

// This can be used in components to access the store
export function useStore() {
	return store
}