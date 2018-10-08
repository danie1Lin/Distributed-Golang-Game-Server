let SessionLoad = 1
if &cp | set nocp | endif
let s:so_save = &so | let s:siso_save = &siso | set so=0 siso=0
let v:this_session=expand("<sfile>:p")
silent only
cd ~/go/src/github.com/daniel840829/gameServer
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
set shortmess=aoO
badd +8 ~/.bash_profile
badd +60 storage/storage.go
badd +1 restart
badd +44 main.go
badd +1 cluster-manager/manager.go
badd +1 Gopkg.toml
badd +18 Dockerfile
badd +91 game/session/room.go
badd +1 game.log
badd +1 entity/entity.go
badd +1 entity/attack.go
badd +1 entity/enemy.go
badd +1 msg/message.proto
badd +67 service/grpc.go
badd +45 session/session.go
badd +95 user/user.go
badd +1 ~/go/src/github.com/daniel840829/remoteProxy/remoteProxy.go
badd +1 data.proto
badd +7 service/ClientToAgent.go
badd +1 msg/message.pb.go
badd +1 uuid/uuid.go
badd +1 service/ClientToGame.go
badd +4 service/AgentToGame.go
badd +1 entity/room.go
badd +1 agent/AgentToGame.go
badd +68 agent/ClientToAgent.go
badd +341 agent/session/session.go
badd +107 agent/session/room.go
badd +4 game/entity/room.go
badd +1 agentServer/agent/AgentToGame.go
badd +1 game/agentToGame.go
badd +1 game/clientToGame.go
badd +1 game/session/session.go
badd +1 util/util.go
badd +16 agent/session/kubermanager.go
argglobal
silent! argdel *
$argadd entity/entity.go
edit agent/ClientToAgent.go
set splitbelow splitright
set nosplitbelow
set nosplitright
wincmd t
set winminheight=1 winheight=1 winminwidth=1 winwidth=1
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 4 - ((3 * winheight(0) + 18) / 37)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
4
normal! 08|
tabnext 1
if exists('s:wipebuf')
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=1 shortmess=aoO
set winminheight=1 winminwidth=1
let s:sx = expand("<sfile>:p:r")."x.vim"
if file_readable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &so = s:so_save | let &siso = s:siso_save
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
