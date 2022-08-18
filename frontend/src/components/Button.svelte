<script lang="ts">
	export let className = '';
</script>

<button
	on:click
	class="bg-teal-400 text-white text-3xl font-inter px-4 py-6 rounded-3xl font-medium {className ??
		''}"><slot /></button
>
