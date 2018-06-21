<h1>Tank Game Server</h1>
  The client is made with unity. Here is the github.
  https://github.com/daniel840829/Tank-Online

<h2>How to Use : </h2>
<h3>client endpoint:</h3>
  download client code<a href="https://github.com/daniel840829/Tank-Online">Here</a>
<h3>server endpoint<h3>
  physic package is base on https://github.com/ianremmler/ode Open Dynamic Engine<br>
  use some struct to wrap ianremmler's ode golang library<br>
  Please read his repository to build ODE's share library<br>
  But I made some change in order to let the collide callback to deal with geom is space
  please install https://github.com/daniel840829/ode

<h2>The file structure:</h2>
- entity
  - GameManager
  - Entity
  - Room
- hmap
- msg
  - use protobuf to define package and rpc interface
- physic<br>
  this package is base on https://github.com/ianremmler/ode Open Dynamic Engine<br>
  use some struct to wrap ianremmler's ode golang library<br>
  Please read his repository to build ODE's share library<br>
   - World
   - Obj
   - ObjData
- rpctest
- service
  <br>implemente gRPC's Service defined in msg package to communicate with client
  - I divide the message into five parts
     - Error
     - Sync Postion
     - CallMathod
     - Regist
     - Login
  - use Game Manager to select three parts of these
    - Error
    - Sync Postion
    - CallMathod
- uuid
  <br>generate different IDs of objects that can be call with reflection 
- user
  - UserManager
  - User
  TODO: There should be a session manager to cache user infomation and state
- storage
  - Use mogdb to storage user infomation
- timeCalibrate
- util
  - some tools
- data
  - use to load some model data e.g. character or map
