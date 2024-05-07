#include "common/Log.h"
#include <iostream>
#include <ctime>

LogLevel Log::m_logLevel = DEBUG;
std::string Log::m_logPath = "./log/";
std::ofstream Log::m_logFile;

void Log::init(const LogLevel& log_level, const std::string& log_path){
    setLogLevel(log_level);
    setLogPath(log_path);
    m_lastLogDir = "";
    m_lastLogFileName = "";
}

void Log::setLogLevel(const LogLevel& level) {
    m_logLevel = level;
}

void Log::setLogPath(const std::string& logPath){
    m_logPath = logPath;
}


void Log::debug(const std::string& message, const char* file, int line) {
    if (m_logLevel <= DEBUG) {
        std::string timeString;
        getCurrentTime(timeString);
        std::string logLine = "[" + timeString + "][DEBUG][" + file + ":" + std::to_string(line) + "] " + message + "\n"; 
        writeLog(logLine);
    }
}

void Log::info(const std::string& message, const char* file, int line) {
    if (m_logLevel <= INFO) {
        std::string timeString;
        getCurrentTime(timeString);
        std::string logLine = "[" + timeString + "][INFO][" + file + ":" + std::to_string(line) + "] " + message + "\n"; 
        writeLog(logLine);
    }
}

void Log::warning(const std::string& message, const char* file, int line) {
    if (m_logLevel <= WARNING) {
        std::string timeString;
        getCurrentTime(timeString);
        std::string logLine = "[" + timeString + "][WARNING][" + file + ":" + std::to_string(line) + "] " + message + "\n"; 
        writeLog(logLine);
    }
}

void Log::error(const std::string& message, const char* file, int line) {
    if (m_logLevel <= ERROR) {
        std::string timeString;
        getCurrentTime(timeString);
        std::string logLine = "[" + timeString + "][ERROR][" + file + ":" + std::to_string(line) + "] " + message + "\n"; 
        writeLog(logLine);
    }
}


bool Log::createLogFile(const std::string& logDir, const std::string& logFilename) {
    if (!directoryExists(logDir)) {
        createDirectory(logDir);
    }

    std::string logFilePath = logDir + "/" + logFilename;

    m_logFile.open(logFilePath, std::ios_base::app);
    if(!m_logFile.is_open()){
        std::cout << "[" << __FILE__ << "][" << __LINE__ << "]" << "Func: createLogFile: Failed to open the file: " << logFilePath << std::endl;
        return false;
    }
    else{
        return true;
    }
}

bool Log::closeLogFile() {
    if (m_logFile.is_open()) {
        m_logFile.close();
    }
    else{
        std::cout << "[" << __FILE__ << "][" << __LINE__ << "]" << "Func: closeLogFile: No logfile open!" << std::endl;
    }
}


// 写入日志
void Log::writeLog(const std::string& logLine) {
    std::time_t now = std::chrono::system_clock::to_time_t(std::chrono::system_clock::now());
    std::tm* timeInfo = std::localtime(&now);
    std::ostringstream oss;

    // 获取path地址
    oss << std::put_time(timeInfo, "%Y%m");
    std::string logDir = m_logPath + "/" + oss.str();
    // 获取log名称
    oss.clear();
    oss << std::put_time(timeInfo, "%Y%m%d");
    std::string logFilename = oss.str() + ".log";

    // 如果尚未打开文件，尝试打开
    if(!m_logFile.is_open()){
        if(!Log::createLogFile(logDir, logFilename)){
            std::cout << "[" << __FILE__ << "][" << __LINE__ << "]" << "Func:writeLog: Failed to open the log file!" << std::endl;
            return;
        }
    }

    // 文件位置变化
    if(m_lastLogDir != logDir ||  m_lastLogFileName != logFilename){
        Log::closeLogFile();
        if(!Log::createLogFile(logDir, logFilename)){
            std::cout << "[" << __FILE__ << "][" << __LINE__ << "]" << "Func:writeLog: Failed to open the log file!" << std::endl;
            return;
        }
        m_lastLogDir = logDir;
        m_lastLogFileName = logFilename;
    }

    // 写入文件
    m_logFile << logLine << std::endl;
}