<script lang="ts">
	import BigCard from '@/components/BigCard.svelte';
	import Button from '@/components/Button.svelte';
	import Countdown from '@/components/Countdown.svelte';
	import CheckIcon from '@/components/icons/CheckIcon.svelte';
	import CrossIcon from '@/components/icons/CrossIcon.svelte';
	import TeamTag from '@/components/TeamTag.svelte';

	import type { Team } from '@/utils';
	import { createEventDispatcher } from 'svelte';

	export let cards: string[];
	export let team: Team;
    export let username: string

	const dispatch = createEventDispatcher();

	function validateCard() {
		if (cards[0] === undefined) return;
		dispatch('validateCard');
	}

	function passCard() {
		if (cards[0] === undefined) return;
		dispatch('passCard');
	}
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center text-2xl">{username}</p>

	<div class="flex justify-center mt-2">
		<TeamTag {team} />
	</div>

	<div class="mt-4 flex justify-center translate-y-2/4">
		<Countdown />
	</div>

	<div class="flex justify-center">
		<BigCard card={cards[0] ?? ''} />
	</div>

	<div class="absolute bottom-8 w-full">
		<div class="m-4 flex justify-center">
			<Button class="w-72 py-2" on:click={validateCard}><CheckIcon /></Button>
		</div>

		<div class="m-4 flex justify-center">
			<Button class="w-72 py-2 bg-red-400" on:click={passCard}><CrossIcon /></Button>
		</div>
	</div>
</div>
