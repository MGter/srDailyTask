#ifndef __COMMON_LOG_H_
#define __COMMON_LOG_H_

#include <map>
#include <string>
#include <fstream>
#include <iostream>
#include <sstream>
#include <chrono>
#include <ctime>
#include <iomanip>
#include "Tool.h"

extern std::map<std::string, std::string> g_configMap;

enum LogLevel {
    DEBUG,
    INFO,
    WARNING,
    ERROR
};

class Log {
public:
    static void initLogLevel();
    static void setLogLevel(LogLevel level);
    static void setLogPath(std::string& logPath);
    static void debug(const std::string& message, const char* file, int line);
    static void info(const std::string& message, const char* file, int line);
    static void warning(const std::string& message, const char* file, int line);
    static void error(const std::string& message, const char* file, int line);
private:
    void static createLogFile();
    void static closeLogFile();
    void static writeLog(const std::string& logLine);
private:
    static LogLevel        m_logLevel;
    static std::string     m_logPath;
    static std::ofstream   m_logFile;
};


#endif //__LOG_H_