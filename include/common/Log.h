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

// 这里就不要引用什么全局变量了init

enum LogLevel {
    DEBUG,
    INFO,
    WARNING,
    ERROR
};

class Log {
public:
    static void init(const LogLevel& log_level, const std::string& log_path);
    static void setLogLevel(const LogLevel& level);
    static void setLogPath(const std::string& logPath);
    static void debug(const std::string& message, const char* file, int line);
    static void info(const std::string& message, const char* file, int line);
    static void warning(const std::string& message, const char* file, int line);
    static void error(const std::string& message, const char* file, int line);
private:
    bool static createLogFile(const std::string& logDir, const std::string& logFilename);
    bool static closeLogFile();
    void static writeLog(const std::string& logLine);
private:
    static LogLevel        m_logLevel;
    static std::string     m_logPath;
    static std::ofstream   m_logFile;
    // 区分是否因时间更新而导致输出的log文件位置变化
    static std::string     m_lastLogDir;
    static std::string     m_lastLogFileName;
};


#endif //__LOG_H_