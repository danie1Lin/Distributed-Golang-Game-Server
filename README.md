## Golang Distributed Game Server
[![CircleCI](https://circleci.com/gh/daniel840829/gameServer/tree/kubernete-intergration.svg?style=svg)](https://circleci.com/gh/daniel840829/gameServer/tree/kubernete-intergration)
## Motivation
At first, I just want to learn Golang.I started to think about which is the best way?
Because the concurrency mechanisms of Golang is very powerful, I choose online game to verify if I can use Golang to make a efficient game server.For me, this is the first time I make this such hard project. I have to learn Unity, Golang, C# At a time. I am glad that I still have full passion to this project and I never give up.
## Tech/framework used
- golang
- gRPC
- Kubernetes
- Protobuf
## Features
- CrossPlatform - Message packet use protobuf which is light, fast, crossplatform. 
- Autoscaling - controller is written with go-client ,you can wirte the strategy to autoscale dedicated game server by your own.
- Lightweight - The image of dedicated game server is less than 40MB.
## Architecture
- Agent server : 
  - match players to join other player's room or create own room
  - control the amount of gameplay server and load balancing. when the amountof a gameplay server's connections exceed maxium connections it should have, agent will create a new pod run gameplay server.
- Gameplay Server :
  - After players are matched successfully ,these players will get the gameplay server's ip and token,and player can start to play.
  ![](https://github.com/daniel840829/gameServer/blob/kubernete-intergration/Golang%20game%20server%20architecture.png?raw=true)
## Installation
- Server
  - Install kubernete
  - Create cluster
  - go run main.go --type=agent
- Client
  - You can run in Unity Editor by open the Prestart.scene as first scene.
  - If you want to test with mutliplayer you can try build Andorid apk because it is likely to be builded successfully.
  ![](https://media.giphy.com/media/ftdlle6pOE6Y8w5bho/giphy.gif)
## How to use?
If you want to make your own game by modifying this project, I am pleasured.
You can throught these step to make it work.
### Modify The msg/message.proto
1. Change the GameFrame message in proto buff.
2. Run './update.sh' in 'msg/'
3. unzip message.zip under 'Asset/gameServer/proto' in the Unity Project.
### Create Your Game Logic
1. add your handler to <a href="https://github.com/daniel840829/gameServer/blob/a218213609e8857f84ffa5516c412922ef9cd4c1/game/session/room.go#L157">"gameServer/game/session/room.go": func (r *Room) Run()</a>
2. modify the <a href="https://github.com/daniel840829/Tank-Online/blob/87be8962024241dff4d8f3de1809fe4ef60f0848/Assets/Scripts/Entity/EntityManager.cs#L188">UpdateFrame</a> fuction to handle packets design by yourself. then,Data Flow to Entity to render the change of entity's properties.
## The file structure:
<ol>
  <li>agent
    <ol>
      <li>session
        <ol>
          <li>room</li>
          <li>session</li>
          <li>kubernete</li>
        </ol>
      </li>
    </ol>
  </li>
  <li>game
    <ol>
      <li>session
        <ol>
          <li>room</li>
          <li>session</li>
        </ol>
      </li>
    </ol>
  </li>
  <li>msg
  <br>  - use protobuf to define package and rpc interface
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
</ol>

## Contribute
## Credits
## License



MIT © [daniel840829]()
