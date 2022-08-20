<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	import { playerId, token, username } from '@/stores';
	import UsernameSelect from '@/views/UsernameSelect.svelte';
	import QrCodeView from '@/views/QrCodeView.svelte';
	import PlayerListView from '@/views/PlayerListView.svelte';
	import { game, GameState } from '@/game';

	const roomId = $page.params.slug;
	const gameState = game.state();
    const { players } = game

	enum View {
		QRCode,
		PlayerList,
	}

	let currentView = View.PlayerList;

	onMount(() => {
        if ($gameState === GameState.NotConnected) return

		if ($gameState === GameState.WaitingRoom) {
			if (game.isHost()) {
				currentView = View.QRCode;
			}
		}
	});

	function nextClick() {
		if (currentView !== View.QRCode) return;
		currentView = View.PlayerList;
	}

	function startGame() {
		console.log('start game');
	}

	function handleBack() {
		//if (currentView !== View.PlayerList || $playerId !== $host) return;

		currentView = View.QRCode;
	}
</script>

{#if currentView === View.QRCode}
    <QrCodeView on:next={nextClick} url={game.getRoomLink()}/>
{:else if currentView === View.PlayerList}
    <PlayerListView on:back={handleBack} on:startGame={startGame} isHost={game.isHost()} players={$players} />
{/if}
