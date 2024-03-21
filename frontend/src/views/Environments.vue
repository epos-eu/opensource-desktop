<script>
import { GetInstalledEnvironments, DeleteInstalledEnvironment } from '../../wailsjs/go/main/App';
import Dialog from '../components/Dialog.vue';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import LoadingSpinner from '../components/LoadingSpinner.vue';

const tips = "Click on an environment on the left to see its details";

const dialogText = "Are you sure you want to delete this environment?";
const dialogConfirmButton = { show: true, text: 'Delete', positive: false };
const dialogCancelButton = { show: true, text: 'Cancel', positive: true };

// Environment structure
// Environment {
// 	    platform: string;
// 	    environmentSetup: EnvironmentSetup;
// 	    variables: Section[];
// EnvironmentSetup {
// 	    name: string;
// 	    version: string;
// 	    context: string;
// Section {
// 	    name: string;
// 	    variables: {[key: string]: string};

export default {
	components: {
		Dialog: Dialog,
		LoadingSpinner,
	},
	data() {
		return {
			tips,
			selectedEnvironment: null,
			dockerEnvironments: [],
			kubernetesEnvironments: [],
			showDialog: false,
			dialogText: dialogText,
			dialogConfirmButton: dialogConfirmButton,
			dialogCancelButton: dialogCancelButton,
			isDeleting: false,
			showErrorDialog: false,
			errorDialogConfirmButton: { show: true, text: 'OK', positive: true },
			id: null,
			loadingEnvironments: true,
		};
	},
	methods: {
		// Select an environment
		selectEnvironment(environment) {
			// If the environment is already selected, deselect it
			if (this.selectedEnvironment === environment) {
				this.selectedEnvironment = null;
				return;
			}

			// Else, select the environment
			this.selectedEnvironment = environment;
			this.$refs.rightContainer.scrollTop = 0;

			// Set the id of the selected event to use in the populate route
			this.id = environment.id;
		},
		confirmDelete() {
			this.showDialog = false;
			let name = this.selectedEnvironment.environmentSetup.name;
			let version = this.selectedEnvironment.environmentSetup.version;
			let platform = this.selectedEnvironment.platform;
			let context = this.selectedEnvironment.environmentSetup.context;

			// Show the spinner
			this.isDeleting = true;

			DeleteInstalledEnvironment(platform, name, version, context).then(() => {
				// Remove the environment from the lists
				if (this.selectedEnvironment.platform === "docker") {
					this.dockerEnvironments = this.dockerEnvironments.filter(environment => environment !== this.selectedEnvironment);
				} else if (this.selectedEnvironment.platform === "kubernetes") {
					this.kubernetesEnvironments = this.kubernetesEnvironments.filter(environment => environment !== this.selectedEnvironment);
				}
				this.selectedEnvironment = null;

				// Hide the spinner
				this.isDeleting = false;

			}).catch((error) => {
				// Show a dialog with the error
				this.errorDialogText = "Error deleting the environment: " + error;
				this.showErrorDialog = true;

				// Hide the spinner
				this.isDeleting = false;
			});
		},
		cancelDelete() {
			this.showDialog = false;
		},
		edit() {
			// Start a new installation with the selected environment as the state and skip the first step
			this.$store.commit('editEnvironmentInit', this.selectedEnvironment);

			// this.$router.push('/environment-setup')	// TODO: allow changing the name/version?
			this.$router.push('/variables')

		},
		openDataPortal() {
			BrowserOpenURL(this.selectedEnvironment.accessPoints.dataPortal);
		},
		openApiGateway() {
			BrowserOpenURL(this.selectedEnvironment.accessPoints.apiGateway);
		},
		closeErrorDialog() {
			this.showErrorDialog = false;
		},
	},
	created() {
		// Get the installed environments and handle the promise
		GetInstalledEnvironments().then(environments => {
			// If there are no environments (null), return
			if (!environments) {
				// Hide the loading spinner after at least 1 second
				setTimeout(() => {
					this.loadingEnvironments = false;
				}, 1000);
				return;
			}

			// Filter the environments and add them to the correct list
			environments.forEach(environment => {
				// Generate the id used by the router
				environment.id = environment.environmentSetup.name + '-' + environment.environmentSetup.version + '-' + environment.platform;
				if (environment.platform === "docker") {
					this.dockerEnvironments.push(environment);
				} else if (environment.platform === "kubernetes") {
					this.kubernetesEnvironments.push(environment);
				}
			});

			// Select the environment if it's in the URL (name-version)
			if (this.$route.params.id) {
				// Split the id into name, version and platform
				let name = this.$route.params.id.split('-')[0];
				let version = this.$route.params.id.split('-')[1];
				let platform = this.$route.params.id.split('-')[2];

				// Search the docker and kubernetes environments for the environment with the same name and version
				let environment = this.dockerEnvironments.find(environment => environment.environmentSetup.name === name && environment.environmentSetup.version === version && environment.platform === platform);
				if (!environment) {
					environment = this.kubernetesEnvironments.find(environment => environment.environmentSetup.name === name && environment.environmentSetup.version === version && environment.platform === platform);
				}

				// If the environment was found, select it
				if (environment) {
					this.selectEnvironment(environment);
				}
			}

			// Hide the loading spinner after at least 1 second
			setTimeout(() => {
				this.loadingEnvironments = false;
			}, 1000);
		});
	},
};
</script>

