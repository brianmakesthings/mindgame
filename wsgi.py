from flask import Flask, render_template, request, redirect, url_for, session
from flask_socketio import SocketIO, send, emit, join_room, leave_room
from flask_talisman import Talisman
import engineio
import datetime
from markupsafe import escape
import uuid
import json
import re
import game


# csp = {
#     'default-src': [
#         '\'self\'',
#         '\'unsafe-inline\'',
#         'stackpath.bootstrapcdn.com',
#         'code.jquery.com',
#         'cdn.jsdelivr.net',
#         'unpkg.com',
#         'cdnjs.cloudflare.com'
#     ]
# }
app = Flask(__name__)
# Talisman(app, content_security_policy=csp)
app.secret_key = b"ikkO\xb8\xca\xec\xa8.\xb0|':\xee\xafM"
app.config['SECRET_KEY'] = b'\xdf\x18u\xdb-\xd1\xf0BBv\x1c\xbbf\xa8i\x9b'

# rpc = DemoRpcClient()
socketio = SocketIO(app)

# # https://stackoverflow.com/questions/40972805/python-capture-contents-inside-curly-braces
# def extractPlayerState(playerState):
#     regexPlayer = r"\{(.*?)\}"
#     regexCard = r"\[(.*?)\]"
#     matches = re.findall(regexPlayer, playerState)
    
#     for i in range(len(matches)):
#         matches[i] = matches[i].split(" ", 1)
#         matches[i][1] = re.search(regexCard, matches[i][1]).group(0)
#         matches[i][1] = matches[i][1][1:-1].split(" ")
#         if "" in matches[i][1]:
#             matches[i][1].remove("")
#         matches[i][1].sort(key=int)
#         # extrach card values
#     print(matches)
#     return matches


@socketio.on('userConnect')
def userConnect(json, methods=['GET', 'POST']):
    print('received my event: ' + str(json))
    room = str(json["roomId"])
    session['room'] = room
    participantData = game.getParticipants(room)
    join_room(room)
    emit('connectResponse', participantData, room=room)

@socketio.on('checkLobby')
def checkLobby(json, methods=['GET', 'POST']):
    room = str(json["roomId"])
    participantData = game.getParticipants(room)
    participants = participantData[0]
    if len(participants) < 2:
        emit('lobbyError', "Not Enough People")
    elif len(participants) > 4:
        emit('lobbyError', "Too Many People")
    else:
        emit('lobbyError')

@socketio.on('playCard')
def playCard(json, methods=['GET', 'POST']):
    print(json)
    room = str(json["roomId"])
    userId = str(json["userId"])
    card = str(json["card"])
    gameState = game.playCard(room, userId, card)
    playersState = [(player.id, player.cards) for player in gameState["playerState"]]
    emit('playCardResponse', (card, gameState['lives'], gameState['level'], playersState, gameState["moveResult"] ), room=room)

@socketio.on('kickUser')
def kickUser(json, methods=['GET', 'POST']):
    print(json)
    room = str(json["roomId"])
    userId = str(json["userId"])
    game.removeUser(int(room), userId)
    participantData = game.getParticipants(int(room))
    join_room(room)
    emit('connectResponse', participantData, room=room)

@socketio.on('restart')
def restart(json, methods=['Get', 'POST']):
    room = int(json["roomId"])
    game.restart(room)

# @socketio.on('disconnect')
# def on_leave():
#     print("SOMEONE DISCONNECTED")
#     room = session["room"]
#     leave_room(room)
#     emit('leaveResponse', room=room)

# From: https://www.bonser.dev/flask-session-timeout.html
@app.before_request
def before_request():
    session.permanent=True
    app.permanent_session_lifetime = datetime.timedelta(minutes=20)
    session.modified = True

@app.route('/')
def home():
    return render_template('home.html', redirect="/login")

@app.route('/rules')
def rules():
    return render_template('rules.html')

@app.route('/login', methods=["POST"])
def login():
    username = request.form['uname']
    username = re.sub(r'\W+', '', username)
    id = game.generateUserId(username)
    session['username'] = username
    session['id'] = id
    return redirect(url_for('join'))

@app.route('/invite/<int:roomNum>', methods=["GET", "POST"])
def invite(roomNum):
    if request.method=="POST":
        username = request.form['uname']
        username = re.sub(r'\W+', '', username)
        id = game.generateUserId(username)
        session['username'] = username
        session['id'] = id
        session['id'] = response['id']
    else:
        if 'id' not in session:
            return render_template('home.html', redirect="/invite/%d" % roomNum)
            
    if (roomNum == ""):
        return "Please enter a valid room number"
        
    roomNum = str(roomNum)
    response = game.joinRoom(session['id'], roomNum)
    if response != "Success":
        if response == "Player already in room":
            return redirect('/lobby/'+roomNum)
        else:
            return "error: " + response
    return redirect('/lobby/' + roomNum)


@app.route('/join')
def join():
    if 'username' not in session:
        return redirect('/')
    print(session['id'])
    return render_template('join.html', uname=escape(session['username']))
    

@app.route('/createRoom', methods=["POST"])
def create():
    id = session['id']
    roomId = game.createRoom(id)
    print(roomId)
    return redirect('/lobby/'+str(roomId))

@app.route('/searchRoom', methods=["POST"])
def searchRoom():
    roomNum = request.form.get("roomNum")
    if (roomNum == ""):
        return "Please enter a valid room number"
    roomNum = str(roomNum)
    response = game.joinRoom(session['id'], roomNum)
    if response != "Success":
        if response == "Player already in room":
            return redirect('/lobby/'+roomNum)
        else:
            return "error: " + response
    return redirect('/lobby/' + roomNum)

@app.route('/lobby/<int:roomNum>')
def lobby(roomNum):
    if 'username' not in session:
        return redirect('/')
    participantData = game.getParticipants(roomNum)
    print(participantData)
    return render_template('lobby.html', roomNum=escape(str(roomNum)), userId=session['id'], participants=participantData[0], owner=participantData[1], error="")


@app.route('/startgame/<int:roomNum>', methods=["POST"])
def start_game(roomNum):
    game.startGame(roomNum)
    socketio.emit('startGame', str(roomNum), room=str(roomNum))
    return redirect('/game/%d' % roomNum)

@app.route('/game/<int:roomNum>')
def load_game(roomNum):
    if 'username' not in session:
        return redirect('/')
    gameData = game.getGameState(roomNum)
    # playerState = extractPlayerState(response['playerState'])
    print("Player State")
    print(gameData["playerState"])
    return render_template('game.html', userId = session['id'], roomNum=escape(str(roomNum)), 
                            lives = gameData['lives'], level = gameData['level'], 
                            playerState = gameData["playerState"] )




if __name__ == '__main__':
    socketio.run(app, host='0.0.0.0', port=56247)
    