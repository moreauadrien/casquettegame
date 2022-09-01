<script lang="ts">
	import Button from '@/components/Button.svelte';
	import TeamScores from '@/components/TeamScores.svelte';
	import type { TeamPoints } from '@/utils';
	import { createEventDispatcher } from 'svelte';

	export let isHost: boolean;
	export let score: TeamPoints[];
	export let round: number;
	export let username: string;

	const dispatch = createEventDispatcher();

	function handleClick() {
		dispatch('nextClick');
	}
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">{username}</p>
	<p class="mt-3 text-center font-inter font-medium text-4xl">Scores</p>

	{#each score as elem}
		<TeamScores team={elem.team} scores={elem.points} />
	{/each}

	{#if isHost && round < 3}
		<div class="flex justify-center absolute bottom-28 w-full">
			<Button on:click={handleClick} className="font-semibold">Next</Button>
		</div>
	{/if}
</div>
