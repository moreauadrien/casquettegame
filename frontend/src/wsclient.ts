export type EventHandler = (data: any) => void;

const WEBSOCKET_URL = `ws://${window.location.host}/ws`;

class Client {
	private socket?: WebSocket;
	private username: string;
	private token: string;
	private playerId: string;

	private handlers: Map<string, EventHandler>;

	public constructor() {
		this.username = '';
		this.token = '';
		this.playerId = '';
		this.handlers = new Map<string, EventHandler>();
	}

	public connect(username: string, playerId: string, token: string) {
		this.username = username;
		this.playerId = playerId;
		this.token = token;

		this.socket = new WebSocket(WEBSOCKET_URL);

		return new Promise((resolve, reject) => {
			const onOpen = () => {
				this.socket!.removeEventListener('open', onOpen);
				this.socket!.removeEventListener('error', onError);

				this.setupHandlers();

				resolve('connected');
			};

			const onError = () => {
				this.socket = undefined;
				reject('failed to connect');
			};

			this.socket!.addEventListener('open', onOpen);
			this.socket!.addEventListener('error', onError);
		});
	}

	public joinRoom(roomId: string) {
		this.emmitEvent('join', { roomId });
	}

	public emmitEvent(eventName: string, data: any) {
		if (this.socket === undefined) throw new Error('websocket is not connected');

		const payload = {
			username: this.username,
			id: this.playerId,
			token: this.token,
			event: {
				type: eventName,
				data,
			},
		};

		this.socket.send(JSON.stringify(payload));
	}

	public on(eventName: string, handler: EventHandler) {
		const handlerAlreadyExist = this.handlers.has(eventName);
		if (handlerAlreadyExist) {
			throw new Error('this event already has an handler');
		}

		this.handlers.set(eventName, handler);
	}

	private setupHandlers() {
		this.setupMessageHandler();
		this.setupErrorHandler();
		this.setupCloseHandler();
	}

	private setupMessageHandler() {
		this.socket!.addEventListener('message', (e) => {
			const payload = JSON.parse(e.data);

			if (payload.type === undefined) throw new Error(`the payload ${e.data} could not be parsed`);

			const handler = this.handlers.get(payload.type);
			if (handler !== undefined) {
				handler(payload.data);
			}
		});
	}

	private setupErrorHandler() {
		this.socket!.addEventListener('error', (e) => {});
	}

	private setupCloseHandler() {
		this.socket!.addEventListener('close', (e) => {
			this.socket = undefined;
		});
	}
}

export const client = new Client();
