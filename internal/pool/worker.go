package pool

type Worker struct {
	id        int
	completed int
	failed    int
	total     int
}