<template>
	<!-- Loading spinner while deleting -->
	<LoadingSpinner :isLoading="isDeleting" :text="'Deleting environment...'"></LoadingSpinner>
	<!-- Loading spinner while loading the environments -->
	<LoadingSpinner :isLoading="loadingEnvironments" :text="'Loading installed environments...'"></LoadingSpinner>
	<!-- Confirm delete dialog -->
	<Dialog v-if="showDialog" @confirm="confirmDelete" @cancel="cancelDelete" :text="dialogText"
		:confirmButton="dialogConfirmButton" :cancelButton="dialogCancelButton" :title="'Delete environment'"></Dialog>
	<!-- Error deleting dialog -->
	<Dialog v-if="showErrorDialog" @confirm="closeErrorDialog" :text="errorDialogText"
		:title="'Error deleting environment'" :confirmButton="errorDialogConfirmButton" :cancelButton="null"></Dialog>
	<div class="environments-container-container">
		<div class="environments-container">
			<!-- the left list -->
			<div class="environments-left-container">
				<!-- the logo -->
				<img src="../assets/images/epos_logo.svg" alt="EPOS Data Portal" class="home-logo" />
				<!-- show a list of the environments for docker and one for kubernetes -->
				<div class="environments-list">
					<div class="environments-list-title">Docker Environments</div>
					<div class="environments-list-container">
						<div v-for="environment in dockerEnvironments" :key="environment.environmentSetup.name"
							class="environment-item-container"
							:class="{ 'selected-environment-item-container': environment === selectedEnvironment }"
							@click="selectEnvironment(environment)">
							<div class="environment-item">
								<div class="environment-value">{{ environment.environmentSetup.name }}</div>
								<div class="environment-value">{{ "V. " + environment.environmentSetup.version }}</div>
							</div>
						</div>
					</div>
				</div>
				<div class="environments-list">
					<div class="environments-list-title">Kubernetes Environments</div>
					<div class="environments-list-container">
						<div v-for="environment in kubernetesEnvironments" :key="environment.environmentSetup.name"
							class="environment-item-container"
							:class="{ 'selected-environment-item-container': environment === selectedEnvironment }"
							@click="selectEnvironment(environment)">
							<div class="environment-item">
								<div class="environment-value">{{ environment.environmentSetup.name }}</div>
								<div class="environment-value">{{ "V. " + environment.environmentSetup.version }}</div>
							</div>
						</div>
					</div>
				</div>

			</div>
			<!-- the right content -->
			<div class="environments-right-container" ref="rightContainer">
				<div class="environments-right-container-main-content">
					<div v-if="!selectedEnvironment" class="environments-info">
						<div class="environments-title">Installed Environments</div>
						<div class="environments-tips">{{ tips }}</div>
					</div>
					<div v-if="selectedEnvironment" class="environments-details">
						<div class="environments-details-title">
							<span>{{ selectedEnvironment.environmentSetup.name }}</span>
							<span class="environments-title-space">V.</span>
							<span>{{ selectedEnvironment.environmentSetup.version }}</span>
						</div>
						<table class="environments-details-table">
							<tr>
								<td>Name:</td>
								<td>{{ selectedEnvironment.environmentSetup.name }}</td>
							</tr>
							<tr>
								<td>Version:</td>
								<td>{{ selectedEnvironment.environmentSetup.version }}</td>
							</tr>
							<tr>
								<td>Platform:</td>
								<td>{{ selectedEnvironment.platform }}</td>
							</tr>
							<tr v-if="selectedEnvironment.platform == 'kubernetes'">
								<td>Context:</td>
								<td>{{ selectedEnvironment.environmentSetup.context }}</td>
							</tr>
						</table>
						<div class="environments-details-title">Access Points</div>
						<table class="environments-details-table">
							<tr>
								<td>Data Portal:</td>
								<td>{{ selectedEnvironment.accessPoints.dataPortal }}</td>
								<button class="environments-access-points-button
" @click="openDataPortal">Open in the browser</button>
							</tr>
							<tr>
								<td>API Gateway:</td>
								<td>{{ selectedEnvironment.accessPoints.apiGateway }}</td>
								<button class="environments-access-points-button" @click="openApiGateway">Open in the
									browser</button>
							</tr>
						</table>

						<router-link v-if="selectedEnvironment" class="environments-populate-button"
							:to="'/populate/' + id">Populate
							Environment</router-link>
					</div>

					<div v-if="selectedEnvironment" class="environments-details">
						<!-- the details including the variables using a table-->
						<div class="environments-details-title">Environmental Variables</div>
						<div v-for="(section, sectionIndex) in selectedEnvironment.variables" :key="sectionIndex">
							<h2 class="variables-section-title">{{ section.name }}</h2>
							<table class="environments-details-table">
								<tr v-for="(value, key) in section.variables" :key="key">
									<td class="variables-label">{{ key }}</td>
									<td class="variables-value">{{ value }}</td>
								</tr>
							</table>
						</div>
					</div>
				</div>
				<div class="navigation-container">
					<router-link to="/">
						<button class="primary-button">Home</button>
					</router-link>
					<button v-if="selectedEnvironment" class="secondary-button" @click="showDialog = true">Delete
						environment</button>
					<button v-if="selectedEnvironment" class="primary-button" @click="edit">Edit</button>
				</div>
			</div>
		</div>
	</div>
</template>
