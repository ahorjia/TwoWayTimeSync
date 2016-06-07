package main

import "precisesleep"
import "time"
import "fmt"
func main() {

   currentTime := time.Now();
   sleepDeadline := currentTime.Add(500*time.Millisecond);
   timeNano64 := sleepDeadline.UnixNano();

   precisesleep.SleepUntil64(timeNano64);
   endTime  := time.Now();

   fmt.Printf("Start time: %v\nTargt time: %v\nEnd__ time: %v\n",currentTime.UnixNano(),timeNano64,endTime.UnixNano());

   testTimeTime();

   testSleepFor();

   testSleepFor64();
}

func testTimeTime() {

   currentTime := time.Now();
   sleepDeadline := currentTime.Add(500*time.Millisecond);
   timeNano64 := sleepDeadline.UnixNano();

   precisesleep.SleepUntil(sleepDeadline);
   endTime  := time.Now();

   fmt.Printf("Start time: %v\nTargt time: %v\nEnd__ time: %v\n",currentTime.UnixNano(),timeNano64,endTime.UnixNano());
}
func testSleepFor() {

   currentTime := time.Now();
   sleepDuration := 500 * time.Millisecond;
   sleepDeadline := currentTime.Add(sleepDuration);
   timeNano64 := sleepDeadline.UnixNano();

   precisesleep.SleepFor(sleepDuration);
   endTime  := time.Now();

   fmt.Printf("Start time: %v\nTargt time: %v\nEnd__ time: %v\n",currentTime.UnixNano(),timeNano64,endTime.UnixNano());
}

func testSleepFor64() {

   currentTime := time.Now();
   sleepDuration := 500 * time.Millisecond;
   sleepDeadline := currentTime.Add(sleepDuration);
   timeNano64 := sleepDeadline.UnixNano();

   precisesleep.SleepFor64(sleepDuration.Nanoseconds());
   endTime  := time.Now();

   fmt.Printf("Start time: %v\nTargt time: %v\nEnd__ time: %v\n",currentTime.UnixNano(),timeNano64,endTime.UnixNano());
}
