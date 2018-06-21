<h1>Tank Game Server</h1>
  The client is made with unity. Here is the github.
  https://github.com/daniel840829/Tank-Online

<h2>How to Use : </h2>
<h3>client endpoint:</h3>
  download client code<a href="https://github.com/daniel840829/Tank-Online">Here</a>
<h3>server endpoint</h3>
  physic package is base on https://github.com/ianremmler/ode Open Dynamic Engine<br>
  use some struct to wrap ianremmler's ode golang library<br>
  Please read his repository to build ODE's share library<br>
  But I made some change in order to let the collide callback to deal with geom is space
  please install https://github.com/daniel840829/ode

<h2>The file structure:</h2>
<ol>
  
  <li>
  entity
    <ol>
      <li>GameManager</li>
      <li>Entity</li>
      <li>Room</li>
    </ol>
  </li>
    
  <li>hmap</li>
  
  <li>msg
  <br>  - use protobuf to define package and rpc interface
  </li>
  
  <li>physic<br>
  this package is base on https://github.com/ianremmler/ode Open Dynamic Engine<br>
  use some struct to wrap ianremmler's ode golang library<br>
  Please read his repository to build ODE's share library<br>
    <ul>
      <li>World</li>
      <li>Obj</li>
      <li>ObjData</li>
     </ul>
  </li>
  <li>rpctest</li>
  
  <li>service
  <br>implemente gRPC's Service defined in msg package to communicate with client
  <br>I divide the message into five parts
    <ul>
      <li>Error</li>
      <li>Sync Postion</li>
      <li>CallMathod</li>
      <li>Regist</li>
      <li>Login</li>
     </ul>
  use Game Manager to select three parts of these
    <ul>
      <li>Error</li>
      <li>Sync Postion</li>
      <li>CallMathod</li>
    </ul>
   </li>
   
   <li>
    uuid
    <br>generate different IDs of objects that can be call with reflection 
   </li>

  <li>user
    <ul>
      <li>UserManager</li>
      <li>User<li>
  TODO: There should be a session manager to cache user infomation and state
  <li>
    storage
    <br>Use mogdb to storage user infomation
  </li>
      
  <li>timeCalibrate</li>
      
  <li>
    util
    <br>some tools
  </br>
  
  <li>
    data
    <br>use to load some model data e.g. character or map
  </li>
</ol>
