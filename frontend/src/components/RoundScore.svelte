<script lang="ts">
	import { Team } from '@/utils';

	export let team: Team;
	export let score: number;
    export let round: number

	const teamColors: { [key in Team]: string } = {
		[Team.UNDEFINED]: '',
		[Team.BLUE]: 'bg-indigo-50 text-indigo-800',
		[Team.PURPLE]: 'bg-purple-50 text-purple-800',
		[Team.RED]: 'bg-rose-50 text-rose-800',
		[Team.YELLOW]: 'bg-amber-50 text-amber-800',
	};
</script>

<div class="w-72 h-10 flex justify-between items-center rounded-3xl m-1 {teamColors[team]}">
    <p class="text-base font-regular font-inter mx-8">Manche {round}</p>
    <p class="text-base font-regular font-inter mx-8">{score}</p>
</div>
