#git clone https://github.com/GrantZheng/kit.git
#cd kit
#go install
#kit new service test
./kit n s test
#./kit g s test
./kit g s test --dmw
./kit g s test --gorilla
./kit g s test --endpoint-mdw
./kit g s test --dmw -t grpc
./kit g d test
./kit g c test
./kit g c test -t grpc -i test/pkg/grpc/pb
