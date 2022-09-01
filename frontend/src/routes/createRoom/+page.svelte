<script lang="ts">
	import { goto } from '$app/navigation';
	import UsernameSelect from '@/views/UsernameSelect.svelte';

	import { createRoom } from '@/api';

    let error = ''

	async function handleSubmit(e: CustomEvent<string>) {
        const result = await createRoom(e.detail)

        if (result.err) {
            error = result.val.message
        } else {
            const roomId = result.val.roomId

            goto(`/room/${roomId}`)
        }
	}
</script>

<p>{error}</p>

<UsernameSelect on:submit={handleSubmit} title="Nouvelle partie" />
