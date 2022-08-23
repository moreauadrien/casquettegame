<script lang="ts">
import Card from '@/components/Card.svelte';

import TeamTag from '@/components/TeamTag.svelte';

	import { username } from '@/stores';
	import type { Team } from '@/utils';

	export let wasSpeaker: boolean;
	export let cards: string[];
	export let team: Team;
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">{$username}</p>

	<div class="flex justify-center mt-2">
		<TeamTag {team} />
	</div>

	<div class="mt-8 flex justify-center">
        <div class="w-80 h-14 flex justify-center items-center rounded-3xl mt-8 bg-teal-400">
            <p class="text-slate-50 text-2xl">Le temps est écoulé !</p>
        </div>
	</div>


	<p class="mt-10 text-center font-inter font-medium text-2xl">Cartes validées:</p>

	<div class="flex flex-col items-center mt-8">
		{#each cards as card}
			<Card value={card} />
		{/each}
	</div>
</div>
