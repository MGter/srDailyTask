#include <signal.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <execinfo.h>
#include <iostream>
#include <map>

#include "JsonParse.h"
#include "Log.h"
#include "Tool.h"

#include "srDailyTask.h"

// 设置全局变量g_configMap，储存config.conf,设置日志级别等内容
std::map<std::string, std::string> g_configMap;
std::string g_logPath;

void SigIntHandler(int const a)
{
	signal(SIGINT, SigIntHandler);
    // 在这里发送停止的信号
	// xStreamDeviceTool::GetInstance()->m_bStop = 1;
}


static void WidebrightSegvHandler(int signum)
{
#define TRACK_SIZE 100
	void *array[TRACK_SIZE];
	size_t size;
	time_t crash_time;
	struct tm* time_info;
	int fd = 0;
	char file_name[128] = { 0 };

	/* 还原默认的信号处理handler */
	signal(signum, SIG_DFL);
	time(&crash_time);

	sprintf(file_name, "crash-%d-%d-%lld.dat", signum, getpid(), crash_time);
	fd = open(file_name, O_RDWR | O_CREAT | O_TRUNC);
	if (fd > 0)
	{
		size = backtrace(array, TRACK_SIZE);
		backtrace_symbols_fd(array, size, fd);
		close(fd);
	}

	_exit(1);
}


int main(int argc, char* argv[]){
	Log::info(__FILE__, __LINE__, "srDailyTask starts");

	signal(SIGINT, SigIntHandler);
	signal(SIGSEGV, WidebrightSegvHandler);
	signal(SIGABRT, WidebrightSegvHandler);
	signal(SIGILL, WidebrightSegvHandler);
	signal(SIGFPE, WidebrightSegvHandler);
	signal(SIGPIPE, SIG_IGN);
	sigset_t signal_mask;
	sigemptyset(&signal_mask);
	sigaddset(&signal_mask, SIGPIPE);
	int rc = pthread_sigmask(SIG_BLOCK, &signal_mask, NULL);
	if (rc != 0)
	{
		printf("block sigpipe error\n");
	}

	if(!srDailyTask::getInstance()->init()){
		Log::error(__FILE__, __LINE__, "Failed to start the task");
	}
	auto status = srDailyTask::getInstance()->start();
	
	switch (status)
	{
	case srDailyTask::error:
		Log::error(__FILE__, __LINE__, "srDailyTask exited with errors");
		break;
	case srDailyTask::stop:
		Log::info(__FILE__, __LINE__, "srDailyTask stoped normally");
		break;
	case srDailyTask::unknown:
		Log::warning(__FILE__, __LINE__, "srDailyTask exited with known status");
		break;
	default:
		break;
	}

	Log::info(__FILE__, __LINE__, "srDailyTask stopped!");
    return 0;
}