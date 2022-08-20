import type { Credentials, Payload } from './wsclient';
import { WsClient } from './wsclient';

import type { Writable, Readable } from 'svelte/store';
import { writable, derived } from 'svelte/store';
import type { PlayerInfos } from './api';

export enum GameState {
	NotConnected,
	WaitingRoom,
}

export class Game {
	private socket?: WsClient;
	private credentials?: Credentials;
	private gameState: Writable<GameState>;
	private roomId?: string;
	private hostId?: string;
	private wPlayers: Writable<PlayerInfos[]>;

	public players: Readable<PlayerInfos[]>;

	public constructor() {
		this.gameState = writable(GameState.NotConnected);
		this.wPlayers = writable([]);
		this.players = derived(this.wPlayers, (s) => s);
	}

	public createRoom() {
		this.errorIfNotConnected();

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
				this.hostId = this.credentials?.id;
				this.wPlayers.set(data.players);

				this.setGameState(GameState.WaitingRoom);

				clearTimeout(timeoutId);
				resolve(this.roomId);
			});
		});
	}

	public joinRoom(roomId: string) {
		this.errorIfNotConnected();

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
				console.log(data);
				if (data?.status === 'error') {
					reject(`ERROR: ${data?.message}`);
				}

				this.roomId = roomId;
				this.wPlayers.set(data.players);
				this.hostId = data.host;

				this.setGameState(GameState.WaitingRoom);

				clearTimeout(timeoutId);
				resolve(roomId);
			});
		});
	}

	private registerEventHandlers() {
		this.socket?.on('playerJoin', (data: any) => {
			this.wPlayers.set(data.players);
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

	private errorIfNotConnected() {
		if (this.socket === undefined) {
			throw new Error('the websocket is not connected');
		}
	}
}

export const game = new Game();
