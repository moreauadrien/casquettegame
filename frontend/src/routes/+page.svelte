<script lang="ts">
	import Button from '@/components/Button.svelte';
	import { goto } from '$app/navigation';

</script>

<div class="flex justify-center mt-52">
	<Button on:click={() => goto('createRoom')}>Créer une partie</Button>
</div>
