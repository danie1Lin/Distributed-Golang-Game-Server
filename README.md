## Golang Distributed Game Server

## Motivation
At first, I just want to learn Golang.I started to think about which is the best way?
Because the concurrency mechanisms of Golang is very powerful, I choose online game to verify if I can use Golang to make a efficient game server.For me, this is the first time I make this such hard project. I have to learn Unity, Golang, C# At a time. I am glad that I still have full passion to this project and I have never given up.

[![CircleCI](https://circleci.com/gh/daniel840829/gameServer/tree/kubernete-intergration.svg?style=svg)](https://circleci.com/gh/daniel840829/gameServer/tree/kubernete-intergration)

## Tech/framework used
- golang
- gRPC
- Kubernetes
- protobuf
## Features
- CrossPlatform - Message packet use protobuf which is light, fast, crossplatform. 
- Autoscaling - controller is written with go-client ,you can wirte the strategy to autoscale dedicated game server by your own.
- Lightweight - The image of dedicated game server is less than 40MB.
## Ａrchitecture
- Agent server : 
  - match players to join other player's room or create own room
  - control the amount of gameplay server and load balancing. when the amountof a gameplay server's connections exceed maxium connections it should have, agent will create a new pod run gameplay server.
- Gameplay Server :
  - After players are matched successfully ,these players will get the gameplay server's ip and token,and player can start to play.

## Installation
- Server
  - Install kubernete
  - Create cluster
  - go run main.go --type=agent
- Client
  you can run in Unity Editor by open the Prestart.scene as first scene.
  or you can just download from <a href="https://github.com/daniel840829/Tank-Online/tree/GameStateSystem/Build">build</a>
  Whole project download from <a href="https://github.com/daniel840829/Tank-Online">here</a>
## How to use?
If you want to make your own game by modifying this project, I am pleasured.
You can throught these step to make it work.
### New Function?

### Modify The msg/message.proto
Change the GameFrame message in proto buff
### Creat Your Game Logic


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
A short snippet describing the license (MIT, Apache etc)

MIT © [daniel840829]()