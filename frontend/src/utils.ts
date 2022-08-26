import { v4 as uuidv4 } from 'uuid';

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

export function generateToken() {
	const preShuffled = uuidv4().replaceAll('-', '');
	return preShuffled
		.split('')
		.sort(function () {
			return 0.5 - Math.random();
		})
		.join('');
}

export function generateId() {
	return uuidv4();
}
