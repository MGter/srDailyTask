g++ -c src/main.cpp -o obj/main.o -Iinclude -Iinclude/common -Ideps -Ideps/mysql_server/include -Ideps/rapidjson
g++ -c src/common/JsonParse.cpp -o obj/JsonParse.o -Iinclude -Iinclude/common -Ideps -Ideps/mysql_server/include -Ideps/rapidjson
g++ -c src/common/Log.cpp -o obj/Log.o -Iinclude -Iinclude/common -Ideps -Ideps/mysql_server/include -Ideps/rapidjson
g++ -c src/common/Tool.cpp -o obj/Tool.o -Iinclude -Iinclude/common -Ideps -Ideps/mysql_server/include -Ideps/rapidjson
g++ -c src/srDailyTask.cpp -o obj/srDailyTask.o -Iinclude -Iinclude/common -Ideps -Ideps/mysql_server/include -Ideps/rapidjson
g++ -c sample/test.cpp -o obj/test.o -Iinclude -Iinclude/common -Ideps -Ideps/mysql_server/include -Ideps/rapidjson

g++ obj/main.o obj/JsonParse.o obj/Log.o obj/Tool.o obj/srDailyTask.o obj/test.o -o srDailyTask