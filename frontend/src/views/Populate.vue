<script>
import InstallationStep from '../components/InstallationStep.vue';
import { OpenFolderDialog } from '../../wailsjs/go/main/App';
import { mapState } from 'vuex';

const steps = [
	{ title: 'Select Folder', active: true },
	{ title: 'Populate', active: false }
];

const tips = 'Specify the path to a folder containing the metadata information to populate the environment with.';

export default {
	components: {
		InstallationStep
	},
	data() {
		return {
			steps,
			tips,
		};
	},
	computed: {
		...mapState({
			populateState: state => state.populateState
		}),
		navigation() {
			return {
				next: {
					path: '/populate-install/' + this.populateState.id,
					disabled: this.populateState.path === "" || this.populateState.path === null,
				},
				back: {
					path: `/environments/${this.populateState.id}`,
					disabled: false
				},
				cancel: {
					path: `/environments/${this.populateState.id}`,
					disabled: false
				}
			};
		}
	},
	methods: {
		openPathDialog() {
			OpenFolderDialog('Select the folder containing the metadata').then(path => {
				this.populateState.path = path;
			}).catch(err => {
				console.error(err);
			});
		},
	},
	mounted() {
		// Reset the state
		this.$store.commit('resetPopulateState');

		// Get the environment name and version from the router
		if (this.$route.params.id) {
			this.populateState.id = this.$route.params.id;
			this.populateState.name = this.populateState.id.split('-')[0];
			this.populateState.version = this.populateState.id.split('-')[1];
			this.populateState.platform = this.populateState.id.split('-')[2];
		}
	}
};
</script>


<template>
	<InstallationStep :steps="steps" :tips="tips" :navigation="navigation">
		<!-- The installation step's main content -->
		<div class='populate-main-content-container'>
			<!-- The title-->
			<h1 class="populate-title">Select Folder</h1>
			<!-- The form container -->
			<div class="populate-main-content">
				<form class="populate-form" autocomplete="off">
					<div class="populate-form-field">
						<label for="name">Path:</label>
						<input type="text" id="path" name="path" v-model="populateState.path" />
					</div>
				</form>
				<!-- Specify path button -->
				<button class="primary-button" @click="openPathDialog">Select folder</button>
			</div>
		</div>
	</InstallationStep>
</template>