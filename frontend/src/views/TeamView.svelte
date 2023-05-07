<script lang="ts">
	import type { PlayerInfos } from '@/api';

    export let username: string
	export let players: PlayerInfos[];
	export let team: Team;

    let teamMembers: PlayerInfos[]
    $: {
        teamMembers = players.filter(p => p.team === team)
    }

	import PlayerTag from '@/components/PlayerTag.svelte';

	import TeamTag from '@/components/TeamTag.svelte';
	import type { Team } from '@/utils';
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center text-2xl">{username}</p>
	<p class="mt-20 text-center text-4xl">Ton équipe</p>

	<div class="flex justify-center mt-4">
		<TeamTag {team} />
	</div>

	<div class="flex flex-col items-center mt-8">
		{#each teamMembers as player}
			<PlayerTag username={player.username} {team} />
		{/each}
	</div>
</div>
