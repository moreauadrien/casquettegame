<script lang="ts">
	import { goto } from '$app/navigation';

	import { token, username, playerId } from '@/stores';
	import UsernameSelect from '@/views/UsernameSelect.svelte';
	import { game } from '@/game';

	async function createRoom() {
		try {
			await game.connect({
				id: $playerId,
				token: $token,
				username: $username,
			});

            const roomId = await game.createRoom()
            goto(`/room/${roomId}`)
		} catch (e) {
			console.error(e);
		}
	}
</script>

<UsernameSelect
	on:submit={(e) => {
		username.set(e.detail);
		createRoom();
	}}
	title="Nouvelle partie"
/>
