<script lang="ts">
    import Button from "@/components/Button.svelte";
    import NumberInput from "@/components/NumberInput.svelte";
    import { createEventDispatcher } from "svelte";

    let dispatch = createEventDispatcher<{submit: number}>()

    let numberOfTeams = 2

    function handleClick() {
        dispatch('submit', numberOfTeams)
    }

</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center font-inter font-medium text-2xl">Nouvelle partie</p>
		<p class="mt-20 text-center font-inter font-medium text-4xl">Nombre<br />d’équipes</p>

		<div class="flex justify-center mt-8">
			<NumberInput min={2} max={2} bind:value={numberOfTeams} />
		</div>

	<div class="flex justify-center absolute bottom-28 w-full">
		<Button on:click={handleClick} className="px-8 py-4">Next</Button>
	</div>
</div>
