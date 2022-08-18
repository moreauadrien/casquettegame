<script lang="ts">
	import BigCard from '@/components/BigCard.svelte';
	import Button from '@/components/Button.svelte';
	import Countdown from '@/components/Countdown.svelte';
	import CheckIcon from '@/components/icons/CheckIcon.svelte';
	import CrossIcon from '@/components/icons/CrossIcon.svelte';
	import TeamTag from '@/components/TeamTag.svelte';

	import { Team } from '@/utils';

	let team = Team.BLUE;

	function validateCard() {}

	function passCard() {}
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">Marc</p>

	<div class="flex justify-center mt-2">
		<TeamTag {team} />
	</div>

	<div class="mt-4 flex justify-center translate-y-2/4">
		<Countdown />
	</div>

	<div class="flex justify-center">
		<BigCard card="Harrison Ford" />
	</div>

	<div class="absolute bottom-8 w-full">
		<div class="m-4 flex justify-center">
			<Button className="w-72 py-2"><CheckIcon /></Button>
		</div>

		<div class="m-4 flex justify-center">
			<Button className="w-72 py-2 bg-red-400"><CrossIcon /></Button>
		</div>
	</div>
</div>
