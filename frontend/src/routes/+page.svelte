<script lang="ts">
	import Button from '@/components/Button.svelte';
	import { goto } from '$app/navigation';

</script>

<div class="flex justify-center absolute bottom-28 w-full">
	<Button on:click={() => goto('create')} big>Créer une partie</Button>
</div>
