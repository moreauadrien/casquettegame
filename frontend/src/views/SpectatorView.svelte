<script lang="ts">
	import Card from '@/components/Card.svelte';

	import Countdown from '@/components/Countdown.svelte';

	import TeamTag from '@/components/TeamTag.svelte';

	import type { Team } from '@/utils';

	export let cards: string[];

	export let team: Team;
	export let username: string;

</script>

<div class="flex flex-col justify-center h-full">
	<div class="h-2/5">
		<p class="mt-14 text-center text-2xl">{username}</p>

		<div class="flex justify-center mt-2">
			<TeamTag {team} />
		</div>

		<div class="mt-8 flex justify-center">
			<Countdown />
		</div>

		<p class="mt-10 text-center text-2xl">Cartes validées:</p>
	</div>

	<div class="h-3/5 overflow-auto">
		<div class="flex flex-col items-center mt-8">
			{#each cards as card}
				<Card value={card} />
			{/each}
		</div>
	</div>
</div>
