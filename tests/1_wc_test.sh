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
# first word-count

# generate the correct output
./bin/sequential plugins/wc/wc.so inputs/pg-*.txt || exit 1
sort outputs/mr-out-* > outputs/mr-correct-wc.txt
rm -r outputs/mr-out*.txt intermediate/mr-*.txt

echo '***' Starting wc test.

./bin/cordinator inputs/pg-*.txt &
pid=$!

# give the coordinator time to create the sockets.
sleep 1

# start multiple workers.
./bin/worker plugins/wc/wc.so &
./bin/worker plugins/wc/wc.so &
./bin/worker plugins/wc/wc.so &

# wait for the coordinator to exit.
wait $pid

# since workers are required to exit when a job is completely finished,
# and not before, that means the job has finished.
sort outputs/mr-out* | grep . > outputs/mr-wc-all.txt
if cmp outputs/mr-wc-all.txt outputs/mr-correct-wc.txt
then
  echo '---' wc test: PASS
else
  echo '---' wc output is not the same as mr-correct-wc.txt
  echo '---' wc test: FAIL
  failed_any=1
fi


# wait for remaining workers and coordinator to exit.
wait
