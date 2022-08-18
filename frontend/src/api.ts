export type RoomConfiguration = {
	numberOfTeams: number;
	hostUsername: string;
	hostId: string;
	hostToken: string;
};

export async function createRoom(config: RoomConfiguration) {
	try {
		const response = await fetch('/api/createRoom/', {
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/json',
			},
			method: 'POST',
			body: JSON.stringify(config),
		});

		const jsonResponse = await response.json();
		return jsonResponse.roomId as string;
	} catch (e) {
		console.error(e);
	}
}
