<script>
import InstallationStep from '../components/InstallationStep.vue';
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { InstallEnvironment, DeleteInstalledEnvironment } from '../../wailsjs/go/main/App';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';

const steps = [
	{ title: 'Platform', active: false },
	{ title: 'Environment', active: false },
	{ title: 'Variables', active: false },
	{ title: 'Install', active: true }
];

const tips = '';

export default {
	components: {
		InstallationStep
	},
	data() {
		return {
			terminal: null,
			installing: false,
			steps,
			tips,
			navigation: {
				next: {
					path: `/environments/${this.$store.state.installationState.id}`,
					disabled: false,
					text: 'Install',
					onClick: () => {
						// Start the installation
						this.install();
					}
				},
				back: {
					path: '/variables',
					disabled: false
				},
				cancel: {
					// /environments if editing, / if creating
					path: this.$store.state.editingEnvironment ? `/environments/${this.$store.state.installationState.id}` : '/',
					disabled: false,
				},
			}
		};
	},
	methods: {
		updateTerminal(message) {
			this.terminal.write(message + '\r\n');
		},
		install() {
			// Update the UI
			this.installing = true;
			this.navigation.next.disabled = true;
			this.navigation.cancel.disabled = true;
			this.navigation.back.disabled = true;

			// Start the installation or editing
			if (this.$store.state.editingEnvironment) {
				this.editInstalledEnvironment();
			} else {
				// Start the installation
				this.updateTerminal('Installing EPOS Open Source on your system');
				this.updateTerminal('Installing...');
				this.installNewEnvironment();
			}


		},
		installNewEnvironment() {
			// Get the platform, environment and variables
			let platform = this.$store.state.installationState.platform;
			let environment = this.$store.state.installationState.environmentSetup;
			let variables = this.$store.state.installationState.variables;
			let skipImagesAutoupdate = this.$store.state.installationState.skipImagesAutoupdate;
			let isEditing = this.$store.state.editingEnvironment ? true : false;

			// Listen for the TERMINAL_OUTPUT event
			EventsOn('TERMINAL_OUTPUT', (output) => {
				this.updateTerminal(output);
			});

			// Install from the Go backend
			InstallEnvironment(platform, environment, variables, skipImagesAutoupdate, isEditing).then(() => {
				// Finish the installation
				this.finishInstallation();
			}).catch((error) => {
				this.updateTerminal('\nInstallation failed:\r\n');
				this.updateTerminal('Error: ' + error + '\r\n');
				this.installationError();
			});
		},
		editInstalledEnvironment() {
			// Start the installation of the edited environment
			this.updateTerminal('\nEditing the installed environment\r\n');
			this.updateTerminal('Editing...\r\n');

			// Install the new environment
			this.installNewEnvironment();
		},
		finishInstallation() {
			// Update the UI
			this.installing = false;
			this.navigation.next.disabled = false;
			// Finish the installation
			this.updateTerminal('\nInstallation finished successfully!');
			// Update the navigation
			this.navigation.next.onClick = null;
			this.navigation.next.text = 'See in Installed Environments';
		},
		installationError() {
			// Update the UI
			this.installing = false;
			this.navigation.next.disabled = true;
			this.navigation.cancel.disabled = false;
			this.navigation.back.disabled = false;
		}
	},
	mounted() {
		// Create the terminal
		this.terminal = new Terminal({
			fontSize: 13,
			theme: {
				background: '#111111',
			},
			rows: 25,
			fontFamily: 'monospace',
			cursorInactiveStyle: 'none',

		});
		const fitAddon = new FitAddon();
		this.terminal.loadAddon(fitAddon);
		this.terminal.open(this.$refs.terminal);
		fitAddon.fit();

		// Resize the terminal when the window is resized
		window.addEventListener('resize', () => {
			fitAddon.fit();
		});
		if (this.$store.state.editingEnvironment) {
			this.updateTerminal('Press \'Install\' to start the editing process.\r\n');
		} else {
			this.updateTerminal('Press \'Install\' to start the installation process.\r\n');
		}
	},
};
</script>


<template>
	<div class="install-container">
		<InstallationStep :navigation="navigation" :steps="steps" :tips="tips">
			<!-- The installation step's main content -->
			<div class='install-main-content-container'>
				<!-- The title-->
				<h1 class="install-title">Install</h1>
				<!-- The main content container -->
				<div class="install-main-content">
					<!-- The terminal output -->
					<div class="terminal-output" ref="terminal"></div>
				</div>
			</div>
		</InstallationStep>
	</div>
</template>