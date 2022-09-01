export type RuleLine = {
	text: string;
	gray?: boolean;
};

type Rules = {
	[key: number]: RuleLine[];
};

const rules: Rules = {
	1: [
		{
			text: 'L’orateur n’a pas de limite de mots.',
		},
		{
			text: 'Les coéquipiers n’ont pas de limite de propositions.',
		},
		{
			text: 'L’orateur ne peut pas :',
		},
		{
			text: 'Prononcer des parties ou des diminutifs et prénoms indiqués sur la carte.',
		},
		{
			text: 'Ex : “C’est la fille de Serge Gainsbourg.”',
			gray: true,
		},
		{
			text: 'Utiliser des traductions directes.',
		},
		{
			text: 'Ex : “C’est Michael White en français.”',
			gray: true,
		},
		{
			text: 'Énumérer des lettres de l’alphabet.',
		},
		{
			text: 'Ex : “Ça commence par un b.”',
			gray: true,
		},
	],
	2: [
		{
			text: 'L’orateur ne peut faire qu’une proposition par carte.',
		},
		{
			text: 'Les coéquipiers ne peuvent donner qu’une seule proposition.',
		},

		{
			text: 'Si l’équipe comporte plus de 2 joueurs, seule la première proposition compte. Mettez vous d’accord avant !',
			gray: true,
		},
		{
			text: 'L’orateur ne peut pas :',
		},
		{
			text: 'Prononcer des parties ou des diminutifs des noms et prénoms indiqués.',
		},
		{
			text: 'Utiliser des traductions directes.',
		},

		{
			text: 'Ex : “white” pour faire deviner Michel Blanc',
			gray: true,
		},
	],
	3: [
		{
			text: 'L’orateur n’a plus le droit de parler. Il peut mimer ou faire des bruitages (onomatopées).',
		},
		{
			text: 'Les coéquipiers n’ont pas de limite de propositions.',
		},
		{
			text: 'L’orateur ne peut pas fredonner de chanson.',
		},
	],
};

export default rules;
