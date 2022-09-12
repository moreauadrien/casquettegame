<script lang="ts">
    import type { PlayerInfos } from '@/api';

	export let isHost: boolean;
	import BackButton from '@/components/BackButton.svelte';
	import Button from '@/components/Button.svelte';
	import PlayerTag from '@/components/PlayerTag.svelte';

    export let players: PlayerInfos[]
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	function startGame() {
		dispatch('startGame');
	}

	function back() {
		dispatch('back');
	}
</script>

{#if isHost}
	<div class="absolute left-8 top-8">
		<BackButton on:click={back} />
	</div>
{/if}

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">Nouvelle partie</p>

	<p class="mt-20 short:mt-12 text-center font-inter font-medium text-4xl">Liste des<br />joueurs</p>

	<div class="flex flex-col items-center mt-8">
		{#each players as player}
			<PlayerTag username={player.username} team={player.team} />
		{/each}
	</div>

	{#if isHost}
		<div class="flex justify-center absolute bottom-12 w-full">
			<Button on:click={startGame} className="font-semibold">Here we go !</Button>
		</div>
	{/if}
</div>
