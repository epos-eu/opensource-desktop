<script>
import InstallationStep from '../components/InstallationStep.vue';
import { IsEnvironmentInstalled, GetKubernetesContexts } from '../../wailsjs/go/main/App'
import { useStore } from '../store'
import LoadingSpinner from '../components/LoadingSpinner.vue'

const steps = [
	{ title: 'Platform', active: false },
	{ title: 'Environment', active: true },
	{ title: 'Variables', active: false },
	{ title: 'Install', active: false }
];

const tips = 'Choose a name and version for the environment. The name and version will be used to identify the environment so they have to be unique. <br><br>You will find the environment by its name and version in the Installed Environments section in the Home.';

export default {
	components: {
		InstallationStep,
		LoadingSpinner
	},
	data() {
		return {
			steps,
			tips,
			warning: false,
			contexts: null,
			loadingContexts: false,
			errorLoadingContexts: false,
		};
	},
	computed: {
		platform: {
			get() {
				return this.$store.state.installationState.platform;
			},
			set(value) {
				this.$store.commit('setPlatform', value);
			}
		},
		name: {
			get() {
				return this.$store.state.installationState.environmentSetup.name;
			},
			set(value) {
				this.$store.commit('setEnvironmentName', value);
			}
		},
		version: {
			get() {
				return this.$store.state.installationState.environmentSetup.version;
			},
			set(value) {
				this.$store.commit('setEnvironmentVersion', value);
			}
		},
		context: {
			get() {
				return this.$store.state.installationState.environmentSetup.context;
			},
			set(value) {
				this.$store.commit('setEnvironmentContext', value);
			}
		},
		navigation() {
			return {
				next: {
					path: '/variables',
					disabled: this.name === "" || this.version === "" || this.name === null || this.version === null || this.warning,
					onClick: this.nextButtonClick
				},
				back: {
					// /environments/:id if editing, /platform if creating
					path: this.$store.state.editingEnvironment ? `/environments/${this.$store.state.installationState.id}` : '/platform',
					disabled: false
				},
				cancel: {
					// /environments/:id if editing, / if creating
					path: this.$store.state.editingEnvironment ? `/environments/${this.$store.state.installationState.id}` : '/',
					disabled: false
				}
			};
		}
	},
	methods: {
		nextButtonClick() {
			// Generate the id
			useStore().state.installationState.id = useStore().state.installationState.environmentSetup.name + '-' + useStore().state.installationState.environmentSetup.version + '-' + useStore().state.installationState.platform;

			// Redirect to the 'next' path
			// (in the navigation component if the 'onClick' function is defined, it will not redirect to the 'next' path automatically, so we have to do it manually here)
			this.$router.push(this.navigation.next.path);
		},
		validateAndSetName(input) {
			let cleanedInput;
			// If the platform is kubernetes, allow only lowercase letters, numbers, and hyphens
			if (this.platform === 'kubernetes') {
				// Allow only the characters that can be used in a file name
				cleanedInput = input.target.value.replace(/[^a-z0-9-]/g, '');
			} else {
				// Also allow uppercase letters
				cleanedInput = input.target.value.replace(/[^a-zA-Z0-9-]/g, '');
			}
			input.target.value = cleanedInput;
			this.name = cleanedInput;

			// Check if the environment already exists
			this.isEnvironmentInstalled();

		},
		validateAndSetVersion(input) {
			// Permit only numbers and dots, can't start with a dot
			let cleanedInput = input.target.value.replace(/[^0-9.]/g, '');
			input.target.value = cleanedInput;
			this.version = cleanedInput;

			// Check if the environment already exists
			this.isEnvironmentInstalled();
		},
		isEnvironmentInstalled() {
			// If editing, only check if the environment already exists if the name or version is different from the original
			if (this.$store.state.editingEnvironment) {
				// Get the original name and version
				let originalName = this.$store.state.editingEnvironment.environmentSetup.name;
				let originalVersion = this.$store.state.editingEnvironment.environmentSetup.version;

				// If the name or version is not changed, skip the check
				if (this.name === originalName && this.version === originalVersion) return
			}
			// Check if the environment already exists
			IsEnvironmentInstalled(this.name, this.version, this.platform, this.context).then((result) => {
				this.warning = result;
			});
		},
		loadContexts() {
			// Set the loading spinner
			this.loadingContexts = true;
			// Get the contexts
			GetKubernetesContexts().then((result) => {
				this.contexts = result;

				// Set the first context as the default
				this.context = this.contexts[0];
				// Wait at least 1 second before removing the loading spinner
				setTimeout(() => {
					this.loadingContexts = false;
				}, 1000);
			}).catch((error) => {
				console.error(error);
				setTimeout(() => {
					// Remove the loading spinner
					this.loadingContexts = false;

					// Show an error message
					this.errorLoadingContexts = true;

					// Disable the next button
					this.navigation.next.disabled = true;
				}, 1000);
			});
		}
	},
	mounted() {
		// Check if the environment already exists
		this.isEnvironmentInstalled();

		// If the platform is kubernetes, get the contexts
		if (this.platform === 'kubernetes') {
			this.loadContexts();
		}
	}
};
</script>


<template>
	<LoadingSpinner :isLoading="loadingContexts" text="Loading Kubernetes Contexts" />
	<InstallationStep :steps="steps" :tips="tips" :navigation="navigation">
		<!-- The installation step's main content -->
		<div class='environment-setup-main-content-container'>
			<!-- The title-->
			<h1 class="environment-setup-title">Environment setup</h1>
			<!-- The form container -->
			<div class="environment-setup-main-content">
				<form autocomplete="off">
					<div class="environment-setup-form-field">
						<label for="name">Name:</label>
						<input type="text" id="name" name="name" maxlength="40" @input="validateAndSetName"
							v-model="name">
					</div>
					<div class="environment-setup-form-field">
						<label for="version">Version:</label>
						<input type="text" id="version" name="version" maxlength="20" @input="validateAndSetVersion"
							v-model="version">
					</div>
					<!-- show only if kubernetes is the chosen platform  -->
					<div v-if="platform === 'kubernetes'" class="environment-setup-form-field">
						<label for="context">Context:</label>
						<select id="context" name="context" v-model="context" class="environment-setup-select">
							<option disabled value="">Please select one</option>
							<option v-for="ctx in contexts" :value="ctx" :key="ctx">
								{{ ctx }}
							</option>
						</select>
					</div>
				</form>
				<p v-if="warning" class="environment-setup-warning">
					An environment with the same name and version already exists. Please choose a different name or
					version.
				</p>
				<p v-if="errorLoadingContexts" class="environment-setup-warning">
					Error loading the Kubernetes contexts.
					<button @click="loadContexts" class="primary-button">
						Try again
					</button>
				</p>
			</div>
		</div>
	</InstallationStep>
</template>