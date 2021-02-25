#git clone https://github.com/GrantZheng/kit.git
#cd kit
#go install
#kit new service test
./kit n s test
./kit g s test
./kit g s test --dmw # to create the default middleware
./kit g s test --gorilla # to create the default middleware
./kit g s test --endpoint-mdw  # to create the default middleware
./kit g s test -t grpc # specify the transport (default is http)
./kit g d test
./kit g c test
./kit g c test -t grpc -i test/pkg/grpc/pb
