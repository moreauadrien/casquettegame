<script lang="ts">
	export let size: number;
	export let data: string;

	import { onMount } from 'svelte';
	import * as qr from 'qr-ts';

	let outputArea: HTMLDivElement;

	onMount(() => {
		const options = {
			c1: '#000000',
			c2: '#FFFFFF',
			size: size / 33,
		};

		const code = qr.renderOnCanvas(qr.generate(data), 'output', options);
		outputArea.appendChild(code);
	});
</script>

<div bind:this={outputArea} class="inline-block" />
