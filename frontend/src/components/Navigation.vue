<script>
export default {
	props: {
		navigation: {
			type: Object,
			default: null
		},
	},
	mounted() {
		window.addEventListener('keyup', this.onKeyUp);
	},
	beforeDestroy() {
		window.removeEventListener('keyup', this.onKeyUp);
	},
	methods: {
		onKeyUp(event) {
			// Check if the 'Enter' key was pressed
			if (event.key === 'Enter' || event.keyCode === 13) {
				// Check if the 'Next' button is not disabled and if it is rendered
				if (this.navigation && !this.navigation.next.disabled && this.$refs.nextButton) {
					// Simulate a click event on the 'Next' button
					this.$refs.nextButton.click();
				}
			}
		},
		nextButtonClick() {
			// Check if the 'Next' button is not disabled and if it is rendered
			if (this.navigation && !this.navigation.next.disabled && this.navigation.next.onClick) {
				// Call the 'onClick' function
				this.navigation.next.onClick();
			} else {
				// Redirect to the 'next' path
				this.$router.push(this.navigation.next.path);
			}
		},
		cancelButtonClick() {
			// Check if the 'Cancel' button is not disabled and if it is rendered
			if (this.navigation && this.navigation.cancel.onClick && !this.navigation.cancel.disabled) {
				// Redirect to the 'cancel' path
				this.navigation.cancel.onClick();
			} else {
				// Redirect to the 'cancel' path
				this.$router.push(this.navigation.cancel.path);
			}
		},
		backButtonClick() {
			// Check if the 'Back' button is not disabled and if it is rendered
			if (this.navigation && this.navigation.back.onClick && !this.navigation.back.disabled) {
				// Call the 'onClick' function
				this.navigation.back.onClick();
			} else {
				// Redirect to the 'back' path
				this.$router.push(this.navigation.back.path);
			}
		},
	},
	computed: {
		nextButtonText() {
			return this.navigation && this.navigation.next.text ? this.navigation.next.text : 'Next';
		},
		backButtonText() {
			return this.navigation && this.navigation.back.text ? this.navigation.back.text : 'Back';
		},
		cancelButtonText() {
			return this.navigation && this.navigation.cancel.text ? this.navigation.cancel.text : 'Cancel';
		},
	},
};

</script>

<template>
	<div v-if="navigation" class="navigation-container">
		<button
			:class="{ 'secondary-button': !navigation.cancel.disabled, 'secondary-button-inactive': navigation.cancel.disabled }"
			:disabled="navigation.cancel.disabled" @click="cancelButtonClick">{{ cancelButtonText }}</button>
		<div class="navigation-back-next-container">
			<button
				:class="{ 'primary-button': !navigation.back.disabled, 'primary-button-inactive': navigation.back.disabled }"
				:disabled="navigation.back.disabled" @click="backButtonClick">{{ backButtonText }}</button>
			<button ref="nextButton"
				:class="{ 'primary-button': !navigation.next.disabled, 'primary-button-inactive': navigation.next.disabled }"
				:disabled="navigation.next.disabled" @click="nextButtonClick">{{ nextButtonText }}</button>
		</div>
	</div>
</template>