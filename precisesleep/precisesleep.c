#include "precisesleep.h"
int sleepUntil( int32_t seconds,int32_t nanoseconds)
{
   struct timespec sSleepTime,sRemainder;
   int returnVal;
   sSleepTime.tv_sec = seconds;
   sSleepTime.tv_nsec = nanoseconds;
//   printf("target sec: %d, target nsec: %d\n",seconds,nanoseconds);
   returnVal = clock_nanosleep(CLOCK_REALTIME,TIMER_ABSTIME,&sSleepTime,&sRemainder);
   return returnVal;
}

int sleepFor( int32_t seconds,int32_t nanoseconds)
{
   struct timespec sSleepTime,sRemainder;
   int returnVal;
   sSleepTime.tv_sec = seconds;
   sSleepTime.tv_nsec = nanoseconds;
//   printf("target sec: %d, target nsec: %d\n",seconds,nanoseconds);
   returnVal = clock_nanosleep(CLOCK_REALTIME,0,&sSleepTime,&sRemainder);
   return returnVal;
}
