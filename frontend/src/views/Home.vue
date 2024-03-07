<script>
import { IsInternetConnected, CheckForUpdates, DoUpdate } from '../../wailsjs/go/main/App'
import { BrowserOpenURL, Environment, Quit } from '../../wailsjs/runtime/runtime'
import ConfirmDialog from '../components/Dialog.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'

const bannerInfoText = "Unlock the power of the EPOS Data Portal\nInstall now for seamless data integration driven by metadata."

export default {
	components: {
		ConfirmDialog,
		LoadingSpinner,
	},
	data() {
		return {
			updateUrl: "https://github.com/",
			macDialog: {
				show: false,
				text: "There is a new version available. Do you want to download it?",
				confirm: { show: true, text: 'Download', },
				cancel: { show: true, text: 'Ignore', },
			},
			updateDialog: {
				show: false,
				text: "There is a new version available. Do you want to install it?",
				confirm: { show: true, text: 'OK', },
				cancel: { show: true, text: 'Ignore' },
			},
			updateFinishedDialog: {
				show: false,
				text: "The update has been installed. Please restart the application to apply the changes.",
				confirm: { show: true, text: 'OK', },
				cancel: null,
			},
			isCheckingForUpdate: false,
			isUpdating: false,
			bannerInfoText,
		};
	},
	computed: {
		// Check if the update dialog should be shown
		homeBannerFullscreen: {
			get() {
				return this.$store.state.homeBannerFullscreen;
			},
			set(value) {
				this.$store.commit('setHomeBannerFullscreen', value);
				// wait for the animation to finish
				setTimeout(() => {
					// Do the init
					this.init();
				}, 1200);
			}
		}
	},
	methods: {
		goToView(view) {
			this.$router.push(view);
		},
		openDocumentation() {
			BrowserOpenURL('https://epos-eu.github.io/epos-open-source/doc.html');
		},
		updateConfirm() {
			this.updateDialog.show = false;
			this.isUpdating = true;
			DoUpdate(this.updateUrl).then(() => {
				// Hide the loading spinner after 1 second
				setTimeout(() => {
					this.isUpdating = false;
					// If the update was successful, show the update dialog
					this.updateFinishedDialog.show = true;
				}, 1000);
			}).catch((error) => {
				console.error(error);
				this.updateFinishedDialog.show = true;
			});
		},
		updateCancel() {
			this.updateDialog.show = false;
		},
		macCancel() {
			this.macDialog.show = false;
		},
		macConfirm() {
			BrowserOpenURL(this.updateUrl);
		},
		updateFinishedConfirm() {
			this.updateFinishedDialog.show = false;
			// Close the application
			Quit();
		},
		init() {
			// Check if the internet is connected, if there isn't, show a message and exit
			IsInternetConnected();

			// Check if the check for updates has already been done
			if (!this.$store.state.checkForUpdatesDone) {

				Environment().then((systemEnvironment) => {
					this.isCheckingForUpdate = true;
					// Check for updates
					CheckForUpdates(systemEnvironment).then((updateUrl) => {
						// Hide the loading spinner after 1 second
						setTimeout(() => {
							this.isCheckingForUpdate = false;
							// Set the flag for updates done to true
							this.$store.commit('setCheckForUpdatesDone', true);
							// If there is an update available
							if (updateUrl !== "") {
								this.updateUrl = updateUrl;
								// If the user is on a Mac, show the Mac dialog
								if (systemEnvironment.platform === "darwin") {
									this.macDialog.show = true;
								} else {
									this.updateDialog.show = true;
								}
							}
						}, 1000);
					}).catch((error) => {
						console.error(error);
						this.$store.commit('setCheckForUpdatesDone', true);
						setTimeout(() => {
							this.isCheckingForUpdate = false;
						}, 1000);
					});
				});
			}
		},
	},
}

</script>

<template>
	<!-- Loading spinner while checking for the updates-->
	<LoadingSpinner :isLoading="isCheckingForUpdate" text="Checking for updates..." />
	<!-- Loading spinner while updating -->
	<LoadingSpinner :isLoading="isUpdating" text="Updating..." />
	<!-- Mac -->
	<ConfirmDialog v-if="macDialog.show" @confirm="macConfirm" @cancel="macCancel" :text="macDialog.text"
		:confirmButton="macDialog.confirm" :cancelButton="macDialog.cancel" :title="'Update available'" />
	<!-- Other -->
	<ConfirmDialog v-if="updateDialog.show" @confirm="updateConfirm" @cancel="updateCancel" :text="updateDialog.text"
		:title="'Update available'" :confirmButton="updateDialog.confirm" :cancelButton="updateDialog.cancel" />
	<ConfirmDialog v-if="updateFinishedDialog.show" @confirm="updateFinishedConfirm" :text="updateFinishedDialog.text"
		:title="'Update finished'" :confirmButton="updateFinishedDialog.confirm"
		:cancelButton="updateFinishedDialog.cancel" />
	<div class="home">
		<div class="home-left-banner">
			<!-- The logo -->
			<img src="../assets/images/epos_logo.svg" alt="EPOS Data Portal" class="home-logo" />
			<h2 class="home-banner-info-title">Data Portal <br>Installer</h2>

			<!-- The fullscreen banner info -->
			<div v-if="homeBannerFullscreen" class="home-banner-info-container">
				<p class="home-banner-info">
					Unlock the power of the EPOS Data Portal.<br>Install now for seamless data integration driven by
					metadata.
				</p>
				<button class="home-banner-button" @click="homeBannerFullscreen = !homeBannerFullscreen">
					<img src="../assets/images/rocket-icon-white.png" alt="EPOS Data Portal" class="home-icon" />Start
				</button>
			</div>
		</div>
		<!-- The buttons -->
		<transition name="expand">
			<div
				:class="{ 'home-main-content': !homeBannerFullscreen, 'home-main-content-collapsed': homeBannerFullscreen }">
				<div class="home-button-container">
					<button class="home-button" @click="openDocumentation()">
						<img src="../assets/images/document-icon.png" alt="EPOS Data Portal" class="home-icon" />
						Documentation
					</button>
				</div>
				<div class="home-button-container">
					<button class="home-button" @click="goToView('/platform')">
						<img src="../assets/images/rocket-icon.png" alt="EPOS Data Portal" class="home-icon" />
						Install </button>
				</div>
				<div class="home-button-container">
					<button class="home-button" @click="goToView('/environments')">
						<img src="../assets/images/database-icon.png" alt="EPOS Data Portal" class="home-icon" />
						Installed Environments </button>
				</div>
				<div class="home-button-container">
					<button class="home-button" @click="goToView('/about')">
						<img src="../assets/images/about-icon.png" alt="EPOS Data Portal" class="home-icon" />
						About</button>
				</div>
			</div>
		</transition>
	</div>
</template>
