<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	import QrCodeView from '@/views/QrCodeView.svelte';
	import PlayerListView from '@/views/PlayerListView.svelte';
	import { game, GameState } from '@/game';
	import TeamView from '@/views/TeamView.svelte';
	import PreTurnView from '@/views/PreTurnView.svelte';
	import SpeakerView from '@/views/SpeakerView.svelte';
	import SpectatorView from '@/views/SpectatorView.svelte';
	import TurnRecapView from '@/views/TurnRecapView.svelte';
	import ScoreView from '@/views/ScoreView.svelte';
	import { getRoomLink } from '@/utils';
	import { qrCodeIsVisible } from '@/stores';
	import RulesRecap from '@/views/RulesRecap.svelte';
	import CardSelectView from '@/views/CardSelectView.svelte';
import WaitView from '@/views/WaitView.svelte';

	const roomId = $page.params.slug;
	const { players, cards, gameState, round, team, swapsRemaining } = game;

	let error = '';

	onMount(async () => {
		const result = await game.connect();
		if (result.err) {
			error = result.val.message;
			return;
		}
	});
</script>

<p>{error}</p>

{#if $gameState === GameState.WaitingRoom}
	{#if $qrCodeIsVisible}
		<QrCodeView on:next={() => qrCodeIsVisible.set(false)} url={getRoomLink(roomId)} />
	{:else}
		<PlayerListView
			on:back={() => qrCodeIsVisible.set(true)}
			on:startGame={() => game.startGame()}
			isHost={game.isHost()}
			players={$players}
		/>
	{/if}
{:else if $gameState === GameState.TeamsRecap}
	<TeamView players={$players} team={$team} username={game.getUsername()} />
{:else if $gameState === GameState.WaitPlayers}
    <WaitView username={game.getUsername()} />
{:else if $gameState === GameState.CardSelection}
	<CardSelectView
		username={game.getUsername()}
		cards={$cards}
        swapsRemaining={$swapsRemaining}
		on:validate={() => game.validateCardSwitch()}
		on:change={(e) => game.changeCard(e.detail)}
	/>
{:else if $gameState === GameState.RulesRecap}
	<RulesRecap
		username={game.getUsername()}
		round={$round}
		isHost={game.isHost()}
		on:startRound={() => game.acceptRules()}
	/>
{:else if $gameState === GameState.WaitTurnStart}
	<PreTurnView
		on:startTurn={() => game.startTurn()}
		isSpeaker={game.isSpeaker()}
		speaker={game.getSpeaker()}
		username={game.getUsername()}
	/>
{:else if $gameState === GameState.Turn}
	{#if game.isSpeaker()}
		<SpeakerView
			cards={$cards}
			team={$team}
			username={game.getUsername()}
			on:passCard={() => game.passCard()}
			on:validateCard={() => game.validateCard()}
		/>
	{:else}
		<SpectatorView cards={$cards} team={game.getSpeaker().team} username={game.getUsername()} />
	{/if}
{:else if $gameState === GameState.TurnRecap}
	<TurnRecapView
		wasSpeaker={game.isSpeaker()}
		team={game.getSpeaker().team}
		cards={$cards}
		username={game.getUsername()}
		on:handOver={() => game.handOver()}
	/>
{:else if $gameState === GameState.ScoreRecap}
	<ScoreView
		score={game.getScore()}
		isHost={game.isHost()}
		on:nextClick={() => game.startNextRound()}
		round={$round}
		username={game.getUsername()}
	/>
{/if}
