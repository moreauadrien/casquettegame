<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	import { playerId, token, username } from '@/stores';
	import { host } from '@/room';

	import { client } from '@/wsclient';
	import UsernameSelect from '@/views/UsernameSelect.svelte';
	import QrCodeView from '@/views/QrCodeView.svelte';
	import PlayerListView from '@/views/PlayerListView.svelte';

	const roomId = $page.params.slug;

	enum View {
		QRCode,
		PlayerList,
	}

	let currentView = View.QRCode;

	onMount(() => {
		if ($username.length === 0) return;

		joinRoom();
	});

	$: {
		if ($playerId !== $host && currentView === View.QRCode) {
			currentView = View.PlayerList;
		}
	}

	async function joinRoom() {
		try {
			await client.connect($username, $playerId, $token);

			client.joinRoom(roomId);
		} catch (e) {
			console.error(e);
		}
	}

	function nextClick() {
		if (currentView !== View.QRCode) return;
		currentView = View.PlayerList;
	}

	function startGame() {
		console.log('start game');
	}

	function handleBack() {
		if (currentView !== View.PlayerList || $playerId !== $host) return;

		currentView = View.QRCode;
	}
</script>

{#if $username.length === 0}
	<UsernameSelect
		on:submit={(e) => {
			username.set(e.detail);
			joinRoom();
		}}
		title="Rejoins une partie"
	/>
{:else if currentView === View.QRCode}
	<QrCodeView on:next={nextClick} />
{:else if currentView === View.PlayerList}
	<PlayerListView on:back={handleBack} on:startGame={startGame} isHost={$playerId === $host} />
{/if}
