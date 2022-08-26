<script lang="ts">
	import type { Team } from '@/utils';
	import RoundScore from './RoundScore.svelte';
	import TeamTag from './TeamTag.svelte';

	export let team: Team;
	export let scores: number[];
</script>

<div class="flex justify-center mt-5">
	<TeamTag {team} />
</div>

<div class="flex flex-col items-center">
	{#each scores as score, i}
		<RoundScore {team} {score} round={i + 1} />
	{/each}
</div>
