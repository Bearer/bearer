export function randomNames(quantity: number = 30) {
  return Array.from(Array(quantity).keys()).map(() => [bearSample(), nameSample()].join(' '))
}

function bearSample() {
  return bears[Math.floor(Math.random() * bears.length)]
}

function nameSample() {
  return names[Math.floor(Math.random() * names.length)]
}

const bears = [
  'Cinnamon',
  'Florida',
  'Glacier',
  'Haida Gwaii',
  'Kermode',
  'Spirit',
  'Louisiana',
  'Newfoundland',
  'Baluchistan',
  'Formosan',
  'Himalayan',
  'Ussuri',
  'Alaska Peninsula',
  'Atlas',
  'Bergman',
  'Cantabrian',
  'Gobi',
  'Grizzly',
  'Kamchatka',
  'Kodiak',
  'Marsican',
  'Sitka',
  'Stickeen',
  'Ussuri',
  'Giant',
  'Qinling',
  'Sloth',
  'Sun',
  'Polar',
  'Ursid hybrid',
  'Spectacled'
]
const names = [
  'Alpha',
  'Bravo',
  'Charlie',
  'Delta',
  'Echo',
  'Foxtrot',
  'Golf',
  'Hotel',
  'India',
  'Juliet',
  'Kilo',
  'Lima',
  'Mike',
  'November',
  'Oscar',
  'Papa',
  'Quebec',
  'Romeo',
  'Sierra',
  'Tango',
  'Uniform',
  'Victor',
  'Whiskey',
  'X-ray',
  'Yankee',
  'Zulu'
]
