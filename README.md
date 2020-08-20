# Mind Game
## Run
```
yourmachine $ vagrant up
yourmachine $ vagrant ssh
vagrant $ cd ~/project/project
vagrant $ ./start.sh
```
Go to localhost:8080 in a browser outside of vagrant

## Stop
```
vagrant $ cd ~/project/project
vagrant $ ./kill
```
# Purpose
Create an adaptation of [The Mind](https://boardgamegeek.com/boardgame/244992/mind) card game.

## How to play
Full rules in rules.html  
In a team of 2-4, place the cards in ascending order, WITHOUT COMMUNICATING.

### Scenario
You and a friend are sitting across from each other at a table.  
You both launch the website on your phone/computer.  
Player 1 has cards 37 and 99.  
Player 2 has cards 5 and 30.  
For the team to advance to the next level, player 2 must place cards 5 and 30, then player 1 must place cards 37 and 99.

# Languages
## Python
Using the Flask framework, Python handles the web server to handle http and REST calls.

## Go
Keeps track of the game state. Essentially acts as a forgetful database.

## Javascript
Handles client-side real time updates and form verification.

# Cross Language Communication
## Rabbit MQ
Communicates game data between the Go and Python server. Accomplished via a Remote Procedure Call method.

## SocketIO and REST
Handles communication between the client and the Python server. Game input, serving HTML pages, etc.

# Cool Things I Did
- Real time game play. When you play a card or join a lobby, your fellow players get instant feedback without having to referesh the page
- Individual rooms. You can run multiple games at once!
    - Make sure you launch in separate browsers (ex. Firefox and Chrome) or private/incognito windows
    - NOTE: Private/Incognito windows share cookies so no more than two games can run from a single browser
- It's mobile friendly :)
- You can kick players from the lobby
- Share a link to your lobby so your friends can join

# Screenshots
<img src="/uploads/fbdc8306641d0b5fd34102c81f6218e2/image.png"  width="800">
<br>
<img src="/uploads/73498defab287103bc5d96c09b1d19d1/Screenshot_20200816-234253_Chrome.jpg"  width="240">
<img src="/uploads/caaf6f67827561441d34d18d38548cf4/Screenshot_20200816-234347_Chrome.jpg "  width="240">  

# Reference/Disclaimer
Design language of project is very similar to my CMPT 276 group project. I did use some of the same boostrap classes.  
Fellow group members of that project are Oliver Xie and Adrian Nguyen
