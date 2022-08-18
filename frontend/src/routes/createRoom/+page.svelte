<script lang="ts">
	import { createRoom } from '@/api';
	import { goto } from '$app/navigation';

	import { token, username, playerId } from '@/stores';
	import TeamsNumberSelect from '@/views/TeamsNumberSelect.svelte';
	import UsernameSelect from '@/views/UsernameSelect.svelte';

	enum State {
		ChooseNumberOfTeams,
		ChooseYourUserName,
	}

	let currentState = State.ChooseNumberOfTeams;

	let numberOfTeams: number;

	async function nextClick() {
		switch (currentState) {
			case State.ChooseNumberOfTeams:
				currentState = State.ChooseYourUserName;
				break;

			case State.ChooseYourUserName:
				const roomId = await createRoom({
					numberOfTeams,
					hostUsername: $username,
					hostId: $playerId,
					hostToken: $token,
				});

				if (roomId !== undefined) {
					goto(`/room/${roomId}`);
				}

				break;
		}
	}
</script>

{#if currentState === State.ChooseNumberOfTeams}
	<TeamsNumberSelect
		on:submit={(e) => {
			numberOfTeams = e.detail;
			nextClick();
		}}
	/>
{:else if currentState === State.ChooseYourUserName}
	<UsernameSelect
		on:submit={(e) => {
			username.set(e.detail);
			nextClick();
		}}
		title="Nouvelle partie"
	/>
{/if}
