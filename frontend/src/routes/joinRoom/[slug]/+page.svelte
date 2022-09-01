<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
    import { joinRoom } from '@/api';
	import UsernameSelect from '@/views/UsernameSelect.svelte';

	const roomId = $page.params.slug;

    let error = ''

    async function handleSubmit(e: CustomEvent<string>) {
        const result = await joinRoom(e.detail, roomId)

        if (result.err) {
            error = result.val.message
        } else {
            goto(`/room/${roomId}`)
        }
    }
</script>

<p>{error}</p>

<UsernameSelect
	on:submit={handleSubmit}
	title="Rejoins une partie"
/>
