from flask import Flask, render_template, request, redirect, url_for, session
from flask_socketio import SocketIO, send, emit, join_room, leave_room
import engineio
import datetime
from markupsafe import escape
import pika
import uuid
import json
import re

# python3 amqp_client.py
# based on https://www.rabbitmq.com/tutorials/tutorial-six-python.html
class DemoRpcClient(object):
    def __init__(self):
        self.connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost', heartbeat=12000))
        self.channel = self.connection.channel()

        result = self.channel.queue_declare(queue='', exclusive=True)
        self.callback_queue = result.method.queue

        self.channel.basic_consume(
            queue=self.callback_queue,
            on_message_callback=self.on_response,
            auto_ack=True)

    def on_response(self, ch, method, props, body):
        if self.corr_id == props.correlation_id:
            self.response = json.loads(body.decode('utf-8'))

    def call(self, arg):
        self.response = None
        self.corr_id = str(uuid.uuid4())
        body = json.dumps(arg).encode('utf8')
        
        self.channel.basic_publish(
            exchange='',
            routing_key='rpc_queue',
            properties=pika.BasicProperties(reply_to=self.callback_queue, correlation_id=self.corr_id),
            body=body)
        while self.response is None:
            self.connection.process_data_events()
        return self.response

app = Flask(__name__)
app.secret_key = b"ikkO\xb8\xca\xec\xa8.\xb0|':\xee\xafM"
app.config['SECRET_KEY'] = b'\xdf\x18u\xdb-\xd1\xf0BBv\x1c\xbbf\xa8i\x9b'

rpc = DemoRpcClient()
socketio = SocketIO(app)

# https://stackoverflow.com/questions/40972805/python-capture-contents-inside-curly-braces
def extractPlayerState(playerState):
    regexPlayer = r"\{(.*?)\}"
    regexCard = r"\[(.*?)\]"
    matches = re.findall(regexPlayer, playerState)
    
    for i in range(len(matches)):
        matches[i] = matches[i].split(" ", 1)
        matches[i][1] = re.search(regexCard, matches[i][1]).group(0)
        matches[i][1] = matches[i][1][1:-1].split(" ")
        if "" in matches[i][1]:
            matches[i][1].remove("")
        matches[i][1].sort(key=int)
        # extrach card values
    print(matches)
    return matches
    


@socketio.on('userConnect')
def userConnect(json, methods=['GET', 'POST']):
    print('received my event: ' + str(json))
    room = str(json["roomId"])
    session['room'] = room
    response = rpc.call({'request' : 'getParticipants', 'roomId' : room})
    if "error" in response:
        return "error: " + response["error"]
    participants = response["participants"][1:-1].split(" ")
    join_room(room)
    emit('connectResponse', (participants, response["owner"]), room=room)

@socketio.on('checkLobby')
def checkLobby(json, methods=['GET', 'POST']):
    room = str(json["roomId"])
    response = rpc.call({'request' : 'getParticipants', 'roomId' : room})
    if "error" in response:
        return "error: " + response["error"]
    participants = response["participants"][1:-2].split(" ")
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
    response = rpc.call({'request' : 'playCard', 'roomId' : room, 'userId' : userId, 'card' : card})
    print(response)
    playerState = extractPlayerState(response["playerState"])
    emit('playCardResponse', (card, response['lives'], response['level'], playerState, response["moveResult"] ), room=room)

@socketio.on('kickUser')
def kickUser(json, methods=['GET', 'POST']):
    print(json)
    room = str(json["roomId"])
    userId = str(json["userId"])
    response = rpc.call({'request' : 'removeUser', 'roomId' : room, 'userId': userId})
    print(response)
    response = rpc.call({'request' : 'getParticipants', 'roomId' : room})
    if "error" in response:
        return "error: " + response["error"]
    participants = response["participants"][1:-1].split(" ")
    join_room(room)
    emit('connectResponse', (participants, response["owner"]), room=room)

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
    response = rpc.call({'request' : 'generateUserId', 'string': username})
    if ("error" in response):
        return "error: " + response["error"]
    session['username'] = username
    session['id'] = response['id']
    return redirect(url_for('join'))

@app.route('/invite/<int:roomNum>', methods=["GET", "POST"])
def invite(roomNum):
    if request.method=="POST":
        username = request.form['uname']
        username = re.sub(r'\W+', '', username)
        response = rpc.call({'request' : 'generateUserId', 'string': username})
        if ("error" in response):
            return "error: " + response["error"]
        session['username'] = username
        session['id'] = response['id']
    else:
        if 'id' not in session:
            return render_template('home.html', redirect="/invite/%d" % roomNum)
            
    if (roomNum == ""):
        return "Please enter a valid room number"
    roomNum = str(roomNum)
    response = rpc.call({'request' : 'joinRoom', 'userId' : session['id'], 'roomId' : roomNum})
    if "error" in response:
        if response["error"] == "user already in room":
            return redirect('/lobby/'+roomNum)
        else:
            return "error: " + response["error"]
    return redirect('/lobby/' + roomNum)


@app.route('/join')
def join():
    if 'username' not in session:
        return redirect('/')

    print(session['username'])
    return render_template('join.html', uname=escape(session['username']))
    

@app.route('/createRoom', methods=["POST"])
def create():
    id = session['id']
    response = rpc.call({'request' : 'createRoom', 'ownerId' : id})
    if ("error" in response):
        return "error: " + response["error"]
    print(response['roomId'])
    return redirect('/lobby/'+response['roomId'])

@app.route('/searchRoom', methods=["POST"])
def searchRoom():
    roomNum = request.form.get("roomNum")
    if (roomNum == ""):
        return "Please enter a valid room number"
    roomNum = str(roomNum)
    response = rpc.call({'request' : 'joinRoom', 'userId' : session['id'], 'roomId' : roomNum})
    if "error" in response:
        if response["error"] == "user already in room":
            return redirect('/lobby/'+roomNum)
        else:
            return "error: " + response["error"]
    return redirect('/lobby/' + roomNum)

@app.route('/lobby/<int:roomNum>')
def lobby(roomNum):
    if 'username' not in session:
        return redirect('/')
    response = rpc.call({'request' : 'getParticipants', 'roomId' : str(roomNum)})
    if "error" in response:
        return "error: " + response["error"]
    print(response["participants"])
    participants = response["participants"][1:-1].split(" ")
    print(participants)
    # isOwner = (session["id"] == response["owner"])
    return render_template('lobby.html', roomNum=escape(str(roomNum)), userId=session['id'], participants=participants, owner=response["owner"], error="")


@app.route('/startgame/<int:roomNum>', methods=["POST"])
def start_game(roomNum):
    response = rpc.call({'request' : 'startGame', 'roomId' : str(roomNum)})
    if "error" in response:
        return "error: " + response["error"]
    socketio.emit('startGame', str(roomNum) ,room=str(roomNum))
    return redirect('/game/%d' % roomNum)

@app.route('/game/<int:roomNum>')
def load_game(roomNum):
    if 'username' not in session:
        return redirect('/')
    response = rpc.call({'request' : 'getGameState', 'roomId' : str(roomNum)})
    if "error" in response:
        return "error: " + response["error"]
    # should add a check for unwanted guests
    playerState = extractPlayerState(response['playerState'])
    print(playerState)
    return render_template('game.html', userId = session['id'], roomNum=escape(str(roomNum)), 
                            lives = response['lives'], level = response['level'], 
                            playerState = playerState )




if __name__ == '__main__':
    socketio.run(app, host='0.0.0.0', port=5000)
    