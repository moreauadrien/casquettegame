<script lang="ts">
	import Button from '@/components/Button.svelte';
	import { goto } from '$app/navigation';

	import { client, type Payload } from '@/wsclient';
	import { generateId, generateToken } from '@/utils';
	import { onMount } from 'svelte';

	onMount(async () => {
		const resp = await client.connect('adrien', generateId(), generateToken());
		console.log(resp);

		client.sendEvent({ type: 'test' }, (payload: Payload) => {
			console.log(payload);
		});
	});
</script>

<div class="flex justify-center mt-52">
	<Button on:click={() => goto('createRoom')}>Créer une partie</Button>
</div>
