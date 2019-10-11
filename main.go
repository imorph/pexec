package main

func jobExec2(jobPac []Job, poolSize int, maxErrNum int) {
	wp := NewWorkerPool(poolSize, len(jobPac))
	wp.Start()
	defer wp.Stop()
	for _, j := range jobPac {
		wp.SubmitJob(j)
	}
	wp.WaitResults(maxErrNum, len(jobPac))
}

func main() {

	// no errors case, all jobs equal
	var myJobs2 []Job
	for i := 0; i < 10; i++ {
		if i == 5 {
			myJobs2 = append(myJobs2, job2sErr)
		} else {
			myJobs2 = append(myJobs2, job1s)
		}
	}
	jobExec2(myJobs2, 2, 1)
}
