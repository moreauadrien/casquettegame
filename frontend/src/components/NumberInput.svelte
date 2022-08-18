<script lang="ts">
	export let min: number;
	export let max: number;
	export let value: number;

	function decrement() {
		if (value <= min) return;
		value--;
	}

	function increment() {
		if (value >= max) return;
		value++;
	}
</script>

<div class="flex">
	<div
		on:click={decrement}
		class="w-24 h-24 m-2 bg-teal-100 rounded-3xl flex justify-center items-center font-inter font-medium text-4xl"
		class:text-gray-400={value <= min}
	>
		-
	</div>
	<div
		class="w-24 h-24 m-2 bg-teal-50 rounded-3xl flex justify-center items-center font-inter font-medium text-4xl"
	>
		{value}
	</div>
	<div
		on:click={increment}
		class="w-24 h-24 m-2 bg-teal-100 rounded-3xl flex justify-center items-center font-inter font-medium text-4xl"
		class:text-gray-400={value >= max}
	>
		+
	</div>
</div>
