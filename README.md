# Golang Distributed Game Server
[![Go Report Card](https://goreportcard.com/badge/github.com/daniel840829/Distributed-Golang-Game-Server)](https://goreportcard.com/report/github.com/daniel840829/Distributed-Golang-Game-Server)
[![CircleCI](https://circleci.com/gh/daniel840829/Distributed-Golang-Game-Server.svg?style=svg)](https://circleci.com/gh/daniel840829/Distributed-Golang-Game-Server)
[![Gitter](https://badges.gitter.im/daniel840829/Distributed-Golang-Game-Server.svg)](https://gitter.im/daniel840829/Distributed-Golang-Game-Server?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
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
### Agent server : 
  - match players to join other player's room or create own room
  - control the amount of gameplay server and load balancing. when the amountof a gameplay server's connections exceed maxium connections it should have, agent will create a new pod run gameplay server.
### Gameplay Server :
  - After players are matched successfully ,these players will get the gameplay server's ip and token,and player can start to play.
  ![](https://github.com/daniel840829/gameServer/blob/kubernete-intergration/Golang%20Game%20ServerArchitecture.png?raw=true)
### Packet Validating 
In the branch master, I use ODE to simulate the physics on server.It is the most safe way to keep game fair.However, I found the memory server use is too much for me, Because I don't have money to maintance the server. So I started to design a way to let client validate packet and simulate physics separatly to reduce the heavy load on the server. I just complete the entities can attack each others so far. I will start to design aftewards:
  - The validation part preventing form players cheating 
  - The interface connecting a physics simulator
## Installation
- Server :
Two Way to run server:
  - Run distributed server using Kubernete cluster
    1. Use Kops to install kubernete on AWS
    2. Create cluster
    3. Install Mongodb 
    4. edit setupEnv.sh with your setting and bash setupEnv.sh
    4. ```go run main.go --type=agent``` on your local machine (Must on Where you install Kops)
  - Run Standalone Server on local machine
    1. Install Mongodb 
    3. edit setupEnv.sh( DONT_USE_KUBE = true )with your setting and bash setupEnv.sh
    2. ```go run main.go --type=standalone``` on your local machine
- Client :
  - [Download this project](https://github.com/daniel840829/Tank-Online) 
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
- agent
  - session
    - room.go
    - session.go
    - kubernetes.go
- game
  - room.go
  - session.go
  - kubernetes.go
- msg Use Protobuf to define package and RPC service interface
- uuid generate different IDs of objects that can be call with reflection 
- user
  - UserManager
  - User
- storage Use mongoDB to storage user infomation
## Support me
 <a href="https://www.buymeacoffee.com/yEKnuC6" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/yellow_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>
