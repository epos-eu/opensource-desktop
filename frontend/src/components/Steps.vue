<script>
import Step from './Step.vue';

export default {
	// what the component expects to receive
	props: {
		steps: {
			type: Array,
			required: true
		}
	},
	// what the component provides
	components: {
		Step
	}
};

</script>

<template>
	<div class="steps-container">
		<div v-for="(step, index) in steps" :key="index">
			<Step :step="step" />
		</div>
	</div>
</template>