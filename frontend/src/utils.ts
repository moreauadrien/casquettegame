export const enum Team {
	UNDEFINED = 'undefined',
	BLUE = 'blue',
	PURPLE = 'purple',
	RED = 'red',
	YELLOW = 'yellow',
}

export type TeamPoints = {
	team: Team;
	points: number[];
};

export function getRoomLink(roomId: string) {
	const { protocol, host } = window.location;

	return `${protocol}//${host}/join/${roomId}`;
}

export function sleep(time: number) {
	return new Promise<void>((resolve) => {
		setTimeout(() => {
			resolve();
		}, time);
	});
}
