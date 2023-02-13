curl -X POST http://192.168.28.185:7777/robot/create -d '{"pool":"41ee7f5f79e8cac83c7942120c088630ec5ef47fa0","token":"TRX","start_num":5,"num_of_bets":7,"odd_chips":[21,41,81,161,321,641,1281],"even_chips":[20,40,60,80,160,320,640,1280],"take_profit":50000,"stop_loss":1000}'
curl  http://192.168.28.185:7777/robot/run?id=1623169938061529088
curl  http://192.168.28.185:7777/robot/stop?id=1623169938061529088
curl  http://192.168.28.185:7777/robot/stop?id=1623169938061529088