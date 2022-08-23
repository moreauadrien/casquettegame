<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	import QrCodeView from '@/views/QrCodeView.svelte';
	import PlayerListView from '@/views/PlayerListView.svelte';
	import { game, GameState } from '@/game';
	import TeamView from '@/views/TeamView.svelte';
	import PreTurnView from '@/views/PreTurnView.svelte';
	import SpeakerView from '@/views/SpeakerView.svelte';
    import { goto } from '$app/navigation';
import SpectatorView from '@/views/SpectatorView.svelte';
import TurnRecapView from '@/views/TurnRecapView.svelte';

	const roomId = $page.params.slug;
	const gameState = game.state();
	const { players, cards } = game;

	let showQrCode = false;

	onMount(() => {
        if ($gameState === GameState.NotConnected) {
            goto(`/joinRoom/${roomId}`)
        }

		if ($gameState === GameState.WaitingRoom) {
			if (game.isHost()) {
				showQrCode = true;
			}
		}
	});
</script>

{#if $gameState === GameState.WaitingRoom}
	{#if showQrCode}
		<QrCodeView on:next={() => (showQrCode = false)} url={game.getRoomLink()} />
	{:else}
		<PlayerListView
			on:back={() => (showQrCode = true)}
			on:startGame={() => game.startRoom()}
			isHost={game.isHost()}
			players={$players}
		/>
	{/if}
{:else if $gameState === GameState.TeamsRecap}
	<TeamView players={$players} team={game.getTeam()} />
{:else if $gameState === GameState.WaitTurnStart}
	<PreTurnView
		on:startTurn={() => game.startTurn()}
		isSpeaker={game.isSpeaker()}
		speaker={game.getSpeaker()}
	/>
{:else if $gameState === GameState.Turn}
	{#if game.isSpeaker()}
		<SpeakerView
			cards={$cards}
			team={game.getTeam()}
			on:passCard={() => game.passCard()}
			on:validateCard={() => game.validateCard()}
		/>
    {:else}
        <SpectatorView cards={$cards} team={game.getSpeaker().team}/>
	{/if}
{:else if $gameState === GameState.TurnRecap}
    <TurnRecapView wasSpeaker={game.isSpeaker()} team={game.getSpeaker().team} cards={$cards}/>
{/if}
