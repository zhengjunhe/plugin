default:
	@cp /home/hezhengjun/work/c_code/openjdk/jvm4chain33/build/linux-x86_64-normal-server-release/jdk/lib/amd64/jli/libjli.so ../../so
	@go build -o jvm_run_tx_query jvm_run_tx_query.go

jvm_one_thread:
	@go build -o jvm_one_thread jvm_one_thread.go

single:
	@go build -o singleJvm

clean:
	@rm jvm_run_tx_query
