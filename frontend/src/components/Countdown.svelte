<script lang="ts">

    let value = 30

    const timer = setInterval(() => {
        value--
        if (value === 0) {
            clearInterval(timer)
        }
    }, 1000)
</script>

<div class="inline-block">
	<div class="w-16 h-16 bg-teal-500 flex justify-center items-center rounded-full inline-block">
        <p class="text-white font-inter font-medium text-3xl">{value}</p>
	</div>
</div>
