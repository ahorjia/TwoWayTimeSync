package precisesleep

//#include "precisesleep.h"
import "C"
import "time"

func SleepUntil(deadline time.Time ) {

   timeNano64 := deadline.UnixNano();
   timeSecs := C.int32_t( timeNano64 / 1e9);
   timeNsecs := C.int32_t( timeNano64 % 1e9);
   C.sleepUntil(timeSecs,timeNsecs);

}

func SleepUntil64(timeNano64 int64 ) {

   timeSecs := C.int32_t( timeNano64 / 1e9);
   timeNsecs := C.int32_t( timeNano64 % 1e9);
   C.sleepUntil(timeSecs,timeNsecs);

}

func SleepFor(amount time.Duration ) {
   timeNano64 := amount.Nanoseconds();
   timeSecs := C.int32_t( timeNano64 / 1e9);
   timeNsecs := C.int32_t( timeNano64 % 1e9);
   C.sleepFor(timeSecs,timeNsecs);

}

func SleepFor64(timeNano64 int64 ) {

   timeSecs := C.int32_t( timeNano64 / 1e9);
   timeNsecs := C.int32_t( timeNano64 % 1e9);
   C.sleepFor(timeSecs,timeNsecs);

}
