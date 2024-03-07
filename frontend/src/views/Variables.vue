<script>
import InstallationStep from '../components/InstallationStep.vue';
import { ReadEnvVariables, GetIp, IsPortAvailable, GetAvailablePort } from '../../wailsjs/go/main/App'
import LoadingSpinner from '../components/LoadingSpinner.vue';
const steps = [
	{ title: 'Platform', active: false },
	{ title: 'Environment', active: false },
	{ title: 'Variables', active: true },
	{ title: 'Install', active: false }
];

const tips = 'Set the environment variables and other settings for the installation. You will be able to see and change these values later in the Installed Environments section in the Home. <br><br>If you are not sure about the values, you can leave the default ones.';

export default {
	components: {
		InstallationStep,
		LoadingSpinner
	},
	data() {
		return {
			steps,
			tips,
			// A map of the variables that are ports with the variable name as key and if the port is not ok as value
			portsNotOk: new Map(),
			isLoadingVariables: false
		};
	},
	computed: {
		navigation() {
			return {
				next: {
					path: '/install',
					disabled: !this.areVariablesOk(),
				},
				back: {
					// path: '/environment-setup',
					path: this.$store.state.editingEnvironment ? `/environments/${this.$store.state.installationState.id}` : '/environment-setup',
					disabled: false
				},
				cancel: {
					// /environments if editing, / if creating
					path: this.$store.state.editingEnvironment ? `/environments/${this.$store.state.installationState.id}` : '/',
					disabled: false
				},
			}
		},
		// Map the variables from the store
		variables: {
			get() {
				return this.$store.state.installationState.variables;
			},
			set(value) {
				this.$store.commit('setVariables', value);
			}
		},
		skipImagesAutoupdate: {
			get() {
				return this.$store.state.installationState.skipImagesAutoupdate;
			},
			set(value) {
				this.$store.commit('setSkipImagesAutoupdate', value);
			}
		},
		platform() {
			return this.$store.state.installationState.platform;
		},
	},
	methods: {
		readEnvFileAndShowVariables() {
			let platform = this.$store.state.installationState.platform;
			// Read the env file and handle the promise
			ReadEnvVariables(platform).then(sections => {
				this.variables = sections;
				// Get the IP address and handle the promise
				GetIp().then(ip => {
					// Replace every variable that has ${API_HOST_ENV} as value with the IP address
					this.variables.forEach(section => {
						// Create a new empty object for the new variables
						let updatedVariables = {};
						for (let key in section.variables) {
							let originalValue = section.variables[key];

							// Replace '${API_HOST_ENV}' in the original value with the IP address
							let updatedValue = originalValue.replace('${API_HOST_ENV}', ip);

							updatedVariables[key] = updatedValue;
						}
						section.variables = updatedVariables;
					});

					let portPromises = [];
					// Populate the ports map
					for (let section of this.variables) {
						for (let key in section.variables) {
							portPromises.push(this.isPortOk(key, section.variables[key]));
						}
					}

					// Wait for the ports to be checked
					Promise.all(portPromises).then(() => {
						// For all the ports that are not ok, get an available port suggestion
						this.portsNotOk.forEach((value, key) => {
							if (value) {
								GetAvailablePort().then(port => {
									this.variables.forEach(section => {
										// If in the section there's a variable with the key, set the port as the value
										if (section.variables[key]) {
											section.variables[key] = port;
											// Set the port as ok
											this.portsNotOk.set(key, false);
										}
									});
								}).catch(() => {
									// If there's an error, set the port as not ok
									this.portsNotOk.set(key, true);
								});;
							}
						});
						// Wait 1 second before hiding the spinner (just to make sure it's visible for a short time)
						setTimeout(() => {
							this.isLoadingVariables = false;
						}, 1000);
					});

				});
			});
		},
		// Key is the variable name and value is the value of the variable
		isPortOk(key, value) {
			// If editing, check the port only if is different from the original value
			if (this.$store.state.editingEnvironment) {
				// Find the original value
				for (let section of this.$store.state.editingEnvironment.variables) {
					// If the value is the same, return a resolved promise
					if (section.variables[key] === value) return Promise.resolve();
				}
			}

			// If the key finishes with '_PORT' but skips the POSTGRESQL_PORT
			if (key.endsWith('_PORT') && !(key.startsWith('POSTGRESQL'))) {
				// Check if the port is available
				return IsPortAvailable(value).then(available => {
					// TODO: Check if the port is already being used by another variable in this environment
					// If the port is available, set the port as ok
					this.portsNotOk.set(key, !available);
				}).catch(() => {
					this.portsNotOk.set(key, true);
				});
			}
			// If the key doesn't finish with '_PORT', return a resolved promise
			return Promise.resolve();
		},
		// Check if all the variables are ok
		areVariablesOk() {
			// Check if all the ports are ok
			let allPortsOk = true;
			this.portsNotOk.forEach((value, key) => {
				if (value) {
					allPortsOk = false;
				}
			});
			if (!allPortsOk) {
				return false;
			}
			// TODO: Check if all the variables are ok?

			return true;
		},
	},
	created() {
		if (!this.$store.state.installationState.variables) {
			this.isLoadingVariables = true;
			// Call the method to read the env file and show the variables
			this.readEnvFileAndShowVariables();
		}
	}
};
</script>


<template>
	<LoadingSpinner :isLoading="isLoadingVariables" text="Loading environment variables" />
	<InstallationStep :steps="steps" :tips="tips" :navigation="navigation">
		<!-- The installation step's main content -->
		<div class='variables-main-content-container'>
			<!-- The title-->
			<h1 class="variables-title">Environment variables</h1>
			<!-- The form container -->
			<form class="variables-main-content" autocomplete="off">
				<!-- Checkbox skip autoupdate docker images -->
				<div class="variables-section-container">
					<label class="variables-label">Skip autoupdate docker images</label>
					<div class="variables-checkbox-container">
						<input type="checkbox" class="variables-checkbox" v-model="skipImagesAutoupdate" />
					</div>
				</div>
				<!-- The variables sections -->
				<div v-for="(section, sectionIndex) in variables" :key="sectionIndex">
					<h2 class="variables-section-title">{{ section.name }}</h2>
					<div v-for="(value, key) in section.variables" :key="key" class="variables-section-container">
						<label class="variables-label">{{ key }}</label>
						<input v-model="section.variables[key]" type="text" class="variables-input"
							@input="isPortOk(key, section.variables[key])" />
						<span v-if="portsNotOk.get(key)" class="tooltiptext">This port is not available, please choose
							another port</span>
					</div>
				</div>
			</form>
		</div>
	</InstallationStep>
</template>