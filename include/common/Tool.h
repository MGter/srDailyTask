#ifndef __COMMON_TOOL_H_
#define __COMMON_TOOL_H_

#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <map>


extern std::map<std::string, std::string> g_configMap;
extern std::string g_logPath;


// 读取conf文件，进行初始化
void initConf(std::string& conf_path);

// 获取当前时间，存入timeString
void getCurrentTime(std::string& timeString);












#endif  //__COMMON_TOOL_H_