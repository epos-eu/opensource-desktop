<template>
	<div :class="{ 'step-active': step.active, 'step-inactive': !step.active }">{{ step.title }}</div>
</template>

<script>
export default {
	// what the component expects to receive
	props: {
		step: {
			type: Object,
			required: true
		}
	}
};
</script>