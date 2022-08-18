<script lang="ts">
	import { Team } from '@/utils';

	export let team: Team;

	const teamBackgroundColor: { [key in Team]: string } = {
		[Team.BLUE]: 'bg-indigo-500',
		[Team.PURPLE]: 'bg-purple-500',
		[Team.RED]: 'bg-rose-500',
		[Team.YELLOW]: 'bg-amber-500',
	};

	const teamNumber: { [key in Team]: number } = {
		[Team.BLUE]: 1,
		[Team.PURPLE]: 2,
		[Team.RED]: 3,
		[Team.YELLOW]: 4,
	};
</script>

<div class="w-28 h-7 flex justify-center items-center rounded-3xl m-1 {teamBackgroundColor[team]}">
	<p class="text-white text-xl">Equipe {teamNumber[team]}</p>
</div>
