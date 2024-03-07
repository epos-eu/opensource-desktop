<script>
export default {
	// what the component expects to receive
	props: {
		tips: {
			type: String,
			required: true
		},
	},
};
</script>

<template>
	<p v-if="tips" class="tips" v-html="tips"></p>
</template>