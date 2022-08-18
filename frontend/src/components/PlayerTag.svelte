<script lang="ts">
	import { Team } from '@/utils';

	export let username: string;
	export let team: Team;

	const teamColors: { [key in Team]: string } = {
		[Team.BLUE]: 'bg-indigo-200 text-indigo-900',
		[Team.PURPLE]: 'bg-purple-200 text-purple-900',
		[Team.RED]: 'bg-rose-200 text-rose-900',
		[Team.YELLOW]: 'bg-amber-200 text-amber-900',
	};
</script>

<div class="w-72 h-10 flex justify-center items-center rounded-3xl m-1 {teamColors[team]}">
	<p class="text-2xl">{username}</p>
</div>
