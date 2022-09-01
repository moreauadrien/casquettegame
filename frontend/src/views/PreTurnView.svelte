<script lang="ts">
	import type { PlayerInfos } from '@/api';
	import Button from '@/components/Button.svelte';
	import PlayerTag from '@/components/PlayerTag.svelte';

	import TeamTag from '@/components/TeamTag.svelte';
    import { createEventDispatcher } from 'svelte';

	export let isSpeaker: boolean;
	export let speaker: PlayerInfos;
    export let username: string

    const dispatch = createEventDispatcher()

    function startTurn() {
        dispatch('startTurn')
    }
</script>

<div class="flex flex-col justify-center">
    <p class="mt-14 text-center font-inter font-medium text-2xl">{username}</p>
	<p class="mt-3 text-center font-inter font-medium text-4xl">Première<br />manche</p>
	<p class="mt-20 text-center font-inter font-medium text-2xl">C'est au tour de :</p>

	<div class="flex justify-center mt-1">
		<TeamTag team={speaker.team} />
	</div>

	<p class="mt-4 text-center font-inter font-medium text-2xl">avec :</p>

	<div class="flex flex-col items-center mt-1">
		<PlayerTag username={speaker.username} team={speaker.team} />
	</div>

	{#if isSpeaker}
		<div class="flex justify-center absolute bottom-28 w-full">
			<Button on:click={startTurn} className="font-semibold">Ready ? Go !</Button>
		</div>
	{/if}
</div>
