#ifndef __COMMON_TOOL_H_
#define __COMMON_TOOL_H_

#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <map>
#include <sys/stat.h>
#include <sys/types.h>

extern std::map<std::string, std::string> g_configMap;
extern std::string g_logPath;


// 读取conf文件，进行初始化
void initConf(const std::string& conf_path);

// 获取当前时间，存入timeString
void getCurrentTime(std::string& timeString);

// 确认文件是否存在
bool directoryExists(const std::string& path);

// 创建相应的文件夹
void createDirectory(const std::string& path);









#endif  //__COMMON_TOOL_H_