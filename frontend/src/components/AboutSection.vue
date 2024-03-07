<script>
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

export default {
	props: ['title', 'description', 'url'],
	methods: {
		openUrl(url) {
			BrowserOpenURL(url);
		},
	},
};
</script>

<template>
	<h3>{{ title }}</h3>
	<div class="about-link-container">
		<p class="about-paragraph">
			{{ description }}
		</p>
		<button class="about-button" @click="openUrl(url)">Visit website</button>
	</div>
</template>