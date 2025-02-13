syntax = "proto3";
package user;

service UserService {
    rpc GetUser (UserRequest) returns (UserResponse);
}

message UserRequest {
    string user_id = 1;
}

message UserResponse {
    string id = 1;
    string name = 2;
    string email = 3;
}



# Build and run services

cd microservices-demo
docker-compose build
docker-compose up

# Test the API Gateway

curl http://localhost:8080/users
curl http://localhost:8080/orders

# Scale a service

docker-compose scale order-service=3
