The client is made with unity. Here is the github.
https://github.com/daniel840829/Tank-Online

The file structure:
- data
- entity
- hmap
- msg
- physic
this package is base on https://github.com/ianremmler/ode
use some struct to wrapper
  - World
  - Obj
  - ObjData
- rpctest
- service
  - Use gRPC to communicate with client
  -- I divide the message into __ part
  ---Error
  ---Sync Postion
  ---CallMathod
  ---Regist
  ---Login
  - use Game Manager to select three part of these
  --Error
  --Sync Postion
  --CallMathod
- storage
- timeCalibrate
- user
- util
- uuid
