#include <iostream>
#include <string>
#include <chrono>


typedef enum CircleMode{
    onece,
    weekly,
    workday,
    weekend,
    custom
}CircleMode;


typedef struct Task{
    uint64_t            taskID;             // task的唯一id
    uint64_t            userID;             // 使用人的唯一id

    std::time_t         createTime;         // task的创建时间
    std::time_t         changeTime;         // task的修改时间

    CircleMode          circleMode;         // 循环模式
    std::string         message;            // task消息内容

    bool                isTired;            // task是否过期
    bool                beSilence;          // 静音，暂不提醒
}Task;