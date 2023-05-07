<script lang="ts">
	import { goto } from '$app/navigation';
	import UsernameForm from '@/views/UsernameForm.svelte';

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

<UsernameForm on:submit={handleSubmit} title="Nouvelle partie" />
