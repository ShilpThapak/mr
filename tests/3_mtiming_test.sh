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
echo '***' Starting map parallelism test.

./bin/cordinator inputs/pg-*.txt &
sleep 1

./bin/worker plugins/mtiming/mtiming.so &
./bin/worker plugins/mtiming/mtiming.so

NT=`cat outputs/mr-out* | grep '^times-' | wc -l | sed 's/ //g'`
if [ "$NT" != "2" ]
then
  echo '---' saw "$NT" workers rather than 2
  echo '---' map parallelism test: FAIL
  failed_any=1
fi

if cat outputs/mr-out* | grep '^parallel.* 2' > /dev/null
then
  echo '---' map parallelism test: PASS
else
  echo '---' map workers did not run in parallel
  echo '---' map parallelism test: FAIL
  failed_any=1
fi

wait

exit 1