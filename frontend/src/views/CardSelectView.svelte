<script lang="ts">
	import Button from '@/components/Button.svelte';
	import Card from '@/components/Card.svelte';
	import { createEventDispatcher } from 'svelte';

	export let username: string;
	export let swapsRemaining: number;

    export let cards: string[];

	const dispatch = createEventDispatcher<{ change: number; validate: null }>();

	function plural(str: string, n: number) {
		return n > 1 ? `${str}s` : str;
	}

	function handleClick(index: number) {
		if (swapsRemaining <= 0) return;

		dispatch('change', index);
	}

	function handleValidate() {
		dispatch('validate');
	}
</script>

<div class="flex flex-col justify-center h-full">
    <div class="h-2/5">
        <p class="mt-14 text-center font-inter font-medium text-2xl">{username}</p>

        {#if swapsRemaining > 0}
            <p class="mt-8 text-center font-inter font-medium text-4xl">
                Tu peux encore échanger {swapsRemaining}
                {plural('carte', swapsRemaining)}
            </p>
        {:else}
            <p class="mt-8 text-center font-inter font-medium text-4xl">
                Tu ne peux plus échanger de carte
            </p>
        {/if}
    </div>
    
    <div class="h-3/5 overflow-auto">
        <div class="flex flex-col items-center mt-8">
            {#each cards as card, i}
                <Card value={card} on:click={() => handleClick(i)} />
            {/each}
        </div>
    </div>

    <div class="flex justify-center m-4">
        <Button on:click={handleValidate} className="font-semibold">Valider</Button>
    </div>
</div>
