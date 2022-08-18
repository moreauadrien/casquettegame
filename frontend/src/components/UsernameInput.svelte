<script lang="ts">
	export let value = '';
	export let placeholder = '';
</script>

<input
	class="bg-teal-50 rounded-3xl w-64 h-16 text-center font-inter font-medium text-2xl text-teal-800"
	type="text"
	bind:value
	{placeholder}
/>
