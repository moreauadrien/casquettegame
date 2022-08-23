<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	import QrCodeView from '@/views/QrCodeView.svelte';
	import PlayerListView from '@/views/PlayerListView.svelte';
	import { game, GameState } from '@/game';
	import TeamView from '@/views/TeamView.svelte';
    import PreTurnView from '@/views/PreTurnView.svelte';

	const roomId = $page.params.slug;
	const gameState = game.state();
	const { players } = game;

	let showQrCode = false;

	onMount(() => {
		if ($gameState === GameState.NotConnected) return;

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
    <TeamView players={$players} team={game.getTeam()}/>
{:else if $gameState === GameState.WaitTurnStart}
    <PreTurnView on:startTurn={game.startTurn} isSpeaker={game.isSpeaker()} speaker={game.getSpeaker()}/>
{/if}
