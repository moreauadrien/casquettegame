<script>
	import PlayerTag from '@/components/PlayerTag.svelte';

	import TeamTag from '@/components/TeamTag.svelte';
	import { Team } from '@/utils';

	const team = Team.BLUE;
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">Marc</p>
	<p class="mt-20 text-center font-inter font-medium text-4xl">Ton équipe</p>

	<div class="flex justify-center mt-4">
		<TeamTag {team} />
	</div>

	<div class="flex flex-col items-center mt-8">
		<PlayerTag username="Adrien" {team} />
		<PlayerTag username="John" {team} />
	</div>
</div>
