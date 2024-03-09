#include "Log.h"
#include <iostream>
#include <ctime>


void Log::initLogLevel() {
    if (g_configMap.find("LogLevel") != g_configMap.end()) {
        std::string levelStr = g_configMap["LogLevel"];
        if (levelStr == "0") {
            Log::setLogLevel(DEBUG);
        } else if (levelStr == "1") {
            Log::setLogLevel(INFO);
        } else if (levelStr == "2") {
            Log::setLogLevel(WARNING);
        } else if (levelStr == "3") {
            Log::setLogLevel(ERROR);
        }
        else{
            Log::setLogLevel(DEBUG);
        }
    }
}


void Log::setLogLevel(LogLevel level) {
    m_logLevel = level;
}

void Log::setLogPath(std::string& logPath){
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
        std::cout << "[" << timeString << "][INFO][" << file << ":" << line << "] " << message << std::endl;
    }
}

void Log::warning(const std::string& message, const char* file, int line) {
    if (m_logLevel <= WARNING) {
        std::string timeString;
        getCurrentTime(timeString);
        std::cout << "[" << timeString << "][WARNING][" << file << ":" << line << "] " << message << std::endl;
    }
}

void Log::error(const std::string& message, const char* file, int line) {
    if (m_logLevel <= ERROR) {
        std::string timeString;
        getCurrentTime(timeString);
        std::cerr << "[" << timeString << "][ERROR][" << file << ":" << line << "] " << message << std::endl;
    }
}


void Log::createLogFile() {
    std::time_t now = std::chrono::system_clock::to_time_t(std::chrono::system_clock::now());
    std::tm* timeInfo = std::localtime(&now);
    std::ostringstream oss;
    oss << std::put_time(timeInfo, "%Y%m");
    std::string logDir = m_logPath + "/" + oss.str();
    std::string logFilename = logDir + "/" + oss.str() + ".log";

    if (!std::filesystem::exists(logDir)) {
        std::filesystem::create_directories(logDir);
    }

    m_logFile.open(logFilename, std::ios_base::app);
}