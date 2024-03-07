<script>
export default {
	name: 'ConfirmDialog',
	props: {
		text: {
			type: String,
			default: 'Are you sure?'
		},
		confirmButton: {
			type: Object,
			default: () => ({ show: true, text: 'Confirm' }),
		},
		cancelButton: {
			type: Object,
			default: () => ({ show: true, text: 'Cancel' }),
		},
		title: {
			type: String,
			default: null,
		},
	},
	computed: {
		// default to positive, if not defined
		confirmButtonClass() {
			return {
				'dialog-button-positive': this.confirmButton.positive !== false,
				'dialog-button-negative': this.confirmButton.positive === false,
			};
		},
		// default to negative, if not defined
		cancelButtonClass() {
			return {
				'dialog-button-positive': this.cancelButton.positive === true,
				'dialog-button-negative': this.cancelButton.positive !== true,
			};
		},
	}
};
</script>

<template>
	<div class="dialog-overlay" v-if="(confirmButton && confirmButton.show) || (cancelButton && cancelButton.show)">
		<div class="dialog">
			<h3 v-if="title" class="dialog-title">{{ title }}</h3>
			<p class="dialog-text">{{ text }}</p>
			<div class="dialog-content">
				<button class="dialog-cancel-button" v-if="cancelButton && cancelButton.show" @click="$emit('cancel')"
					:class="cancelButtonClass">{{
						cancelButton.text }}</button>
				<button class="dialog-confirm-button" v-if="confirmButton && confirmButton.show" @click="$emit('confirm')"
					:class="confirmButtonClass">{{
						confirmButton.text }}</button>
			</div>
		</div>
	</div>
</template>