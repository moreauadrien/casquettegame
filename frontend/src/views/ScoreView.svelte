<script lang="ts">
	import Button from '@/components/Button.svelte';
	import TeamScores from '@/components/TeamScores.svelte';
	import { game } from '@/game';
	import { username } from '@/stores';
	import type { TeamPoints } from '@/utils';
	import { createEventDispatcher } from 'svelte';

	export let isHost: boolean;
	export let scores: TeamPoints[];

	const dispatch = createEventDispatcher();

	function handleClick() {
		dispatch('nextClick');
	}
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">{$username}</p>
	<p class="mt-3 text-center font-inter font-medium text-4xl">Scores</p>

	{#each scores as elem}
		<TeamScores team={elem.team} scores={elem.points} />
	{/each}

	{#if isHost && game.getRoundNumber() <= 3}
		<div class="flex justify-center absolute bottom-28 w-full">
			<Button on:click={handleClick} className="font-semibold">Next</Button>
		</div>
	{/if}
</div>
