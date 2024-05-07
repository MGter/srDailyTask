#ifndef __SR_DAILY_TASK_
#define __SR_DAILY_TASK_



class srDailyTask{
public:
    enum Status{
        init,
        running,
        stop,
        err,
        unknown
    };

    srDailyTask();
    ~srDailyTask();
    bool start();
    bool stop();
    bool addTask();
    bool delTask();
    bool findTask();
    bool changeTask();
    static srDailyTask* getInstance();

private:
    Status ThreadRunning();

public:
    static srDailyTask* g_srDailyTask;

private:
    bool m_started;
};



#endif //__SR_DAILY_TASK_