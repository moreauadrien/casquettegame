import type { Credentials } from './wsclient';
import { WsClient } from './wsclient';

import type { Writable, Readable } from 'svelte/store';
import { get } from 'svelte/store';
import { writable, derived } from 'svelte/store';
import type { PlayerInfos } from './api';
import { Team, type TeamPoints } from './utils';

export const enum GameState {
	NotConnected = 'notConnected',
	WaitingRoom = 'waitingRoom',
	TeamsRecap = 'teamsRecap',
	RulesRecap = 'rulesRecap',
	WaitTurnStart = 'waitTurnStart',
	Turn = 'turn',
	TurnRecap = 'turnRecap',
	ScoreRecap = 'scoreRecap',
}

export class Game {
	private socket?: WsClient;
	private credentials?: Credentials;
	private gameState: Writable<GameState>;
	private roomId?: string;
	private hostId?: string;
	private roundNumber: number;
	private team: Team;

	private wPlayers: Writable<PlayerInfos[]>;
	public players: Readable<PlayerInfos[]>;

	private wCards: Writable<string[]>;
	public cards: Readable<string[]>;

	private wScores: Writable<TeamPoints[]>;
	public scores: Readable<TeamPoints[]>;

	private speaker?: PlayerInfos;

	public constructor() {
		this.gameState = writable(GameState.NotConnected);
		this.roundNumber = 1;

		this.team = Team.UNDEFINED;

		this.wPlayers = writable([]);
		this.players = derived(this.wPlayers, (s) => s);

		this.wCards = writable([]);
		this.cards = derived(this.wCards, (s) => s);

		this.wScores = writable([]);
		this.scores = derived(this.wScores, (s) => s);
	}

	public createRoom() {
		const createEvent = {
			type: 'create',
		};

		return new Promise<string>((resolve, reject) => {
			const timeoutId = setTimeout(() => {
				reject('timeout');
			}, 5000);

			this.socket?.sendEvent(createEvent, (data: any) => {
				if (data?.status === 'error') {
					reject(`ERROR: ${data?.message}`);
				}

				if (typeof data?.roomId !== 'string' || (data?.roomId as string).length === 0) {
					reject('roomId is not valid');
				}

				this.roomId = data.roomId as string;
				this.team = data.team as Team;
				this.hostId = this.credentials?.id;
				this.wPlayers.set(data.players);

				this.setGameState(GameState.WaitingRoom);

				clearTimeout(timeoutId);
				resolve(this.roomId);
			});
		});
	}

	public joinRoom(roomId: string) {
		const joinEvent = {
			type: 'join',
			data: {
				roomId,
			},
		};

		return new Promise<string>((resolve, reject) => {
			const timeoutId = setTimeout(() => {
				reject('timeout');
			}, 5000);

			this.socket!.sendEvent(joinEvent, (data: any) => {
				if (data?.status === 'error') {
					reject(`ERROR: ${data?.message}`);
				}

				this.roomId = roomId;
				this.team = data.team;
				this.wPlayers.set(data.players);
				this.hostId = data.host;

				this.setGameState(GameState.WaitingRoom);

				clearTimeout(timeoutId);
				resolve(roomId);
			});
		});
	}

	public startRoom() {
		if (this.isHost() === false) {
			throw new Error('you have to be the host to start the game');
		}

		this.socket!.sendEvent({ type: 'startGame' });
	}

	public startTurn() {
		if (this.isSpeaker() === false) {
			throw new Error('you have to be the speaker to start the turn');
		}

		this.socket!.sendEvent({ type: 'startTurn' }, (data: any) => {
			if (data.status === 'success') {
				this.wCards.set(data.cards);

				this.setGameState(GameState.Turn);
			}
		});
	}

	public passCard() {
		if (this.isSpeaker() === false) {
			throw new Error("you can't pass card if you are not the speaker");
		}

		if (get(this.cards).length !== 0) {
			this.socket!.sendEvent({ type: 'passCard', data: { cards: get(this.cards)[0] } });

			this.wCards.update((current) => {
				current.push(current.shift()!);
				return current;
			});
		}
	}

	public handOver() {
		if (this.isSpeaker() === false) {
			throw new Error("you can't hand over when you are not the speaker");
		}

		this.socket!.sendEvent({ type: 'handOver' });
	}

	public validateCard() {
		if (this.isSpeaker() === false) {
			throw new Error("you can't validate card if you are not the speaker");
		}

		if (get(this.cards)) {
			this.wCards.update((current) => {
				current.shift();
				return current;
			});

			this.socket!.sendEvent({ type: 'validateCard', data: { cards: get(this.cards)[0] } });
		}
	}

	public nextRound() {
		this.socket!.sendEvent({ type: 'nextRound' });
	}

	private registerEventHandlers() {
		this.socket?.on('playerJoin', (data: any) => {
			this.wPlayers.set(data.players);
		});

		this.socket?.on('stateUpdate', (data: any) => {
			const state = data.state as GameState;
			this.setGameState(state);

			switch (state) {
				case GameState.TeamsRecap:
					this.wPlayers.set(data.players);
					break;

				case GameState.WaitTurnStart:
					this.speaker = data.speaker;
					this.wCards.set([]);
					break;

				case GameState.TurnRecap:
					this.wCards.set(data.cards ?? []);
					break;

				case GameState.ScoreRecap:
					this.roundNumber++;
					this.wScores.set(data.scores);
					break;
			}
		});

		this.socket?.on('turnUpdate', (data: any) => {
			this.wCards.set(data.cards);
		});
	}

	public getRoomId() {
		return this.roomId;
	}

	public getRoomLink() {
		const { protocol, host } = window.location;

		return `${protocol}//${host}/joinRoom/${this.roomId}`;
	}

	public state() {
		return derived(this.gameState, (s) => s);
	}

	private setGameState(gameState: GameState) {
		this.gameState.set(gameState);
	}

	public connect(credentials: Credentials) {
		this.credentials = credentials;

		this.socket = new WsClient(this.credentials);
		this.registerEventHandlers();

		return this.socket.connect();
	}

	public isHost() {
		return this.hostId !== undefined && this.hostId === this.credentials?.id;
	}

	public isSpeaker() {
		return this.speaker !== undefined && this.speaker.id === this.credentials?.id;
	}

	public getTeam() {
		return this.team;
	}

	public getSpeaker() {
		return this.speaker ?? { id: '', username: '', team: Team.UNDEFINED };
	}

	public getRoundNumber() {
		return this.roundNumber;
	}
}

export const game = new Game();
