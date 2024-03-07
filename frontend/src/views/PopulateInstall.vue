<script>
import InstallationStep from '../components/InstallationStep.vue';
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { PopulateEnvironment } from '../../wailsjs/go/main/App';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';

const steps = [
	{ title: 'SpecifyPath', active: false },
	{ title: 'Populate', active: true }
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
					path: `/environments/${this.$store.state.populateState.id}`,
					disabled: false,
					text: 'Populate',
					onClick: () => {
						// Start the installation
						this.install();
					}
				},
				back: {
					path: '/populate/' + this.$store.state.populateState.id,
					disabled: false
				},
				cancel: {
					path: '/environments/' + this.$store.state.populateState.id,
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

			// Start the population
			this.updateTerminal('\nPopulating the environment\r\n');
			this.updateTerminal('Populating...\r\n');
			this.populate();

		},
		populate() {
			// Get the platform, environment and variables
			let envName = this.$store.state.populateState.name;
			let envVersion = this.$store.state.populateState.version;
			let platform = this.$store.state.populateState.platform;
			let path = this.$store.state.populateState.path;

			// Listen for the TERMINAL_OUTPUT event
			EventsOn('TERMINAL_OUTPUT', (output) => {
				this.updateTerminal(output);
			});

			// Populate the environment
			PopulateEnvironment(envName, envVersion, path, platform).then(() => {
				// Finish the installation
				this.finishInstallation();
			}).catch((error) => {
				this.updateTerminal('\nPopulation failed:\r\n');
				this.updateTerminal('Error: ' + error + '\r\n');
				this.installationError();
			});
		},
		finishInstallation() {
			// Update the UI
			this.installing = false;
			this.navigation.next.disabled = false;
			// Finish the installation
			this.updateTerminal('\nPopulation finished successfully!');
			// Update the navigation
			this.navigation.next.onClick = null;
			this.navigation.next.text = 'Back to Installed Environments';
		},
		installationError() {
			// Update the UI
			this.installing = false;
			this.navigation.cancel.disabled = false;
			this.navigation.back.disabled = false;
			this.navigation.next.disabled = false;
			this.navigation.next.text = 'Retry';
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

		this.updateTerminal('Press \'Populate\' to start the population process.\r\n')
	},
};
</script>


<template>
	<div class="install-container">
		<InstallationStep :navigation="navigation" :steps="steps" :tips="tips">
			<!-- The installation step's main content -->
			<div class='install-main-content-container'>
				<!-- The title-->
				<h1 class="install-title">Populate</h1>
				<!-- The main content container -->
				<div class="install-main-content">
					<!-- The terminal output -->
					<div class="terminal-output" ref="terminal"></div>
				</div>
			</div>
		</InstallationStep>
	</div>
</template>