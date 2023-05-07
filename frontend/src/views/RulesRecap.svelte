<script lang="ts">
	import Button from '@/components/Button.svelte';
	import rules from '@/rules';
import { createEventDispatcher } from 'svelte';

	export let username: string;
	export let round: number;
	export let isHost: boolean;

	const title: { [key: number]: string } = {
		1: 'Première<br/>manche',
		2: 'Deuxième<br/>manche',
		3: 'Troisième<br/>manche',
	};

    const dispatch = createEventDispatcher()

    function handleClick() {
        dispatch("startRound")
    }
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center text-2xl">{username}</p>
	<p class="mt-3 text-center text-4xl">{@html title[round]}</p>

	<div class="flex justify-center mt-4">
		<div class="w-10/12">
			{#each rules[round] as elem}
                <p class="mt-2 {elem.gray === true ? "text-gray-500": ""}">{elem.text}</p>
			{/each}
		</div>
	</div>

	{#if isHost}
		<div class="flex justify-center absolute bottom-28 w-full">
			<Button on:click={handleClick}>Chef, oui chef</Button>
		</div>
	{/if}
</div>
