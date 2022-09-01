<script lang="ts">
	import { sleep } from '@/utils';

	export let value: string;
	let displayedValue: string;

	let rotated = false;
	let div: HTMLDivElement;

	$: {
		animateTo(value);
	}

	async function animateTo(value: string) {
		if (displayedValue === undefined) {
			displayedValue = value;
			return;
		}

		rotated = true;
		await sleep(100);
		displayedValue = value;
		rotated = false;
	}
</script>

<div
	on:click
	bind:this={div}
	class="w-72 h-11 flex justify-center items-center rounded-3xl m-1 bg-teal-50 text-teal-800 {rotated
		? 'rotated'
		: ''}"
>
	<p class="text-2xl font-medium font-inter">{displayedValue}</p>
</div>

<style>
	div {
		transition: transform 100ms linear;
	}

	.rotated {
		transform: rotateX(90deg);
	}
</style>
