#!/usr/bin/env bash

#
# map-reduce tests
#
ROOT_DIR=$(pwd)
echo $ROOT_DIR


# run the test in a fresh sub-directory.
rm -rf outputs/mr-*.txt 
rm -rf intermediate/mr-*.txt

# make sure software is freshly built.
go build -buildmode=plugin -race -o plugins/wc/wc.so plugins/wc/wc.go || exit 1
go build -buildmode=plugin -race -o plugins/indexer/indexer.so plugins/indexer/indexer.go || exit 1
go build -buildmode=plugin -race -o plugins/mtiming/mtiming.so plugins/mtiming/mtiming.go || exit 1
go build -buildmode=plugin -race -o plugins/rtiming/rtiming.so plugins/rtiming/rtiming.go || exit 1
go build -buildmode=plugin -race -o plugins/jobcount/jobcount.so plugins/jobcount/jobcount.go || exit 1
go build -buildmode=plugin -race -o plugins/early_exit/early_exit.so plugins/early_exit/early_exit.go || exit 1
go build -buildmode=plugin -race -o plugins/crash/crash.so plugins/crash/crash.go || exit 1
go build -buildmode=plugin -race -o plugins/nocrash/nocrash.so plugins/nocrash/nocrash.go || exit 1

cd cmd/sequential && go build -race -o $ROOT_DIR/bin sequential.go && cd $ROOT_DIR
cd cmd/cordinator && go build -race -o $ROOT_DIR/bin cordinator.go && cd $ROOT_DIR
cd cmd/worker && go build -race -o $ROOT_DIR/bin worker.go && cd $ROOT_DIR

failed_any=0

#########################################################
# test whether any worker or coordinator exits before the
# task has completed (i.e., all output files have been finalized)
rm -f mr-*

echo '***' Starting early exit test.

DF=anydone$$
rm -f $DF

(./bin/cordinator inputs/pg-*.txt; touch $DF) &
sleep 1

# start multiple workers.
(./bin/worker plugins/early_exit/early_exit.so; touch $DF) &
(./bin/worker plugins/early_exit/early_exit.so; touch $DF) &
(./bin/worker plugins/early_exit/early_exit.so; touch $DF) &

# wait for any of the coord or workers to exit.
# `jobs` ensures that any completed old processes from other tests
# are not waited upon.
jobs &> /dev/null
if [[ "$OSTYPE" = "darwin"* ]]
then
  # bash on the Mac doesn't have wait -n
  while [ ! -e $DF ]
  do
    sleep 0.2
  done
else
  # the -n causes wait to wait for just one child process,
  # rather than waiting for all to finish.
  wait -n
fi

# a process has exited. this means that the output should be finalized
# otherwise, either a worker or the coordinator exited early
sort outputs/mr-out* | grep . > outputs/mr-wc-all-initial.txt

# wait for remaining workers and coordinator to exit.
wait

# compare initial and final outputs
sort outputs/mr-out* | grep . > outputs/mr-wc-all-final.txt
if cmp outputs/mr-wc-all-final.txt outputs/mr-wc-all-initial.txt
then
  echo '---' early exit test: PASS
else
  echo '---' output changed after first worker exited
  echo '---' early exit test: FAIL
  failed_any=1
fi

rm -rf outputs/mr-*.txt 
rm -rf intermediate/mr-*.txt
rm -rf $DF*

exit 1
