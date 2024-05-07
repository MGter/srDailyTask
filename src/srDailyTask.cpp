#include "srDailyTask.h"


srDailyTask::srDailyTask(){
    m_started = false;
}


srDailyTask::~srDailyTask(){
    ;
}

srDailyTask* srDailyTask::getInstance(){
    if(!g_srDailyTask){
        g_srDailyTask = new srDailyTask;
    }
    return g_srDailyTask;
}


bool srDailyTask::start(){
    ;
}

bool srDailyTask::stop(){
    ;
}

bool srDailyTask::addTask(){
    ;
}

bool srDailyTask::delTask(){
    ;
}

bool srDailyTask::findTask(){
    ;
}

bool srDailyTask::changeTask(){
    ;
}