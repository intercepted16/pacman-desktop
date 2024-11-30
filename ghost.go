package main

var Ghost = [4]Pacman{
	{x: cageX + 50, y: cageY + 50, radius: pacmanRadius, angle: 0, color: lightBlue, speed: 1, variety: CHASER},
	{x: cageX + 100, y: cageY + 50, radius: pacmanRadius, angle: 0, color: red, speed: 1, variety: AMBUSH},
	{x: cageX + 150, y: cageY + 50, radius: pacmanRadius, angle: 0, color: green, speed: 1, variety: PATROL},
	{x: cageX + 100, y: cageY + 100, radius: pacmanRadius, angle: 0, color: orange, speed: 1, variety: RANDOM},
}
