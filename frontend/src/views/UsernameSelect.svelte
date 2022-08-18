<script lang="ts">
	export let title: string;

	import Button from '@/components/Button.svelte';
	import UsernameInput from '@/components/UsernameInput.svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<{ submit: string }>();

	let username = 'Mickey';

	function handleClick() {
		dispatch('submit', username);
	}
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">{title}</p>
	<p class="mt-20 text-center font-inter font-medium text-4xl">Entre<br />ton prénom</p>

	<div class="flex justify-center mt-8">
		<UsernameInput bind:value={username} placeholder="Maaaarc" />
	</div>

	<div class="flex justify-center absolute bottom-28 w-full">
		<Button on:click={handleClick} className="px-8 py-4">Next</Button>
	</div>
</div>
