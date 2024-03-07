<script>
import InstallationStep from '../components/InstallationStep.vue';
import { mapState } from 'vuex'
import { IsDockerInstalled, IsDockerRunning, IsKubernetesInstalled, SpecifyPlatformPath } from '../../wailsjs/go/main/App';
import LoadingSpinner from '../components/LoadingSpinner.vue';

const steps = [
	{ title: 'Platform', active: true },
	{ title: 'Environment', active: false },
	{ title: 'Variables', active: false },
	{ title: 'Install', active: false }
];

const tips = 'Select the platform you want to use for the installation. <br><br>If you want to use Docker, make sure it is installed and running.<br>If you want to use Kubernetes, make sure it is installed and properly configured.<br><br>If the installer doesn\'t find the platform you want to use, but you have it installed, you can specify the path to the platform\'s executable.';

export default {
	// what the component provides
	components: {
		InstallationStep,
		LoadingSpinner
	},
	// component data
	data() {
		return {
			steps,
			tips,
			// Using null to indicate that the checks have not been performed yet
			isDockerInstalled: null,
			isDockerRunning: null,
			isKubernetesInstalled: null,
			isChecking: false
		};
	},
	// computed properties
	computed: {
		...mapState({
			selectedPlatform: state => state.installationState.platform
		}),
		navigation() {
			return {
				next: {
					path: '/environment-setup',
					disabled: !this.selectedPlatform
				},
				back: {
					path: '/',
					disabled: false
				},
				cancel: {
					path: '/',
					disabled: false
				}
			};
		}
	},
	// component methods
	methods: {
		selectPlatform(platform) {
			if ((platform === 'docker' && this.isDockerInstalled) || (platform === 'kubernetes' && this.isKubernetesInstalled)) {
				this.$store.commit('setPlatform', platform);
			}
		},
		// TODO: make sure this works
		specifyKubernetesPath() {
			SpecifyPlatformPath('kubernetes').then(path => {
				if (path !== '') {
					this.getInstalledPlatforms();
				}
			}).catch((error) => {
				this.isKubernetesInstalled = false;
			});
		},
		// TODO: make sure this works
		specifyDockerPath() {
			SpecifyPlatformPath('docker').then(path => {
				if (path !== '') {
					this.getInstalledPlatforms();
				}
			}).catch((error) => {
				this.isDockerInstalled = false;
			});
		},
		getInstalledPlatforms() {
			this.isChecking = true;
			// Create promises for the Docker and Kubernetes checks
			const dockerCheck = IsDockerInstalled().then(isInstalled => {
				this.isDockerInstalled = isInstalled;
			});

			const isDockerRunningCheck = IsDockerRunning().then(isRunning => {
				this.isDockerRunning = isRunning;
			});

			const kubernetesCheck = IsKubernetesInstalled().then(isInstalled => {
				this.isKubernetesInstalled = isInstalled;
			});

			// Wait for both checks to finish
			Promise.all([dockerCheck, isDockerRunningCheck, kubernetesCheck]).then(() => {
				// if only one platform is installed, select it
				if (this.isDockerInstalled && !this.isKubernetesInstalled) {
					this.selectPlatform('docker');
				} else if (!this.isDockerInstalled && this.isKubernetesInstalled) {
					this.selectPlatform('kubernetes');
				}
				// Wait 1 second before hiding the spinner (just to make sure it's visible for a short time)
				setTimeout(() => {
					this.isChecking = false;
				}, 1000);
			});
		},
	},
	beforeCreate() {
		// Reset the editing environment state, if i'm here it means i'm creating a new environment
		this.$store.commit('resetEditingEnvironment');
	},
	created() {
		// Get the installed platforms
		this.getInstalledPlatforms();
	}
};
</script>


<template>
	<!-- Loading spinner while checking for the installed platforms -->
	<LoadingSpinner :isLoading="isChecking" text="Checking for installed platforms..." />
	<!-- The installation step stub -->
	<InstallationStep :steps="steps" :tips="tips" :navigation="navigation">
		<!-- The installation step's main content -->
		<!-- A container used to layout the platform selection buttons -->
		<div class='platform-main-content-container'>

			<!-- The title of the platform selection step -->
			<h1 class="platform-title">Choose a platform</h1>

			<!-- The container for the platform selection buttons -->
			<div class="platform-main-content">

				<div class="platform-main-content-buttons-container">
					<!-- Docker -->
					<button class="platform-button"
						:class="{ 'grayed-out': (selectedPlatform && selectedPlatform !== 'docker') || !isDockerInstalled || !isDockerRunning}"
						title="Docker" @click="selectPlatform('docker')" :disabled="!isDockerInstalled || !isDockerRunning">
						<img src="../assets/images/docker-logo.png" alt="Docker logo" class="platform-button-image" />
					</button>

					<!-- Kubernetes -->
					<button class="platform-button"
						:class="{ 'grayed-out': (selectedPlatform && selectedPlatform !== 'kubernetes') || !isKubernetesInstalled }"
						title="Kubernetes" @click="selectPlatform('kubernetes')" :disabled="!isKubernetesInstalled">
						<img src="../assets/images/kubernetes-logo.png" alt="Kubernetes logo"
							class="platform-button-image" />
					</button>
				</div>
				<div class="platform-main-content-warnings-container">
					<!-- Show message if Docker is not installed -->
					<p v-if="!isDockerInstalled && isDockerInstalled != null" class="platform-warning-message">
						No valid Docker installation found in your system
						<button class="primary-button" @click="specifyDockerPath">Specify path</button>
					</p>

					<!-- Show message if Docker is not running -->
					<p v-if="!isDockerRunning && isDockerRunning != null && isDockerInstalled" class="platform-warning-message">
						A valid Docker installation was found, but it's not running. Please start it and try again.
						<button class="primary-button" @click="getInstalledPlatforms">Check again</button>
					</p>

					<!-- Show message if Kubernetes is not installed -->
					<p v-if="!isKubernetesInstalled && isKubernetesInstalled != null" class="platform-warning-message">
						No valid Kubernetes installation found in your system:
						<button class="primary-button" @click="specifyKubernetesPath">Specify path</button>
					</p>
				</div>
			</div>
		</div>
	</InstallationStep>
</template>