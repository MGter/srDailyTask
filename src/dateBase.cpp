#include "dataBase.h"

DataBase::DataBase(){
    ;
}


DataBase::~DataBase(){

}

bool DataBase::init(const std::string name, const std::string passwd){
    m_dbName = name;
    m_dbPassWD = passwd;
    if(!connect()){
        Log::error(__FILE__, __LINE__, "Failed to connect to the database, User: %s", m_dbName.c_str());
        return false;
    }
    Log::info(__FILE__, __LINE__, "Success to connect to the database, User: %s", m_dbName.c_str());
    return true;
}

DataBase* DataBase::getDBInstance(){
    if(!g_DB){
        g_DB = new DataBase;
    }
    return g_DB;
}

bool DataBase::connect(){
    ;
}


bool DataBase::addMessage(){
    
}

bool DataBase::deleteMessage(){

}

bool DataBase::findMessage(){

}

bool DataBase::changeMessage(){

}
