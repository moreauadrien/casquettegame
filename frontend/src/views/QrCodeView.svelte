<script lang="ts">
	export let url: string;

	import Button from '@/components/Button.svelte';
	import QrCode from '@/components/QrCode.svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<{ next: undefined }>();

	function handleClick() {
		dispatch('next');
	}
</script>

<div class="h-full">
	<p class="mt-14 text-center font-inter font-medium text-2xl">Nouvelle partie</p>

	<p class="mt-20 short:mt-12 text-center font-inter font-medium text-4xl">Invite les<br />joueurs</p>

	<p class="mt-12 text-center font-inter font-medium text-xl">Scan le QR code</p>
	<div class="flex justify-center">
        <QrCode data={url} class="w-64 short:w-48"/>
	</div>

	<div class="flex justify-center w-full">
		<Button on:click={handleClick} className="px-6 py-2">Next</Button>
	</div>
</div>
