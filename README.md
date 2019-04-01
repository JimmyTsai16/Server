### User API
need authenticated.
* GET /user <br>
Return current user information.
* POST /user <br>
Create new user.
* GET /user/{id} <br>
Return user information of id.

### Chat API
* GET /rooms <br />
Return rooms of current user created.
* WebSocket /chatws/{roomid}/{token} <br />
Connect the room.
* POST /room <br/>
Create new room.

### System Information API
* GET /multiple/{infoType}/{startDate}/{endDate} <br />
Return history system information of {infoType} from {startDate} to {endDate}.
* GET /host <br />
Return system information of host.
* Websocket /ws <br />
Return real-time system information from host continuously.
