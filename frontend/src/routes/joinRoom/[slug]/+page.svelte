<script lang="ts">
	import { goto } from '$app/navigation';

	import { page } from '$app/stores';
	import { game, GameState } from '@/game';
    import { username, token, playerId } from '@/stores';
	import UsernameSelect from '@/views/UsernameSelect.svelte';

	const roomId = $page.params.slug;
	const gameState = game.state();

	$: {
		if ($gameState !== GameState.NotConnected) {
			goto(`/room/${game.getRoomId()}`);
		}
	}

    async function handleUsernameSubmit(e: CustomEvent<string>) {
        username.set(e.detail)

        try {
            await game.connect({username: $username, token: $token, id: $playerId})
            await game.joinRoom(roomId)
        } catch (e) {
            console.log(e)
        }
    }
</script>

<UsernameSelect
	on:submit={handleUsernameSubmit}
	title="Rejoins une partie"
/>
