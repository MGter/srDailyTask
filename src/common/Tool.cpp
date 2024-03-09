#include "Tool.h"


// 初始化conf文件
void initConf(const std::string& conf_path) {
    std::ifstream configFile(conf_path);
    std::string line;
    while (std::getline(configFile, line)) {
        // 解析配置项，假设配置项格式为 key=value
        std::istringstream iss(line);
        std::string key, value;
        if (std::getline(iss, key, '=') && std::getline(iss, value)) {
            // 存储到全局的map中
            g_configMap[key] = value;
        }
    }
    configFile.close();
}

void getCurrentTime(std::string& timeString) {
    time_t now = time(0);
    struct tm* timeinfo;
    char buffer[80];

    timeinfo = localtime(&now);
    strftime(buffer, sizeof(buffer), "%Y-%m-%d %H:%M:%S", timeinfo);
    timeString = buffer;
}
