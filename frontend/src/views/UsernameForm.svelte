<script lang="ts">
	export let title: string

	import Button from '@/components/Button.svelte'
	import UsernameInput from '@/components/UsernameInput.svelte'
	import RightArrow from '@/components/icons/RightArrow.svelte'
	import { createEventDispatcher } from 'svelte'

	const dispatch = createEventDispatcher<{ submit: string }>()

	let username = ''

	function handleClick() {
		dispatch('submit', username)
	}
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center text-2xl">{title}</p>
	<p class="mt-20 text-center text-4xl">Entre<br />ton prénom</p>

	<div class="flex justify-center mt-8">
		<UsernameInput bind:value={username} placeholder="Maaaarc" />
	</div>

	<div class="flex justify-center absolute bottom-14 w-full">
		<Button on:click={handleClick} rightArrow>Next <RightArrow /></Button>
	</div>
</div>
