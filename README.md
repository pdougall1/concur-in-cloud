# Someone was talkin' bout go routines don't matter!

The concern here is that cloud deployed apps might not reserve enough compute for concurrency to matter, there would only be at most one "process" at a time to work with.  The concern might be only relevant to k8s and freactional microprocesses or whatever.  What I care about is if it make sense to build data pipelines in cloud run in go that dependend on concurrent processing in a way that speeds things up.

Admittedly, there are at least 2 relevant use cases.  Goroutines can be held up buy IO (db calls) or through actual processing (computations on data for reports we'll need (in my awesome startup)).  Gonna have to check on both.

Gonna use gcloud pub/sub to frovide load, it should provide continuous demand.