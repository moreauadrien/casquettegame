<script lang="ts">
    export let big = false
    
    let defaultClasses = 'flex items-center justify-center bg-teal-400 text-white rounded-3xl font-semibold'
    let classes = ''

    $: {
        if (big) {
            classes = `${defaultClasses} text-3xl px-4 py-6`
        } else {
            classes = `${defaultClasses} text-2xl px-8 py-3`
        }

        if ($$props.class !== undefined) {
            classes += ` ${$$props.class}`
        }
    }

</script>

<button on:click class={classes}>
    <slot />
</button>
