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
echo '***' Starting crash test.

# generate the correct output
./bin/sequential plugins/nocrash/nocrash.so inputs/pg-*.txt || exit 1
sort outputs/mr-out-* > outputs/mr-correct-crash.txt
rm -r outputs/mr-out*.txt intermediate/mr-*.txt

rm -f mr-done
(./bin/cordinator inputs/pg-*.txt; touch mr-done) &
sleep 1

./bin/worker plugins/crash/crash.so &

# mimic rpc.go's coordinatorSock()
SOCKNAME=/var/tmp/ShilpThapak-mr-`id -u`

( while [ -e $SOCKNAME -a ! -f mr-done ]
  do
    ./bin/worker plugins/crash/crash.so
    sleep 1
  done ) &

( while [ -e $SOCKNAME -a ! -f mr-done ]
  do
    ./bin/worker plugins/crash/crash.so
    sleep 1
  done ) &

while [ -e $SOCKNAME -a ! -f mr-done ]
do
  ./bin/worker plugins/crash/crash.so
  sleep 1
done

wait

rm $SOCKNAME
sort outputs/mr-out* | grep . > outputs/mr-crash-all.txt
if cmp outputs/mr-crash-all.txt outputs/mr-correct-crash.txt
then
  echo '---' crash test: PASS
else
  echo '---' crash output is not the same as outputs/mr-correct-crash.txt
  echo '---' crash test: FAIL
  failed_any=1
fi

rm -f mr-done

