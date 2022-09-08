<script lang="ts">
	export let data: string;

	import * as qr from 'qr-ts';

	let outputArea: HTMLDivElement;

	$: {
		generateQrcode(data, outputArea);
	}

	function generateQrcode(data: string, outputArea: HTMLDivElement | undefined) {
		if (outputArea === undefined) return;

		outputArea.innerHTML = '';

		const code = qr.renderOnCanvas(qr.generate(data), 'output');
        code.className = $$props.class
		outputArea.appendChild(code);
	}
</script>

<div bind:this={outputArea} class="inline-block" />
