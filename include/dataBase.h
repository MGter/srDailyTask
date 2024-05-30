#ifndef __DATA_BASE_H__
#define __DATA_BASE_H__

#include <string>

#include "common/JsonParse.h"
#include "common/Log.h" 
#include "common/Tool.h"
#include "module/point.h"
#include "module/task.h"
#include "module/user.h"
#include "module/wallet.h"

class DataBase{
public:
    enum Status{
        start,
        stop,
        error,
        unknown
    };

    DataBase();
    ~DataBase();

    bool init(const std::string name, const std::string passwd);
    bool connect();

    bool addMessage();
    bool deleteMessage();
    bool findMessage();
    bool changeMessage();

    DataBase* getDBInstance();

private:
    std::string m_dbName;
    std::string m_dbPassWD;
    static DataBase* g_DB;
};

#endif // __DATA_BASE_H__