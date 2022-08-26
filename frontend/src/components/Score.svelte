<script lang="ts">
	import type { Team } from '@/utils';
import TeamTag from './TeamTag.svelte';

	export let team: Team;
    export let score: number[]
</script>

<div>
	<div class="flex justify-center mt-4">
		<TeamTag {team} />
	</div>
</div>
