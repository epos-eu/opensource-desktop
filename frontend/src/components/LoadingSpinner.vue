<template>
	<div class="loading-spinner" v-if="isLoading">
		<div class="spinner"></div>
		<p class="spinner-text">{{ text }}</p>
	</div>
</template>

<script>
export default {
	props: ['isLoading', 'text'],
};
</script>
