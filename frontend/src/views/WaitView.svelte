<script lang="ts">
	export let username: string;
</script>

<div class="flex flex-col justify-center">
	<p class="mt-14 text-center text-2xl">{username}</p>

	<p class="mt-20 text-center text-4xl">En attente des autres joueurs</p>

	<div class="flex justify-center mt-16">
		<div class="loader" />
	</div>
</div>

<style>

    @keyframes spin {
      from {
        transform: rotate(0deg);
      }
      to {
        transform: rotate(360deg);
      }
    }

	.loader {
		border: 12px solid #f3f3f3;
		border-top: 12px solid #2dd4bf;
		border-radius: 50%;
		width: 80px;
		height: 80px;
		animation: spin 900ms linear infinite;
	}
</style>
