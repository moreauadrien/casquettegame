import { WsClient } from './wsclient';

import { type Writable, type Readable, get } from 'svelte/store';
import { writable, derived } from 'svelte/store';
import type { PlayerInfos } from './api';
import type { TeamPoints } from './utils';
import { Team } from './utils';
import { qrCodeIsVisible } from './stores';

export const enum GameState {
	WaitingRoom = 'waitingRoom',
	CardSelection = 'cardSelection',
	TeamsRecap = 'teamsRecap',
	RulesRecap = 'rulesRecap',
	WaitTurnStart = 'waitTurnStart',
	Turn = 'turn',
	TurnRecap = 'turnRecap',
	ScoreRecap = 'scoreRecap',
	WaitPlayers = 'waitPlayers',
}

export class Game {
	private socket: WsClient;

	private host: PlayerInfos;
	private username: string;

	private wPlayers: Writable<PlayerInfos[]>;
	public players: Readable<PlayerInfos[]>;

	private wSwapsRemaining: Writable<number>;
	public swapsRemaining: Readable<number>;

	private wCards: Writable<string[]>;
	public cards: Readable<string[]>;

	private score: TeamPoints[];

	private speaker: PlayerInfos;

	private wGameState: Writable<GameState | undefined>;
	public gameState: Readable<GameState | undefined>;

	private wRound: Writable<number>;
	public round: Readable<number>;

	private wTeam: Writable<Team>;
	public team: Readable<Team>;

	public constructor() {
		this.socket = new WsClient();
		this.registerEventHandlers();

		this.wPlayers = writable([]);
		this.players = derived(this.wPlayers, (s) => s);

		this.wSwapsRemaining = writable(0);
		this.swapsRemaining = derived(this.wSwapsRemaining, (s) => s);

		this.wCards = writable([]);
		this.cards = derived(this.wCards, (s) => s);

		this.speaker = { username: '', id: '', team: Team.UNDEFINED };

		this.wGameState = writable(undefined);
		this.gameState = derived(this.wGameState, (s) => s);

		this.wRound = writable(0);
		this.round = derived(this.wRound, (s) => s);

		this.wTeam = writable(Team.UNDEFINED);
		this.team = derived(this.wTeam, (s) => s);

		this.host = { username: '', id: '', team: Team.UNDEFINED };
		this.score = [];
		this.username = '';
	}

	private registerEventHandlers() {
		this.socket.on('stateUpdate', (data: any) => {
			console.log({ data });
			const { state, speaker, players, host, round, team, username, cards, score, swapsRemaining } =
				data;

			if (username !== undefined) {
				this.username = username;
			}

			if (swapsRemaining !== undefined) {
				this.wSwapsRemaining.set(swapsRemaining);
			}

			if (host !== undefined) {
				this.host = host;

				qrCodeIsVisible.set(this.host.username === this.username);
			}

			if (speaker !== undefined && speaker !== null) {
				this.speaker = speaker;
			}

			if (players !== undefined) {
				this.wPlayers.set(players);
			}

			if (round !== undefined) {
				this.wRound.set(round);
			}

			if (team !== undefined) {
				this.wTeam.set(team as Team);
			}

			if (cards !== undefined) {
				this.wCards.set(cards);
			}

			if (score !== undefined) {
				this.score = score;
			}

			if (state !== undefined) {
				this.wGameState.set(state as GameState);
			}
		});
	}

	public connect() {
		return this.socket.connect();
	}

	public getTeam() {
		return this.team;
	}

	public getSpeaker() {
		return this.speaker ?? { id: '', username: '', team: Team.UNDEFINED };
	}

	public isSpeaker() {
		return this.username === this.speaker.username && this.username.length > 0;
	}

	public isHost() {
		return this.username === this.host.username && this.username.length > 0;
	}

	public getUsername() {
		return this.username;
	}

	public getScore() {
		return this.score;
	}

	public startGame() {
		this.socket.sendMessage('startGame');
	}

	public startTurn() {
		this.socket.sendMessage('startTurn');
	}

	public validateCard() {
		if (this.isSpeaker() === false) return;

		if (get(this.cards)) {
			this.wCards.update((current) => {
				current.shift();
				return current;
			});

			this.socket.sendMessage('validateCard');
		}
	}

	public passCard() {
		if (this.isSpeaker() === false) return;

		if (get(this.cards).length !== 0) {
			this.wCards.update((current) => {
				current.push(current.shift()!);
				return current;
			});

			this.socket.sendMessage('passCard');
		}
	}

	public handOver() {
		this.socket.sendMessage('handOver');
	}

	public startNextRound() {
		this.socket.sendMessage('startNextRound');
	}

	public acceptRules() {
		this.socket.sendMessage('acceptRules');
	}

	public changeCard(index: number) {
		this.socket.sendMessage(`changeCard:${index}`);
	}

	public validateCardSwitch() {
		this.socket.sendMessage('validateCardSwitch');
	}
}

export const game = new Game();
