<script>
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import AboutSection from '../components/AboutSection.vue';
import { GetVersion } from '../../wailsjs/go/main/App'

export default {
	components: {
		AboutSection,
	},
	name: 'About',
	data() {
		return {
			title: 'About',
			description: 'EPOS Data portal is an open source, service-based data integration system. It is based on a microservices architecture.',
			version: '0.0.0',
		};
	},
	methods: {
		goToView(view) {
			this.$router.push(view);
		},
		openUrl(url) {
			BrowserOpenURL(url);
		},
	},
	created() {
		this.version = GetVersion().then((version) => {
			this.version = version;
		});
	},
};

</script>

<template>
	<div class="about">
		<div class="about-left-banner">
			<!-- the logo -->
			<img src="../assets/images/epos_logo.svg" alt="EPOS Data Portal" class="home-logo" />
			<h2 class="home-banner-info-title">Data Portal <br>Installer</h2>
		</div>
		<div class="about-main-content-container">
			<div class="about-main-content">
				<h1>{{ title }}</h1>
				<p>{{ description }}</p>
				<h3>License</h3>
				<p>
					The EPOS Data Portal is released under the GPLv3 license.
				</p>
				<h3>Application version</h3>
				<p>
					{{ version }}
				</p>
				<AboutSection title="Icons"
					description="All the icons in this application have been taken from the Icons8 website."
					url="https://icons8.it/icons"></AboutSection>
				<AboutSection title="Banner image" description="The banner image has been generated using PromeAI."
					url="https://www.promeai.com/image-variation">
				</AboutSection>
				<AboutSection title="EPOS EU" description="The EPOS EU website." url="https://www.epos-eu.org/">
				</AboutSection>
				<AboutSection title="EPOS Data Portal" description="The EPOS Data Portal website."
					url="https://www.epos-eu.org/dataportal">
				</AboutSection>
				<AboutSection title="EPOS Open Source documentation" description="The EPOS Open Source documentation."
					url="https://epos-eu.github.io/epos-open-source/">
				</AboutSection>
			</div>
			<div class="about-back-button-container">
				<button class="about-back-button" @click="goToView('/')"> Back </button>
			</div>
		</div>
	</div>
</template>