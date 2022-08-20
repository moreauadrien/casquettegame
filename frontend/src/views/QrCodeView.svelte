<script lang="ts">
    export let url: string
    
	import Button from '@/components/Button.svelte';
	import QrCode from '@/components/QrCode.svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<{ next: undefined }>();

	function handleClick() {
		dispatch('next');
	}
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">Nouvelle partie</p>

	<p class="mt-20 text-center font-inter font-medium text-4xl">Invite les<br />joueurs</p>

	<p class="mt-12 text-center font-inter font-medium text-xl">Scan le QR code</p>
	<div class="flex justify-center">
		<QrCode data={url} size={250} />
	</div>

	<div class="flex justify-center absolute bottom-28 w-full">
		<Button on:click={handleClick} className="px-8 py-4">Next</Button>
	</div>
</div>
