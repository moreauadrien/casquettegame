<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	import { playerId, token, username } from '@/stores';
	import UsernameSelect from '@/views/UsernameSelect.svelte';
	import QrCodeView from '@/views/QrCodeView.svelte';
	import PlayerListView from '@/views/PlayerListView.svelte';
    import { Room } from '@/room';

	const roomId = $page.params.slug;

	enum View {
		QRCode,
		PlayerList,
	}

	let currentView = View.QRCode;

    let room = new Room()

	onMount(() => {
		if ($username.length === 0) return;
        
        joinRoom()
	});

	async function joinRoom() {
		try {
            await room.connect({username: $username, token: $token, id: $playerId})
            await room.join(roomId)

            console.log("ok")
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
		//if (currentView !== View.PlayerList || $playerId !== $host) return;

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
	<PlayerListView on:back={handleBack} on:startGame={startGame} isHost={true} />
{/if}
