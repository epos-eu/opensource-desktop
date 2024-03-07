<script>
import Steps from './Steps.vue';
import Tips from './Tips.vue';

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
		}
	},
	// what the component provides
	components: {
		Steps,
		Tips
	}
};

</script>


<template>
	<div class="steps-and-tips">
		<div class="steps-and-tips-container">
			<div class="steps-and-tips-top-container">
				<!-- the logo -->
				<img src="../assets/images/epos_logo.svg" alt="EPOS Data Portal" class="home-logo" />
				<Steps :steps="steps" />
			</div>
			<Tips :tips="tips" />
		</div>
	</div>
</template>