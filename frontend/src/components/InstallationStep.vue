<script>
import StepsAndTips from '../components/StepsAndTips.vue';
import Navigation from '../components/Navigation.vue';


export default {
	// what the component expects to receive
	props: {
		tips: {
			type: String,
			required: true
		},
		steps: {
			type: Array,
			required: true
		},
		navigation: {
			type: Object,
			required: false,
			default: null
		}
	},
	// what the component provides
	components: {
		StepsAndTips,
		Navigation
	},
	data() {
		return {
			// just a copy of the steps prop
			thisSteps: this.steps
		}
	},
	mounted() {
		// if editing, remove the platform step and the environment setup step
		if (this.$store.state.editingEnvironment) {
			this.thisSteps = this.thisSteps.slice(2);
		}
	}
};
</script>

<template>
	<div class="installation-step-container">
		<StepsAndTips :steps="thisSteps" :tips="tips" />

		<div class="installation-step-main-content-container">
			<slot></slot>
			<Navigation :navigation="navigation" />
		</div>
	</div>
</template>